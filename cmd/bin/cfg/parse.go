package cfg

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"strings"
)

func ParseConfig(configPath string) (*Config, error) {

	v := viper.NewWithOptions()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c := DefaultConfig()

	// read from config file
	v.SetConfigFile(configPath)

	// merge default config
	defaultConfig := map[string]interface{}{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  &defaultConfig,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(c); err != nil {
		return nil, err
	}
	if err := v.MergeConfigMap(defaultConfig); err != nil {
		return nil, err
	}
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	// sync config value from env
	v.AutomaticEnv()

	// unmarshal value to config instance
	if err := v.Unmarshal(c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
	}); err != nil {
		return nil, err
	}

	return c, nil
}
