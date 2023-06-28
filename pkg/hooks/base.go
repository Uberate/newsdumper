package hooks

import (
	"news/pkg/factory"
	"news/pkg/getter"
)

const V1Str = "v1"

type Hook interface {
	factory.Entity
	Hook(typ string, news []getter.News) error
}

type EmptyHooker struct {
}

func (e EmptyHooker) Kind() string {
	return ""
}

func (e EmptyHooker) Name() string {
	return ""
}

func (e EmptyHooker) Version() string {
	return ""
}

func (e EmptyHooker) Hook(typ string, news []getter.News) error {
	return nil
}
