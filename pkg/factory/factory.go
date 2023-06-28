package factory

import (
	"github.com/sirupsen/logrus"
	"github.com/uberate/gset"
	"sync"
)

type BaseEntity struct {
	kind    string
	name    string
	version string
}

func (be *BaseEntity) Kind() string {
	return be.kind
}
func (be *BaseEntity) Name() string {
	return be.name
}
func (be *BaseEntity) Version() string {
	return be.version
}

type Entity interface {
	Kind() string
	Name() string
	Version() string
}

type Generator[T Entity] func(name string, config interface{}, logger *logrus.Logger) (T, error)

func NewFactor[T Entity](emptyPointer T) *Factory[T] {
	return &Factory[T]{
		emptyPointer: emptyPointer,
		locker:       &sync.RWMutex{},
		models:       map[string]map[string]Generator[T]{},
	}
}

type Factory[T Entity] struct {
	// inner locker
	locker *sync.RWMutex

	// models 代表工厂中的生成函数，用于初始化目标数据。
	// kinds-versions-FactoryGenerator
	models map[string]map[string]Generator[T]

	emptyPointer T
}

// Registry 注册初 Generator（初始化器），在成功注册后，返回 true，如果当前初始化器已经存在，则不会注册并返回 false。
func (f *Factory[T]) Registry(kind, version string, g Generator[T]) bool {
	f.locker.Lock()
	defer f.locker.Unlock()

	if _, ok := f.models[kind]; !ok {
		f.models[kind] = map[string]Generator[T]{}
	}

	if _, ok := f.models[kind][version]; ok {
		// already exists, return false
		return false
	}

	f.models[kind][version] = g
	return true
}

func (f *Factory[T]) RemoveVersion(kind, version string) bool {
	f.locker.Lock()
	defer f.locker.Unlock()

	if _, ok := f.models[kind]; !ok {
		return false
	}

	if _, ok := f.models[kind][version]; !ok {
		return false
	}
	delete(f.models[kind], version)
	return true
}

func (f *Factory[T]) RemoveKind(kind string) (int, bool) {
	f.locker.Lock()
	defer f.locker.Unlock()

	if _, ok := f.models[kind]; !ok {
		return 0, false
	}
	c := len(f.models[kind])

	delete(f.models, kind)
	return c, true
}

func (f *Factory[T]) ContainKind(kind string) bool {
	f.locker.RLock()
	defer f.locker.RUnlock()
	_, ok := f.models[kind]
	return ok
}

func (f *Factory[T]) ContainVersion(kind, version string) bool {
	f.locker.RLock()
	defer f.locker.RUnlock()
	_, ok := f.unsafeGetGenerator(kind, version)
	return ok
}

func (f *Factory[T]) ListKinds() gset.Set[string] {
	f.locker.RLock()
	defer f.locker.RUnlock()

	return gset.FromMapKey(f.models)
}

func (f *Factory[T]) ListVersions(kind string) gset.Set[string] {
	f.locker.RLock()
	defer f.locker.RUnlock()

	if _, ok := f.models[kind]; !ok {
		return gset.FromArray([]string{})
	}
	return gset.FromMapKey(f.models[kind])
}

func (f *Factory[T]) Get(kind, version, name string, config interface{}, logger *logrus.Logger) (T, error, bool) {
	f.locker.RLock()
	defer f.locker.RUnlock()

	if g, exists := f.unsafeGetGenerator(kind, version); exists {
		res, err := g(name, config, logger)
		if err == nil {
			return res, nil, true
		}
		return f.emptyPointer, err, true
	}
	return f.emptyPointer, nil, false
}

func (f *Factory[T]) GetGenerator(kind, version string) (Generator[T], bool) {
	f.locker.RLock()
	defer f.locker.RUnlock()

	return f.unsafeGetGenerator(kind, version)
}

func (f *Factory[T]) unsafeGetGenerator(kind, version string) (Generator[T], bool) {
	if _, ok := f.models[kind]; !ok {
		return nil, false
	}
	if _, ok := f.models[kind][version]; !ok {
		return nil, false
	}

	return f.models[kind][version], true
}
