package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
)

func (TGBOT *TelegramBot) RegisterHandler(){
	TGBOT.bot.Handle("hello", TGBOT.handle_helloworld)
	TGBOT.bot.Handle("example", TGBOT.handle_image_fetcher_example)
}

func (TGBOT *TelegramBot)handle_helloworld(m *tb.Message){
	TGBOT.bot.Send(m.Sender, "Hello world!")
}

func (TGBOT *TelegramBot)handle_image_fetcher_example(m *tb.Message){
	var fetcher Fetcher
	fetcher = &ExampleFetcher{}
	fetcher.Init()
	log.Println("Fetcher init.")
	msg := fetcher.Get()
	log.Println("Image url get.")
	TGBOT.Send(m.Sender, msg)
}
