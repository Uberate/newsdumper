package logics

import (
	"context"
	"news/pkg/getter"
	"news/pkg/hooks"
	"news/pkg/staged"
)

// HookerStaged
//
// Args:
// - HookerInstances: []hooks.Hook{}
// - GroupedNewsKey: map[string][]getter.News{} | FROM split_staged
//
// Result:
func HookerStaged(ctx context.Context) (context.Context, error) {
	log := LogFromCtx(ctx)

	hookers, ok := staged.GetFromContext(ctx, HookerInstances, []hooks.Hook{})
	if !ok {
		return ctx, nil
	}

	types, ok := staged.GetFromContext(ctx, GroupedNewsKey, map[string][]getter.News{})
	if !ok {
		log.Warn("there are no news, skip send staged")
	}

	// sender
	for _, item := range hookers {
		log.Infof("hook: kind: [%s], version: [%s], name: [%s]\n", item.Kind(), item.Version(), item.Name())
		for typ, news := range types {
			if len(news) == 0 {
				continue
			}
			if len(news) >= 5 {
				news = news[:4]
			}
			err := item.Hook(typ, news)
			if err != nil {
				log.Error(err)
			}
		}
	}

	return ctx, nil
}
