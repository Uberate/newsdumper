package staged

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/uberate/gf"
)

type Staged interface {
	gf.Entity
	RunE(ctx context.Context) error
	SetNext(staged Staged)
}

type CommonStaged struct {
	gf.BaseEntity

	RunFunc       func(ctx context.Context) (context.Context, error)
	IgnoreErr     bool
	IgnoreNextErr bool
	NextStaged    Staged
	Logger        *logrus.Logger
}

func (cs *CommonStaged) SetNext(staged Staged) {
	cs.NextStaged = staged
}

func (cs *CommonStaged) RunE(ctx context.Context) error {
	if cs.RunFunc != nil {
		var err error
		ctx, err = cs.RunFunc(ctx)
		if err != nil && !cs.IgnoreErr {
			return err
		}
	}

	if cs.NextStaged != nil {
		var err error
		err = cs.NextStaged.RunE(ctx)
		if err != nil && !cs.IgnoreNextErr {
			return err
		}
	}

	return nil
}

var (
	RunFuncKey         = "run-func-key"
	IgnoreCurrentErr   = "ignore-current-err"
	IgnoreSubStagedErr = "ignore-sub-staged-err"
)

func GenerateCommonStaged(name string, config interface{}, logger *logrus.Logger) (Staged, error) {
	if m, ok := TryToMap(config); ok {
		rf, has := GetFromMap(m, RunFuncKey, func(ctx context.Context) (context.Context, error) { return ctx, nil })
		if !has {
			return nil, fmt.Errorf("can not build common staged: rune func is nil in config, key: %s", RunFuncKey)
		}
		ignoreErr, _ := GetFromMap(m, IgnoreCurrentErr, true)
		ignoreNextStagedErr, _ := GetFromMap(m, IgnoreSubStagedErr, true)
		res := &CommonStaged{
			BaseEntity:    gf.NewBaseEntityGenerator(CommonKind, V1Version)(name),
			RunFunc:       rf,
			Logger:        logger,
			IgnoreErr:     ignoreErr,
			IgnoreNextErr: ignoreNextStagedErr,
		}

		return res, nil
	}

	return nil, fmt.Errorf("config type err, can't convert to map")
}

func AbsStaged(kind, version string,
	runFunc func(ctx context.Context) (context.Context, error),
	ignoreCurrentErr, ignoreNextStagedErr bool) gf.Generator[Staged] {

	metaGenerator := gf.NewBaseEntityGenerator(kind, version)

	return func(name string, config interface{}, logger *logrus.Logger) (Staged, error) {
		res := &CommonStaged{
			BaseEntity:    metaGenerator(name),
			RunFunc:       runFunc,
			Logger:        logger,
			IgnoreErr:     ignoreCurrentErr,
			IgnoreNextErr: ignoreNextStagedErr,
		}

		return res, nil
	}
}

type Builder struct {
	start   Staged
	current Staged
	logger  *logrus.Logger
}

func NewBuilder(logger *logrus.Logger) *Builder {
	start, _, _ := StageFactory.Get(ScopeKind, V1Version, "start", nil, logger)
	return &Builder{
		start:   start,
		current: start,
	}
}

func (b *Builder) NextStaged(next Staged) *Builder {
	b.current.SetNext(next)
	b.current = next
	return b
}

func (b *Builder) Next(name string, f func(ctx context.Context) (context.Context, error), ignoreCurrentError, ignoreNextError bool) *Builder {
	m := map[string]any{
		RunFuncKey:         f,
		IgnoreCurrentErr:   ignoreCurrentError,
		IgnoreSubStagedErr: ignoreCurrentError,
	}
	next, err, _ := StageFactory.Get(CommonKind, V1Version, name, m, b.logger)
	if err != nil {
		panic(err)
	}

	b.current.SetNext(next)
	b.current = next
	return b
}

func (b *Builder) Scope(name string) *Builder {
	next, _, _ := StageFactory.Get(ScopeKind, V1Version, name, nil, b.logger)
	b.current.SetNext(next)
	b.current = next
	return b
}

func (b *Builder) Build() Staged {
	return b.start
}
