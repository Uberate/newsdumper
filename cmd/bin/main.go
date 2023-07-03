package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/cmd/bin/consts"
	"news/internal/logics"
	"news/pkg/getter"
	"news/pkg/hooks"
	"news/pkg/log"
	"os"
)

var (
	configInstance *cfg.Config
	loggerInstance *logrus.Logger
)

func init() {
	showDefaultConfig := false

	configPath := ""

	flag.BoolVar(&showDefaultConfig, "show-default-config", false, "-show-default-config to "+
		"show the default config, this flags will stop the process.")
	// If none env setting, parse flag.
	flag.StringVar(&configPath, "config", "./conf/web.conf.yaml", "-config config-path or set env OM_CONFIG_PATH")

	flag.Parse()

	// init config info, env > flag
	if envStr := os.Getenv(consts.NewDumperConfigPathEnv); len(envStr) != 0 {
		configPath = envStr
	}

	if showDefaultConfig {
		fmt.Print(cfg.ConfigDemo)
		os.Exit(0)
	}

	if len(configPath) == 0 {
		// ignore conf info, use default config
		fmt.Println("load default config")
		c := cfg.DefaultConfig()
		configInstance = c
	} else {
		c, err := cfg.ParseConfig(configPath)
		if err != nil {
			panic(err)
		}
		configInstance = c
	}

	// init logger info
	if lg, err := log.InitLogInstance(configInstance.Log); err != nil {
		panic(err)
	} else {
		loggerInstance = lg
	}
	loggerInstance.Trace("log init done")

	// print the version info
	versionJsonBytes, err := json.Marshal(GetVersionInfo())
	if err != nil {
		loggerInstance.Error(err)
		loggerInstance.Warn("skip version check")
	}

	loggerInstance.Info("version:" + string(versionJsonBytes))
	loggerInstance.Info("init application done")
	loggerInstance.Trace("bootstrap the application")
}

// main is the bootstrap function of cmd.
func main() {
	loggerInstance.Info("scan the cron")

	getters, hookers, err := Load(configInstance, loggerInstance)
	if err != nil {
		loggerInstance.Fatalf("run stop: %v", err)
	}

	if err = logics.Run(getters, hookers, configInstance.RunCron, loggerInstance, configInstance.GroupFilters); err != nil {
		loggerInstance.Fatalf("run stop: %v", err)
	}
}

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

// ----------- version flags.

// build with flag:
// -ldflags "-w -s -X 'main.Version=${VERSION}' -X 'main.HashTag=`git rev-parse HEAD`' -X 'main.BranchName=`git rev-parse --abbrev-ref HEAD`' -X 'main.BuildDate=`date -u '+%Y-%m-%d_%I:%M:%S%p'`' -X 'main.GoVersion=`go version`'"

const (
	VersionTagVersion    = "version"
	VersionTagHashTag    = "hash-tag"
	VersionTagBranchName = "branch-name"
	VersionTagBuildDate  = "build-date"
	VersionTagGoVersion  = "go-version"
)

var Version string
var HashTag string
var BranchName string
var BuildDate string
var GoVersion string

func GetVersionInfo() map[string]string {
	return map[string]string{
		VersionTagVersion:    Version,
		VersionTagHashTag:    HashTag,
		VersionTagBranchName: BranchName,
		VersionTagBuildDate:  BuildDate,
		VersionTagGoVersion:  GoVersion,
	}
}
