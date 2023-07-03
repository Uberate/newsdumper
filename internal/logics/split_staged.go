package logics

import (
	"context"
	"github.com/sirupsen/logrus"
	"news/cmd/bin/cfg"
	"news/pkg/getter"
	"news/pkg/staged"
	"strings"
)

func SplitStage(ctx context.Context) (context.Context, error) {

	var defaultLog *logrus.Logger
	log, ok := staged.GetFromContext(ctx, LoggerInstance, defaultLog)
	if log == nil {
		log = logrus.New()
	}

	newsRes, ok := staged.GetFromContext(ctx, NewsKey, []getter.News{})
	if !ok || len(newsRes) == 0 {
		return ctx, nil
	}

	types := map[string][]getter.News{}
	groups, ok := staged.GetFromContext(ctx, GroupKey, []cfg.MapperSet{})
	if !ok {
		// all in one typ: ""
		types[""] = newsRes
		log.Info("no group config, all news in one group")
	} else {
		// news filter
		for _, item := range newsRes {
			for _, groupKey := range groups {
				if types[groupKey.Key] == nil {
					types[groupKey.Key] = []getter.News{}
				}
				for _, keyWord := range groupKey.Values {
					if strings.Contains(item.Title, keyWord) {
						types[groupKey.Key] = append(types[groupKey.Key], item)
						continue
					}
					if strings.Contains(item.Body, keyWord) {
						types[groupKey.Key] = append(types[groupKey.Key], item)
						continue
					}
				}
			}
		}
	}

	for typ, item := range types {
		log.Debugf("typ: %s, count: %d\n", typ, len(item))
	}

	ctx = staged.SetFromContext(ctx, GroupedNewsKey, types)
	return ctx, nil
}
