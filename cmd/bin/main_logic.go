package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/pkg/getter"
	"news/pkg/hooks"
)

func Load(config *cfg.Config, logger *logrus.Logger) ([]getter.NewsGetter, []hooks.Hook, error) {
	logger.Tracef("start main logic\n")

	getters, err := loadGetters(config.Getters, config.EnableNotFoundGetter, logger)
	if err != nil {
		err = errors.Cause(err)
		logger.Error(err)
		return nil, nil, err
	}

	hookers, err := loadHooker(config.Hookers, config.EnableNotFoundHookers, logger)
	if err != nil {
		err = errors.Cause(err)
		logger.Error(err)
		return nil, nil, err
	}

	return getters, hookers, nil
}

func loadHooker(hookConfigs []cfg.FactoryDesc, enableNotFound bool, logger *logrus.Logger) ([]hooks.Hook, error) {
	var res []hooks.Hook
	for _, item := range hookConfigs {
		hookItem, err, ok := hooks.HookFactory.Get(item.Kind, item.Version, item.Name, item.Config, logger)
		if !ok {
			if enableNotFound {
				loggerInstance.Warnf("not found hook: kind: [%s], version: [%s]", item.Kind, item.Version)
				continue
			}
			return nil, fmt.Errorf("not found hook: kind: [%s], version: [%s]", item.Kind, item.Version)
		} else if err != nil {
			return nil, fmt.Errorf("init hook: [%s], version: [%s] has err: %s", item.Kind, item.Version, err)
		} else {
			res = append(res, hookItem)
		}
	}
	return res, nil
}

func loadGetters(getterConfigs []cfg.FactoryDesc, enableNotFound bool, logger *logrus.Logger) ([]getter.NewsGetter, error) {
	var res []getter.NewsGetter
	for _, item := range getterConfigs {
		getterItem, err, ok := getter.NewGetterFactory.Get(item.Kind, item.Version, item.Name, item.Config, logger)
		if !ok {
			if enableNotFound {
				loggerInstance.Warnf("not found getter: kind: [%s], version: [%s]", item.Kind, item.Version)
				continue
			}
			return nil, fmt.Errorf("not found getter: kind: [%s], version: [%s]", item.Kind, item.Version)
		} else if err != nil {
			return nil, fmt.Errorf("init getter: [%s], version: [%s] has err: %s", item.Kind, item.Version, err)
		} else {
			res = append(res, getterItem)
		}
	}
	return res, nil
}
