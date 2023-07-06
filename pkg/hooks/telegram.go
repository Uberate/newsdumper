package hooks

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"news/pkg/getter"
)

const TelegramChannelHookKind = "telegram-channel"

//func GeneratorLarkHookInstance(name string, config interface{}, logger *logrus.Logger) (Hook, error) {

func GeneratorTelegramChannelHookInstance(name string, config interface{}, logger *logrus.Logger) (Hook, error) {
	o := &TelegramChannelHook{}

	o.name = name
	o.logger = logger

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "yaml",
		Result:  o,
	})
	if err != nil {
		return nil, err
	}

	if err := d.Decode(config); err != nil {
		return nil, err
	}
	return o, nil
}

type TelegramChannelHook struct {
	logger *logrus.Logger
	name   string

	BotToken  string `json:"bot_token" yaml:"bot_token"`
	ChannelId int64  `json:"channel_id" yaml:"channel_id"`
}

func (tch *TelegramChannelHook) Kind() string {
	return TelegramChannelHookKind
}

func (tch *TelegramChannelHook) Name() string {
	return tch.name
}

func (tch *TelegramChannelHook) Version() string {
	return V1Str
}

func (tch *TelegramChannelHook) Hook(typ string, news []getter.News) error {
	bot, err := tgbotapi.NewBotAPI(tch.BotToken)
	if err != nil {
		return err
	}

	messages := typ + ": \n"
	for index, item := range news {
		messages += fmt.Sprintf("%d. %s \nLink: %s\n\n", index+1, item.Title, item.Link)

	}

	message := tgbotapi.NewMessage(tch.ChannelId, messages)
	_, err = bot.Send(message)

	return err
}
