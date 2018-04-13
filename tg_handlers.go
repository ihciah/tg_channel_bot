package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"./fetchers"
)

func (TGBOT *TelegramBot) RegisterHandler(){
	TGBOT.bot.Handle("hello", TGBOT.handle_helloworld)
	TGBOT.bot.Handle("example", TGBOT.handle_image_fetcher_example)
	TGBOT.bot.Handle("v2ex", TGBOT.handle_v2ex)
}

func (TGBOT *TelegramBot)handle_helloworld(m *tb.Message){
	TGBOT.bot.Send(m.Sender, "Hello world!")
}

func (TGBOT *TelegramBot)handle_image_fetcher_example(m *tb.Message){
	var fetcher fetchers.Fetcher
	fetcher = &fetchers.ExampleFetcher{}
	fetcher.Init()
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}

func (TGBOT *TelegramBot)handle_v2ex(m *tb.Message){
	var fetcher fetchers.Fetcher
	fetcher = &fetchers.V2EXFetcher{}
	fetcher.Init()
	msg := fetcher.Get()
	TGBOT.Send(m.Sender, msg)
}
