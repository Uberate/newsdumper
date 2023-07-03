package logics

import (
	"context"
	"github.com/sirupsen/logrus"
	"news/pkg/getter"
	"news/pkg/staged"
	"time"
)

func GetterStaged(ctx context.Context) (context.Context, error) {
	callTime, _ := staged.GetFromContext(ctx, CallTime, time.Now())

	getters, ok := staged.GetFromContext(ctx, GetterInstances, []getter.NewsGetter{})
	if !ok {
		return ctx, nil
	}

	var defaultLog *logrus.Logger
	log, ok := staged.GetFromContext(ctx, LoggerInstance, defaultLog)
	if log == nil {
		log = logrus.New()
	}

	var newsRes []getter.News
	// get news
	for _, item := range getters {
		res, err := item.GetNews(callTime.Unix())
		if err != nil {
			log.Error(err)
			continue
		}
		newsRes = append(newsRes, res...)
	}

	ctx = staged.SetFromContext(ctx, NewsKey, newsRes)

	return ctx, nil
}
