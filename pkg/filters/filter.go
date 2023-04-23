package filters

import (
	"github.com/uberate/gset"
	"news/pkg/getter"
	"strings"
)

type Filter struct {
	Name  string
	Items gset.Set[string]
}

func (f *Filter) FromString(filterStr string) {
	nameAndItems := strings.Split(filterStr, ":")

	f.Name = nameAndItems[0]
	if len(nameAndItems) > 1 {
		strings.Join(nameAndItems[1:], ":")

	}
}

func (f *Filter) Match(news getter.News) bool {
	title := news.Title
	body := news.Body

	for item := range f.Items {
		if strings.Contains(title, item) || strings.Contains(body, item) {
			return true
		}
	}

	return false
}
