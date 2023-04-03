package main

import (
	"fmt"
	"news/pkg/getter"
	"time"
)

func main() {
	for _, item := range getter.NewsGetters {
		res, err := item.GetNews(time.Now().UnixNano())
		fmt.Println(res, err)
	}
}
