package getter

import (
	"news/pkg/factory"
)

var NewGetterFactory *factory.Factory[NewsGetter]

func init() {
	NewGetterFactory = factory.NewFactor[NewsGetter](&EmptyGetter{})

	NewGetterFactory.Registry("sina", "v1", SimpleNewsGetter("sina", "v1", "https://sina-news.vercel.app/rss.json", nil, SinaV1ResParser))
}
