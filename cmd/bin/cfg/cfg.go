package cfg

const ConfigDemo = `## ------------------------------------------------------------------------------------------------------------------ ##
## CONFIG DEMO: config.demo.yaml                                                                                      ##
##                                                                                                                    ##
## Author: Uberate                                                                                                    ##
## Email: <ubserate@gmail.com>                                                                                        ##
##                                                                                                                    ##
## This output show the config demo of application.                                                                   ##
## ------------------------------------------------------------------------------------------------------------------ ##

# group_filters 
# The group_filters will try to group the res of this filters by elements. 
#
group_filters:
  - key: "test1"
    values:
      - "test1"
### If the article body or title has this key word, these news will group to technology group.
### If some article has more than one group key word in different filters, these news will group to these groups.

# disable_getters
# To disable some getters.
#
# If you want to get all getters from the application, start the application with: '-s' or '--show-getters'
disable_getters:
- sina_news_v1
`

// Config of application.
type Config struct {
	// GroupFilters set the keys words to group the news.
	GroupFilters []MapperSet `json:"group_filters" yaml:"group_filters"`

	// DisableWebSites will disable the websites.
	DisableGetters []string `json:"disable_getters" yaml:"disable_getters"`
}

// MapperSet set a key to a set string.
type MapperSet struct {
	Key    string   `json:"key" yaml:"key"`
	Values []string `json:"values" yaml:"values"`
}
