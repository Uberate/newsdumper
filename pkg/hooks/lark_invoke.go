package hooks

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"news/pkg/getter"
	"text/template"
)

const LarkHookKind = "lark"

var defaultLarkTemplate = `{
 "msg_type": "interactive",
 "card": {
   "elements": [
     {{ range $index, $element := .News }}
	   {{- if not (eq $index 0) -}} , {{- end -}}
       {{"tag": "hr"}},
       {"tag": "div", "text": {"content": "**{{$element.Title}}**", "tag": "lark_md" }},
       {"tag": "action", "actions": [{"tag": "button", "text": {"tag": "plain_text","content": "查看详情"}, "url": "{{$element.Link}}", "type": "primary"}]}
     {{ end }}
 ]},
 "header": {"title": {"content": "news: {{.Type}}", "tag": "plain_text"}}
}`

func GeneratorLarkHookInstance(name string, config interface{}, logger *logrus.Logger) (Hook, error) {
	o := &LarkHook{}
	if err := mapstructure.Decode(config, o); err != nil {
		return nil, err
	}

	o.logger = logger
	o.name = name

	if len(o.TemplateValue) == 0 {
		o.TemplateValue = defaultLarkTemplate
	}
	var err error
	o.innerTemplate, err = template.New("lark-template").Parse(o.TemplateValue)
	if err != nil {
		return nil, err
	}
	if len(o.Host) == 0 {
		return nil, fmt.Errorf("lark invoke host can't be empty")
	}

	return o, nil
}

type LarkHook struct {
	logger *logrus.Logger
	name   string

	Host          string `json:"host,omitempty" yaml:"host"`
	TemplateValue string `json:"templateValue,omitempty" yaml:"templateValue"`

	innerTemplate *template.Template
}

func (l LarkHook) Kind() string {
	return LarkHookKind
}

func (l LarkHook) Name() string {
	return l.name
}

func (l LarkHook) Version() string {
	return V1Str
}

func (l LarkHook) Hook(typ string, news []getter.News) error {

	type innerStruct struct {
		Type string
		News []getter.News
	}

	data := innerStruct{
		Type: typ,
		News: news,
	}
	output := bytes.NewBuffer([]byte{})
	if err := l.innerTemplate.Execute(output, data); err != nil {
		return err
	}

	resp, err := http.DefaultClient.Post(l.Host, "application/json", output)

	if err != nil {
		return err
	}

	if !(resp.StatusCode >= 100 && resp.StatusCode < 300) {
		res, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		l.logger.Warnf("response err: %v", string(res))
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	l.logger.Warnf("response err: %v", string(res))
	return nil
}
