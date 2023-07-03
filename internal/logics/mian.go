package logics

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/pkg/getter"
	"news/pkg/hooks"
	"news/pkg/staged"
	"time"
)

func Run(getters []getter.NewsGetter, hooks []hooks.Hook, cronStr string, log *logrus.Logger, groups []cfg.MapperSet) error {
	var InnerStaged = staged.NewBuilder(log)

	ctx := staged.NewShared(context.Background())
	ctx = staged.SetFromContext(ctx, GetterInstances, getters)
	ctx = staged.SetFromContext(ctx, HookerInstances, hooks)
	ctx = staged.SetFromContext(ctx, LoggerInstance, log)
	ctx = staged.SetFromContext(ctx, GroupKey, groups)

	// 优化 Factory
	startStaged, _, _ := staged.StageFactory.Get(FlowStartKind, V1Version, "start", nil, log)
	getterStaged, _, _ := staged.StageFactory.Get(FlowGetNewsKind, V1Version, "get news", nil, log)
	splitStaged, _, _ := staged.StageFactory.Get(FlowSplitNewsKind, V1Version, "group news", nil, log)
	sendStaged, _, _ := staged.StageFactory.Get(FlowSendNewsKind, V1Version, "send news", nil, log)

	flow := InnerStaged.
		NextStaged(startStaged).
		NextStaged(getterStaged).
		NextStaged(splitStaged).
		NextStaged(sendStaged).
		Next("debug", func(ctx context.Context) (context.Context, error) {
			callTime, _ := staged.GetFromContext(ctx, CallTime, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))
			fmt.Println("run time at:", callTime.Unix())
			return ctx, nil
		}, true, true).
		Build()

	cronInstance := cron.New()
	if _, err := cronInstance.AddFunc(cronStr, func() {
		err := flow.RunE(ctx)
		if err != nil {
			log.Error(err)
		}
	}); err != nil {
		return err
	}

	cronInstance.Run()

	return nil
}
