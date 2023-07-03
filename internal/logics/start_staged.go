package logics

import (
	"context"
	"github.com/sirupsen/logrus"
	"news/pkg/getter"
	"news/pkg/hooks"
	"news/pkg/staged"
	"time"
)

func StartStaged(ctx context.Context) (context.Context, error) {
	ctx = staged.SetFromContext(ctx, CallTime, time.Now())

	var initLogger *logrus.Logger
	var ok bool

	if initLogger, ok = staged.GetFromContext(ctx, LoggerInstance, initLogger); !ok {
		initLogger = logrus.New()
		ctx = staged.SetFromContext(ctx, LoggerInstance, initLogger)
	}
	if _, ok = staged.GetFromContext(ctx, GetterInstances, []getter.NewsGetter{}); !ok {
		initLogger.Warnf("not found any getter")
		ctx = staged.SetFromContext(ctx, LoggerInstance, []getter.NewsGetter{})
	}
	if _, ok = staged.GetFromContext(ctx, HookerInstances, []hooks.Hook{}); !ok {
		initLogger.Warnf("not found any hooker")
		ctx = staged.SetFromContext(ctx, HookerInstances, []hooks.Hook{})
	}

	return ctx, nil
}
