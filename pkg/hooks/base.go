package hooks

import (
	"fmt"
	"news/pkg/getter"
)

type InitHook func(config interface{}) (Hook, error)

type Hook interface {
	Hook(typ string, news []getter.News) error
}

func NewAbsHook(invokeHook func(typ string, news []getter.News) error) *AbsHooks {
	return &AbsHooks{
		SendFunction: invokeHook,
	}
}

type AbsHooks struct {
	SendFunction func(typ string, news []getter.News) error
}

func (ah *AbsHooks) Hook(typ string, news []getter.News) error {
	return ah.SendFunction(typ, news)
}

var Hookers map[string]InitHook

func init() {
	Hookers = map[string]InitHook{}
	Hookers[SMTPHookV1] = InitSMTPHook
}

func GetHook(typ string, config interface{}) (Hook, error) {
	if v, ok := Hookers[typ]; ok {
		res, err := v(config)
		return res, err
	}

	return nil, fmt.Errorf("not found type: %s", typ)
}
