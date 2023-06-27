package log

import "github.com/sirupsen/logrus"

// InitLogInstance will init the logger instance by config.
func InitLogInstance(config Config) (*logrus.Logger, error) {
	LoggerInstance := logrus.New()
	LoggerInstance.SetFormatter(&logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             config.DisableColor,
		EnvironmentOverrideColors: config.EnvironmentOverrideColors,
		DisableTimestamp:          config.DisableTimestamp,
		FullTimestamp:             config.FullTimestamp,
		TimestampFormat:           config.TimestampFormat,
	})

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	LoggerInstance.SetLevel(level)
	return LoggerInstance, nil
}
