package cfg

import (
	"news/pkg/log"
	"time"
)

const ConfigDemo = `## ------------------------------------------------------------------------------------------------------------------ ##
## CONFIG DEMO: config.demo.yaml                                                                                      ##
##                                                                                                                    ##
## Author: Uberate                                                                                                    ##
## Email: <ubserate@gmail.com>                                                                                        ##
##                                                                                                                    ##
## This output show the config demo of application.                                                                   ##
## ------------------------------------------------------------------------------------------------------------------ ##

# group_filters 
# The group_filters will try to group the res of this filters by elements. 
#
group_filters:
  - key: "test1"
    values:
      - "test1"
### If the article body or title has this key word, these news will group to technology group.
### If some article has more than one group key word in different filters, these news will group to these groups.

hookers:
- kind: smtp
  name: email-sender
  version: v1
  config:
    host: "host" 
    port: <port>
    username: <send name>
    receivers: []

# If set enable_not_found_getter false, any not found getter will panic process. 
enable_not_found_getter: false

# Define which getter will used.
getters:
- kind: sina
  version: v1
  name: sina-v1
`

func DefaultConfig() *Config {
	return &Config{
		Log: log.Config{
			Level:                     "INFO",
			DisableColor:              false,
			EnvironmentOverrideColors: false,
			DisableTimestamp:          false,
			FullTimestamp:             true,
			TimestampFormat:           time.RFC3339Nano,
		},
	}
}

// Config of application.
type Config struct {
	Log log.Config `json:"log" yaml:"log"`

	// GroupFilters set the keys words to group the news.
	GroupFilters []MapperSet `json:"group_filters" yaml:"group_filters"`

	RunCron string `json:"run_cron" yaml:"run_cron"`

	EnableNotFoundHookers bool          `json:"enable_not_found_hookers" yaml:"enable_not_found_hookers"`
	Hookers               []FactoryDesc `json:"hookers" yaml:"hookers"`

	EnableNotFoundGetter bool          `json:"enable_not_found_getter" yaml:"enable_not_found_getter"`
	Getters              []FactoryDesc `json:"getters" yaml:"getters"`
}

// MapperSet set a key to a set string.
type MapperSet struct {
	Key    string   `json:"key" yaml:"key"`
	Values []string `json:"values" yaml:"values"`
}

type FactoryDesc struct {
	Kind    string      `json:"kind,omitempty" yaml:"kind,inline"`
	Version string      `json:"version,omitempty" yaml:"version,inline"`
	Name    string      `json:"name,omitempty" yaml:"name,inline"`
	Config  interface{} `json:"config,omitempty" yaml:"config,inline"`
}
