package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
	"github.com/uberate/gset"
	"gopkg.in/yaml.v2"
	"news/cmd/bin/cfg"
	"news/pkg/getter"
	"news/pkg/hooks"
	"news/pkg/utils"
	"os"
	"path"
	"strings"
	"time"
)

// Options describe the command flags.
type Options struct {
	ConfigPath string

	// Filters group the news by filters, and will generate the filters. If not set any filter, will do not group by.
	Filters []string

	//FilterConf string

	// AggregateNewsWebSite if enable the aggregate news by website, the output will in one file(if set the Filters too,
	// only group by Filters), else will output to file group by website(and group by filters).
	AggregateNewsWebSite bool

	// OutputDir define the output path. If enable the OutputForce, the output dir will be created by tools when it not
	// exists. The output json file will write to here.
	// And if output file is exists, it will panic(unless set OutputForce, it will cover old value).
	OutputDir string

	// OutputForce will write to file force, if path not exists, the tool will create it. And if file is exists, cover
	// the old value directly.
	OutputForce bool

	// OutputWithDateTime will append the time str to file name.
	//
	// The output format is '"website name"-"filter name".json'. If enable OutputDateTime, the output will like this
	// format: '"website name"-"filter name"-20220202020202.json'.
	OutputWithDateTime bool

	// DisableGetters will disable specify website getter.
	DisableGetters []string

	// CronStr define when time(and how long) to run the application.
	CronStr string

	//=========================== show flags

	// ShowVersion will output the version of tools, and stop the process.
	ShowVersion bool

	// ShowGetters will output the getters of news, and stop the process.
	ShowGetters bool

	ShowConfigDemo bool
}

var OptionsInstance = Options{}

func init() {
	parseFlags()

	showAbleFlagPreCheck()
	checkDir()
}

func parseFlags() {

	pflag.StringVarP(&OptionsInstance.ConfigPath, "config", "c", "", "specify the "+
		"config yaml file.")
	pflag.StringArrayVarP(&OptionsInstance.Filters, "filters", "f", []string{}, "group the news"+
		"by filters, and will generate the filters. If not set any filter, will do no group by. The filter format in'"+
		"-fcar:tsl,TSL' will group to car by filter 'tsl,TSL,car'.")
	//pflag.StringVarP(&OptionsInstance.FilterConf, "")
	pflag.BoolVarP(&OptionsInstance.AggregateNewsWebSite, "aggregate-news-website", "a", false,
		"if enable the aggregate news by web site, the output will in one file(if set the filters too, only "+
			"group by filters), else will output to file group by website(and group by filters).")
	pflag.BoolVarP(&OptionsInstance.OutputForce, "force", "F", false,
		"")
	pflag.StringVarP(&OptionsInstance.OutputDir, "output-dir", "o", "./",
		"")
	pflag.BoolVarP(&OptionsInstance.OutputWithDateTime, "output-with-date-time", "t", false,
		"output file name will append the date-time in format: "+
			"'\"website name\"-\"filter name\"-20220202020202.json'")
	pflag.StringSliceVarP(&OptionsInstance.DisableGetters, "disable-getter", "d", []string{},
		"to disable specify getter, about all getter name, use -s | --show-getters.")
	//pflag.StringSliceVarP(&OptionsInstance.DisableGetters)
	pflag.BoolVarP(&OptionsInstance.ShowVersion, "version", "v", false, "show the "+
		"version of the application, if enable the version flag, the application will stop directly.")
	pflag.BoolVarP(&OptionsInstance.ShowGetters, "show-getters", "s", false, "show"+
		"the news-getters, if enable the show-getters, the application will stop directly.")
	pflag.BoolVarP(&OptionsInstance.ShowConfigDemo, "show-config-demo", "S", false,
		"show config demo of application, if enable the show-config-demo, the application will stop directly.")
	pflag.StringVarP(&OptionsInstance.CronStr, "cron-str", "r", "@every 1h",
		"set when time(and how long) to run the application, follow the cron flag. Set to empty value will run"+
			"once.")

	pflag.Parse()

}

func checkDir() {
	// check output dir is exists.
	if _, err := os.Stat(OptionsInstance.OutputDir); err != nil {
		if os.IsNotExist(err) {
			if OptionsInstance.OutputForce {
				if err = os.MkdirAll(OptionsInstance.OutputDir, os.ModePerm); err != nil {
					panic(err)
				}

				return
			} else {
				fmt.Printf("output dir [%s] not exists, use --force | -F, or create dir by your self.\n",
					OptionsInstance.OutputDir)
				os.Exit(2)
			}
		}
		panic(err)
	}
}

// showAbleFlagPreCheck will check the show flags, and if some flag is enabled, stop the process directly, and output
// value about it.
func showAbleFlagPreCheck() {
	if OptionsInstance.ShowVersion {
		showVersion()
	}
	if OptionsInstance.ShowGetters {
		showGetters()
	}
	if OptionsInstance.ShowConfigDemo {
		showConfigDemo()
	}
}

func showVersion() {
	// todo
	fmt.Println("TODO: show version")
	os.Exit(0)
}

func showGetters() {
	for item := range gset.FromMapKey(getter.NewsGetters) {
		fmt.Println("-", item)
	}
	os.Exit(0)
}

func showConfigDemo() {
	fmt.Println(cfg.ConfigDemo, "")
	os.Exit(0)
}

func parseFiltersOfCommand(value string) (string, gset.Set[string]) {
	values := strings.Split(value, ",")
	k := ""
	set := gset.FromArray([]string{})
	if len(values) > 0 {
		keyAndSets := strings.Split(values[0], ":")
		k = keyAndSets[0]
		if len(keyAndSets) > 1 {
			for _, item := range keyAndSets[1:] {
				set.Push(item)
			}
		}
	}
	if len(values) > 1 {
		set.Push(values[1:]...)
	}

	return k, set
}

// main is the bootstrap function of cmd.
func main() {
	// parse the config
	config := cfg.Config{}
	if len(OptionsInstance.ConfigPath) != 0 {
		// try read config
		configValue, err := os.ReadFile(OptionsInstance.ConfigPath)
		if err != nil {
			// read file error
			panic(err)
		}

		if err = yaml.Unmarshal(configValue, &config); err != nil {
			panic(err)
		}
	}

	// append command filters.
	for _, item := range OptionsInstance.Filters {
		key, filterItem := parseFiltersOfCommand(item)
		for _, value := range config.GroupFilters {
			if value.Key == key {
				value.Values = append(value.Values, filterItem.ToArray()...)
				continue
			}
		}
		config.GroupFilters = append(config.GroupFilters, cfg.MapperSet{
			Key:    key,
			Values: filterItem.ToArray(),
		})
	}

	// append command disable-getters.
	for _, item := range OptionsInstance.DisableGetters {
		config.DisableGetters = append(config.DisableGetters, item)
	}

	var hookers []hooks.Hooks
	for _, item := range config.Hookers {
		sender, err := hooks.GetHook(item.Type, item.Config)
		if err != nil {
			fmt.Printf("init hook err: %s\n", err)
		}
		hookers = append(hookers, sender)
	}

	if len(OptionsInstance.CronStr) == 0 {
		mainLogic(config, hookers)
	} else {
		cInstance := cron.New(cron.WithSeconds())

		if _, err := cInstance.AddFunc(OptionsInstance.CronStr, func() {
			mainLogic(config, hookers)
		}); err != nil {
			fmt.Println(err)
		}

		cInstance.Run()
	}

}

func mainLogic(config cfg.Config, sender []hooks.Hooks) {
	nowTime := time.Now()
	nowTimeSecond := nowTime.Unix()
	nowTimeStr := nowTime.Format("20060102150405")

	disableGetterSets := gset.FromArray(config.DisableGetters)
	groupKeys := map[string]gset.Set[string]{}

	for _, groupFilter := range config.GroupFilters {
		if value, ok := groupKeys[groupFilter.Key]; ok {
			value.Push(groupFilter.Key)
			value.Push(groupFilter.Values...)
		} else {
			groupKeys[groupFilter.Key] = gset.FromArray(groupFilter.Values)
			groupKeys[groupFilter.Key].Push(groupFilter.Key)
		}
	}

	for getterName, item := range getter.NewsGetters {
		if disableGetterSets.Has(getterName) {
			continue
		}

		fmt.Print("Try to get from getter: ", getterName, ": total: ")
		newsGroups := map[string][]getter.News{}
		res, err := item.GetNews(nowTimeSecond)
		if err != nil {
			fmt.Printf("\nGetter [%s] has error: %v", getterName, err)
			continue
		}
		fmt.Println(len(res))

		//split group
		for filterName, filter := range groupKeys {
			fmt.Print("Match for group: ", filterName, ": count: ")
			for _, newsItem := range res {
				found := false
				for keyWord := range filter {
					if len(keyWord) == 0 {
						// skip empty filter
						continue
					}
					if strings.Contains(newsItem.Title, keyWord) ||
						strings.Contains(newsItem.Body, keyWord) {
						found = true
						newsItem.Label = append(newsItem.Label, keyWord)
					}
				}
				if found {
					if _, ok := newsGroups[filterName]; !ok {
						newsGroups[filterName] = []getter.News{}
					}

					newsGroups[filterName] = append(newsGroups[filterName], newsItem)
				}
			}
			fmt.Println(len(newsGroups[filterName]))
		}

		for filterName, news := range newsGroups {
			fileName := fmt.Sprintf("%s-%s.json", getterName, filterName)
			if OptionsInstance.OutputWithDateTime {
				fileName = fmt.Sprintf("%s-%s-%s.json", getterName, filterName, nowTimeStr)
			}
			if err = utils.WriteToJsonFile(path.Join(OptionsInstance.OutputDir, fileName), news); err != nil {
				fmt.Println(err)
			}

			for _, hooker := range sender {
				hooker.Hook(filterName, news)
			}
		}

		fileName := fmt.Sprintf("%s-%s.json", getterName, "all")
		if OptionsInstance.OutputWithDateTime {
			fileName = fmt.Sprintf("%s-%s-%s.json", getterName, "all", nowTimeStr)
		}
		if err = utils.WriteToJsonFile(path.Join(OptionsInstance.OutputDir, fileName), res); err != nil {
			fmt.Println(err)
		}
	}
}
