package getter

import (
	"github.com/uberate/gf"
)

var NewGetterFactory *gf.Factory[NewsGetter]

func init() {
	NewGetterFactory = gf.NewFactor[NewsGetter](&EmptyGetter{})

	NewGetterFactory.Registry("sina", "v1", SimpleNewsGetter("sina", "v1", "https://sina-news.vercel.app/rss.json", nil, SinaV1ResParser))
}
