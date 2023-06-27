package getter

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"news/pkg/factory"
	"strings"
)

type NewsGetter interface {
	Name() string
	Kind() string
	Version() string
	GetNews(callTime int64) ([]News, error)
}

type EmptyGetter struct {
}

func (e EmptyGetter) Name() string {
	return ""
}

func (e EmptyGetter) Kind() string {
	return ""
}

func (e EmptyGetter) Version() string {
	return ""
}

func (e EmptyGetter) GetNews(callTime int64) ([]News, error) {
	return nil, nil
}

//----------------------------------------------------------------------------------------------------------------------

func StableHeader(header map[string]string) func(int64) map[string]string {
	if header == nil {
		header = map[string]string{}
	}
	return func(int64) map[string]string {
		return header
	}
}

func StableParams(params map[string]string) func(int64) map[string]string {
	if params == nil {
		params = map[string]string{}
	}
	return func(int64) map[string]string {
		return params
	}
}

func NewAbsNewsGetter(
	kind, version string, // kind and version
	link, method string, // link and method
	bodyGen func(int642 int64) string, // body gen func
	header func(int642 int64) map[string]string, // header gen func
	paramGen func(int642 int64) map[string]string, // param gen func
	resParser func(response *http.Response) ([]News, error), // parser func
) factory.Generator[NewsGetter] {
	return func(name string, config interface{}, logger *logrus.Logger) NewsGetter {
		return &AbsNewsGetter{
			getterName: name,
			kind:       kind,
			version:    version,
			logger:     logger,

			Link:      link,
			HeaderGen: header,
			Method:    method,
			BodyGen:   bodyGen,
			ParamGen:  paramGen,

			resParser: resParser,
		}
	}
}

func SimpleNewsGetter(
	kind, version string,
	link string,
	headers map[string]string,
	resParse func(response *http.Response) ([]News, error)) factory.Generator[NewsGetter] {
	return NewAbsNewsGetter(
		kind, version,
		link, http.MethodGet,
		nil,
		StableHeader(headers),
		StableParams(map[string]string{}),
		resParse)
}

type AbsNewsGetter struct {
	version    string
	kind       string
	getterName string

	logger *logrus.Logger

	Link      string
	HeaderGen func(endTime int64) map[string]string
	Method    string
	ParamGen  func(endTime int64) map[string]string
	BodyGen   func(endTime int64) string

	resParser func(response *http.Response) ([]News, error)
}

func (getter *AbsNewsGetter) Name() string {
	return getter.getterName
}

func (getter *AbsNewsGetter) Kind() string {
	return getter.kind
}

func (getter *AbsNewsGetter) Version() string {
	return getter.version
}

func (getter *AbsNewsGetter) GetNews(endTime int64) ([]News, error) {
	header := map[string]string{}

	if getter.HeaderGen != nil {
		header = getter.HeaderGen(endTime)
	}

	params := map[string]string{}
	if getter.ParamGen != nil {
		params = getter.ParamGen(endTime)
	}

	body := ""
	if getter.BodyGen != nil {
		body = getter.BodyGen(endTime)
	}

	bodyReader := strings.NewReader(body)

	request, err := http.NewRequest(getter.Method, getter.Link, bodyReader)
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

	getter.logger.Tracef("[time: %d, name: %s] do http request getter\n", endTime, getter.Name())
	getter.logger.Debugf("[time: %d, name: %s] meta: Kind: %s, Version: %s, Name: %s\n", endTime, getter.Name(), getter.kind, getter.version, getter.Name())
	getter.logger.Infof("[time: %d, name: %s] do http request, method: [%s], request: [%s]", endTime, getter.Name(), getter.Method, getter.Link)
	getter.logger.Debugf("[time: %d, name: %s] request headers: %v", endTime, getter.Name(), header)
	getter.logger.Debugf("[time: %d, name: %s] request params: %v", endTime, getter.Name(), params)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return make([]News, 0, 0), err
	}

	if getter.resParser == nil {
		return make([]News, 0, 0), nil
	}

	return getter.resParser(response)
}
