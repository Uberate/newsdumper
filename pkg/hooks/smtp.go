package hooks

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/smtp"
	"news/pkg/getter"
)

const SMTPHookV1 = "smtp-v1"

func InitSMTPHook(config interface{}) (Hooks, error) {
	o := &SMTPHook{}
	err := mapstructure.Decode(config, o)
	return o, err
}

type SMTPHook struct {
	Host      string   `json:"host"`
	Port      string   `json:"port"`
	UserName  string   `json:"username"`
	Receivers []string `json:"receivers"`
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
