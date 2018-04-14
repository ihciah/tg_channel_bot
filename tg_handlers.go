package main

import (
	f "github.com/ihciah/tg_channel_bot/fetchers"
	tb "gopkg.in/tucnak/telebot.v2"
)

type FetcherConfig struct {
	Twitter f.TwitterFetcher `json:"twitter"`
}

func (TGBOT *TelegramBot) RegisterHandler() {
	TGBOT.Bot.Handle("/about", TGBOT.handle_about)
	TGBOT.Bot.Handle("/example", TGBOT.handle_example_fetcher_example)
	TGBOT.Bot.Handle("/v2ex", TGBOT.handle_v2ex)
	TGBOT.Bot.Handle("/twitter", TGBOT.handle_twitter)
}

func (TGBOT *TelegramBot) handle_about(m *tb.Message) {
	TGBOT.Bot.Send(m.Sender, "This is a Bot designed for syncing message(text/image/video) "+
		"from given sites to telegram channel by @ihciah.\n"+
		"Check https://github.com/ihciah/tg_channel_bot for source code and other information.")
}

func (TGBOT *TelegramBot) handle_example_fetcher_example(m *tb.Message) {
	var fetcher f.Fetcher = new(f.ExampleFetcher)
	fetcher.Init(TGBOT.Database)
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}

func (TGBOT *TelegramBot) handle_v2ex(m *tb.Message) {
	var fetcher f.Fetcher = new(f.V2EXFetcher)
	fetcher.Init(TGBOT.Database)
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}

func (TGBOT *TelegramBot) handle_twitter(m *tb.Message) {
	var fetcher f.Fetcher = &TGBOT.FetcherConfigs.Twitter
	fetcher.Init(TGBOT.Database)
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}
