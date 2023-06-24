package getter

import (
	"net/http"
	"strings"
)

type News struct {
	Getter string
	Title  string
	Body   string
	Link   string
	Time   string
	Label  []string
}

// Getter define a getter to get news.
type Getter interface {
	GetNews(callTime int64) ([]News, error)
	Name() string
	LastGet() int64
}

type AbsGetter struct {
	Link       string
	HeaderGen  func(endTime int64) map[string]string
	Method     string
	ParamGen   func(endTime int64) map[string]string
	BodyGen    func(endTime int64) string
	GetterName string

	resParser func(response *http.Response) ([]News, error)

	// lastGet is the last get news time.
	lastGet int64
}

func (ag *AbsGetter) Name() string {
	return ag.GetterName
}

func (ag *AbsGetter) GetNews(endTime int64) ([]News, error) {
	if endTime < ag.lastGet {
		return make([]News, 0, 0), nil
	}

	header := map[string]string{}
	if ag.HeaderGen != nil {
		header = ag.HeaderGen(endTime)
	}

	params := map[string]string{}
	if ag.ParamGen != nil {
		params = ag.ParamGen(endTime)
	}

	body := ""
	if ag.BodyGen != nil {
		body = ag.BodyGen(endTime)
	}

	bodyReader := strings.NewReader(body)

	request, err := http.NewRequest(ag.Method, ag.Link, bodyReader)
	if err != nil {
		return make([]News, 0, 0), err
	}

	q := request.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}

	request.URL.RawQuery = q.Encode()
	for k, v := range header {
		request.Header.Set(k, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return make([]News, 0, 0), err
	}

	if ag.resParser == nil {
		return make([]News, 0, 0), nil
	}

	return ag.resParser(response)
}

func (ag *AbsGetter) LastGet() int64 {
	return ag.lastGet
}

func NewAbsGetter(getterName, link, method string,
	bodyGen func(int642 int64) string,
	header,
	paramGen func(int642 int64) map[string]string,
	resParser func(response *http.Response) ([]News, error)) Getter {
	return &AbsGetter{
		Link:      link,
		HeaderGen: header,
		Method:    method,
		BodyGen:   bodyGen,
		ParamGen:  paramGen,
		resParser: resParser,

		lastGet: 0,
	}
}

func NewSimpleGetGetter(getterName, link string, headers map[string]string, resParse func(response *http.Response) ([]News, error)) Getter {
	return NewAbsGetter(getterName, link, http.MethodGet, nil, func(int642 int64) map[string]string {
		return headers
	}, nil, resParse)
}

// ==================================== getters

const SinaNewV1 = "sina_news_v1"

var NewsGetters = map[string]Getter{
	SinaNewV1: NewSimpleGetGetter(SinaNewV1, "https://sina-news.vercel.app/rss.json", nil, SinaV1ResParser),
}
