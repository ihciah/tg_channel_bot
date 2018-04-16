package main

import (
	f "./fetchers"
	tb "github.com/ihciah/telebot"
	"log"
	"strconv"
	"strings"
)

type FetcherConfig struct {
	Twitter f.TwitterFetcher `json:"twitter"`
}

func (TGBOT *TelegramBot) RegisterHandler() {
	TGBOT.Bot.Handle("/about", TGBOT.handle_about)
	TGBOT.Bot.Handle("/example", TGBOT.handle_example_fetcher_example)
	TGBOT.Bot.Handle("/v2ex", TGBOT.handle_v2ex)
	TGBOT.Bot.Handle("/twitter", TGBOT.handle_twitter)
	TGBOT.Bot.Handle("/tmedia", TGBOT.handle_twitter_media_test_only) // For test only
	TGBOT.Bot.Handle("/twitter_channel", TGBOT.handle_twitter_channel_test_only) // For test only
	TGBOT.Bot.Handle(tb.OnText, TGBOT.handle_controller)
}

func (TGBOT *TelegramBot) handle_about(m *tb.Message) {
	TGBOT.Bot.Send(m.Sender, "This is a Bot designed for syncing message(text/image/video) "+
		"from given sites to telegram channel by @ihciah.\n"+
		"Check https://github.com/ihciah/tg_channel_bot for source code and other information.")
}

func (TGBOT *TelegramBot) handle_example_fetcher_example(m *tb.Message) {
	var fetcher f.Fetcher = new(f.ExampleFetcher)
	fetcher.Init(TGBOT.Database)
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{}))
}

func (TGBOT *TelegramBot) handle_v2ex(m *tb.Message) {
	var fetcher f.Fetcher = new(f.V2EXFetcher)
	fetcher.Init(TGBOT.Database)
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{}))
}

func (TGBOT *TelegramBot) handle_twitter(m *tb.Message) {
	var fetcher f.Fetcher = &TGBOT.FetcherConfigs.Twitter
	fetcher.Init(TGBOT.Database)
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{}))
}

func (TGBOT *TelegramBot) handle_twitter_media_test_only(m *tb.Message) {
	var fetcher f.Fetcher = &TGBOT.FetcherConfigs.Twitter
	fetcher.Init(TGBOT.Database)
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{"cchaai", "MisaCat33"}))
}

func (TGBOT *TelegramBot) handle_twitter_channel_test_only(_ *tb.Message) {
	var fetcher f.Fetcher = &TGBOT.FetcherConfigs.Twitter
	channel := "@FEWSAWD"
	fetcher.Init(TGBOT.Database)
	chat, err := TGBOT.Bot.ChatByID(channel)
	if err != nil{
		log.Println("Error when start chat.", err)
	}
	TGBOT.SendAll(chat, fetcher.GetPush(channel, []string{"cchaai", "MisaCat33", "fleia0124"}))
}

func (TGBOT *TelegramBot) handle_controller(m *tb.Message) {
	handlers := map[string]func(string, *tb.Message)string{
		"addchannel": TGBOT.h_addchannel,
		"addfollow": TGBOT.h_addfollow,
		"delfollow": TGBOT.h_delfollow,
		"listfollow": TGBOT.h_listfollow,
		"setinterval": TGBOT.h_setinterval,
		"listchannel": TGBOT.h_listchannel,
	}
	available_commands := make([]string, 0, len(handlers))
	for c := range handlers{
		available_commands = append(available_commands, c)
	}

	var cmd, params string
	commands := strings.SplitN(m.Text, " ", 2)
	_, command_in := handlers[m.Text]
	if command_in{
		cmd, params = m.Text, ""
	}else if len(commands) == 2{
		cmd, params = commands[0], commands[1]
	}else{
		TGBOT.Send(m.Sender, f.ReplyMessage{Caption: "Command Format: CMD (params)"})
	}

	h_func, ok := handlers[cmd]
	if !ok{
		reply := "Unrecognized command.\nAvailable commands: \n" + strings.Join(available_commands, "\n")
		TGBOT.Send(m.Sender, f.ReplyMessage{Caption: reply})
		return
	}
	TGBOT.Send(m.Sender, f.ReplyMessage{Caption: h_func(params, m)})
}
