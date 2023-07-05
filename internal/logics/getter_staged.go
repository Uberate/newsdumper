package logics

import (
	"context"
	"news/pkg/getter"
	"news/pkg/staged"
	"time"
)

// GetterStaged
//
// Args:
// - CallTime: time.Time | FROM StartStaged
// - GetterInstances []getter.NewsGetter{} | FROM StartStaged
//
// Result:
// - NewsKey: []getter.News
func GetterStaged(ctx context.Context) (context.Context, error) {
	callTime, _ := staged.GetFromContext(ctx, CallTime, time.Now())
	log := LogFromCtx(ctx)

	getters, ok := staged.GetFromContext(ctx, GetterInstances, []getter.NewsGetter{})
	if !ok {
		return ctx, nil
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
