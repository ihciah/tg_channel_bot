package main

import (
	"github.com/ihciah/tg_channel_bot/fetchers"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (TGBOT *TelegramBot) RegisterHandler() {
	TGBOT.Bot.Handle("/about", TGBOT.handle_about)
	TGBOT.Bot.Handle("/example", TGBOT.handle_image_fetcher_example)
	TGBOT.Bot.Handle("/v2ex", TGBOT.handle_v2ex)
}

func (TGBOT *TelegramBot) handle_about(m *tb.Message) {
	TGBOT.Bot.Send(m.Sender, "This is a Bot designed for syncing message(text/image/video) " +
		"from given sites to telegram channel by @ihciah.\n"+
		"Check https://github.com/ihciah/tg_channel_bot for source code and other information.")
}

func (TGBOT *TelegramBot) handle_image_fetcher_example(m *tb.Message) {
	var fetcher fetchers.Fetcher
	fetcher = &fetchers.ExampleFetcher{}
	fetcher.Init()
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}

func (TGBOT *TelegramBot) handle_v2ex(m *tb.Message) {
	var fetcher fetchers.Fetcher
	fetcher = &fetchers.V2EXFetcher{}
	fetcher.Init()
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}
