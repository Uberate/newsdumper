package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/pkg/getter"
	"news/pkg/hooks"
	"time"
)

func Load(config *cfg.Config, logger *logrus.Logger) ([]getter.NewsGetter, []hooks.Hook, error) {
	logger.Tracef("start main logic\n")

	getters, err := loadGetters(config.Getters, config.EnableNotFoundGetter, logger)
	if err != nil {
		err = fmt.Errorf("load getter err: %v\n", err)
		logger.Error(err)
		return nil, nil, err
	}

	return getters, nil, nil
}

func Run(getters []getter.NewsGetter, hooks []hooks.Hook, runCron string, log *logrus.Logger) error {
	cronInstance := cron.New()
	if _, err := cronInstance.AddFunc(runCron, func() {
		var newsRes []getter.News

		getTime := time.Now().Unix()

		for _, item := range getters {
			res, err := item.GetNews(getTime)
			if err != nil {
				log.Error(err)
				continue
			}
			newsRes = append(newsRes, res...)
		}

		for _, item := range newsRes {
			fmt.Println(item)
		}
	}); err != nil {
		return err
	}

	cronInstance.Run()
	return nil
}

func loadGetters(getterConfigs []cfg.FactoryDesc, enableNotFound bool, logger *logrus.Logger) ([]getter.NewsGetter, error) {
	var res []getter.NewsGetter
	for _, item := range getterConfigs {
		if getterItem, ok := getter.NewGetterFactory.Get(item.Kind, item.Version, item.Name, item.Config, logger); !ok {
			if enableNotFound {
				loggerInstance.Warnf("not found getter: kind: [%s], version: [%s]", item.Kind, item.Version)
				continue
			}
			return nil, fmt.Errorf("not found getter: kind: [%s], version: [%s]", item.Kind, item.Version)
		} else {
			res = append(res, getterItem)
		}
	}
	return res, nil
}
