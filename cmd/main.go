package main

import (
	"flag"
	"fmt"
	"news/pkg/getter"
	"news/pkg/utils"
	"path"
	"time"
)

func main() {
	nowTime := time.Now()
	nowTimeSecond := nowTime.Unix()
	nowTimeStr := nowTime.Format(time.RFC3339)

	writePath := flag.String("dump-path", "./output", "write path of output")
	flag.Parse()

	for name, item := range getter.NewsGetters {
		res, err := item.GetNews(nowTimeSecond)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if err = utils.WriteToJsonFile(path.Join(*writePath, name+"_"+nowTimeStr+".json"), res); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
