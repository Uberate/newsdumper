package hooks

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"news/pkg/getter"
)

const SMTPHookKind = "smtp"

func GeneratorSMTPHook(name string, config interface{}, logger *logrus.Logger) (Hook, error) {
	o := &SMTPHook{}
	if err := mapstructure.Decode(config, o); err != nil {
		return nil, err
	}

	// check
	if len(o.Host) == 0 {
		return nil, fmt.Errorf("SMTP config need host param")
	}

	o.name = name
	o.logger = logger
	return o, nil
}

type SMTPHook struct {
	logger *logrus.Logger
	name   string

	Host      string   `json:"host" yaml:"host"`
	Port      string   `json:"port" yaml:"port"`
	UserName  string   `json:"username" yaml:"userName"`
	Receivers []string `json:"receivers" yaml:"receivers"`
}

func (h *SMTPHook) Kind() string {
	return SMTPHookKind
}

func (h *SMTPHook) Version() string {
	return V1Str
}

func (h *SMTPHook) Name() string {
	return h.name
}

func (h *SMTPHook) Hook(typ string, news []getter.News) error {
	msgs := fmt.Sprintf("Subject:%s:%s\r\n\r\n", "News", typ)
	if len(news) == 0 {
		return nil
	}
	for index, item := range news {
		msgs += fmt.Sprintf("\r\n<h4>%d: <a href=\"%s\">%s</a></h4>\r\n", index+1, item.Link, item.Title)
	}
	err := smtp.SendMail(fmt.Sprintf("%s:%s", h.Host, h.Port), nil, h.UserName, h.Receivers, []byte(msgs))
	return err
}
