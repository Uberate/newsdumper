package logics

import (
	"context"
	"github.com/sirupsen/logrus"
	"news/pkg/staged"
)

func LogFromCtx(ctx context.Context) (T *logrus.Logger) {
	var defaultLog *logrus.Logger
	log, _ := staged.GetFromContext(ctx, LoggerInstance, defaultLog)
	if log == nil {
		log = logrus.New()
	}

	return log

}
