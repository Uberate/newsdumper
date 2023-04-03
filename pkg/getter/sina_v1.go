package getter

import (
	"encoding/json"
	"io"
	"net/http"
)

type SinaV1NewStruct struct {
	Items []SianV1Items `json:"items"`
}

type SianV1Items struct {
	Title        string `json:"title"`
	DateModified string `json:"date_modified"`
	Url          string `json:"url"`
}

func SinaV1ResParser(response *http.Response) ([]News, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return make([]News, 0, 0), err
	}

	news := &SinaV1NewStruct{}
	if err = json.Unmarshal(body, news); err != nil {
		return make([]News, 0, 0), err
	}

	res := make([]News, 0, len(news.Items))

	for _, item := range news.Items {
		res = append(res, News{
			Title: item.Title,
			Time:  item.DateModified,
			Link:  item.Url,
		})
	}

	return res, nil
}
