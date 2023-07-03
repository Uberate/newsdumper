package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/cmd/bin/consts"
	"news/internal/logics"
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
