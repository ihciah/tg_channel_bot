package main

import (
	tb "github.com/ihciah/telebot"
	f "github.com/ihciah/tg_channel_bot/fetchers"
	"strconv"
	"strings"
)

const AboutMessage = "This is a Bot designed for syncing message(text/image/video) " +
	"from given sites to telegram channel/user/group by @ihciah.\n" +
	"Check https://github.com/ihciah/tg_channel_bot for source code and other information.\n"

type FetcherConfig struct {
	Base    f.BaseFetcher
	Twitter f.TwitterFetcher `json:"twitter"`
	Tumblr  f.TumblrFetcher  `json:"tumblr"`
	V2EX    f.V2EXFetcher
}

func (TGBOT *TelegramBot) CreateModule(module_id int, channel_id string) f.Fetcher {
	var fetcher f.Fetcher
	switch module_id {
	case MTwitter:
		fetcher = &TGBOT.FetcherConfigs.Twitter
	case MTumblr:
		fetcher = &TGBOT.FetcherConfigs.Tumblr
	case MV2EX:
		fetcher = &TGBOT.FetcherConfigs.V2EX
	default:
		fetcher = &TGBOT.FetcherConfigs.Base
	}
	fetcher.Init(TGBOT.Database, channel_id)
	return fetcher
}

func (TGBOT *TelegramBot) RegisterHandler() {
	TGBOT.Bot.Handle("/about", TGBOT.handle_about)
	TGBOT.Bot.Handle("/id", TGBOT.handle_id)
	//TGBOT.Bot.Handle("/example", TGBOT.handle_example_fetcher_example)
	//TGBOT.Bot.Handle("/v2ex", TGBOT.handle_v2ex)
	TGBOT.Bot.Handle(tb.OnText, TGBOT.handle_controller)
	TGBOT.Bot.Handle(tb.OnPhoto, TGBOT.handle_photo)
}

func (TGBOT *TelegramBot) handle_photo(m *tb.Message) {
	chatid := strconv.FormatInt(m.OriginalChat.ID, 10)
	if m.OriginalChat.Type == "channel" {
		chatid = "@" + m.OriginalChat.Username
	}

	pass := false
	for _, v := range *TGBOT.Channels {
		if v.ID == chatid && auth_user(m.Sender, *v.AdminUserIDs, TGBOT.Admins) {
			pass = true
			break
		}
	}
	if !pass {
		TGBOT.Bot.Send(m.Sender, "Unauthorized.")
		return
	}

	var fetcher f.Fetcher
	if strings.Contains(m.Caption, "tumblr") {
		fetcher = new(f.TumblrFetcher)
	}
	fetcher.Init(TGBOT.Database, chatid)
	TGBOT.Bot.Send(m.Sender, fetcher.Block(m.Caption))
}

func (TGBOT *TelegramBot) handle_about(m *tb.Message) {
	TGBOT.Bot.Send(m.Sender, AboutMessage)
}

func (TGBOT *TelegramBot) handle_id(m *tb.Message) {
	TGBOT.Bot.Send(m.Chat, TGBOT.h_getid([]string{}, m))
}

func (TGBOT *TelegramBot) handle_example_fetcher_example(m *tb.Message) {
	var fetcher f.Fetcher = new(f.ExampleFetcher)
	fetcher.Init(TGBOT.Database, "")
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{}))
}

func (TGBOT *TelegramBot) handle_v2ex(m *tb.Message) {
	var fetcher f.Fetcher = new(f.V2EXFetcher)
	fetcher.Init(TGBOT.Database, "")
	TGBOT.SendAll(m.Sender, fetcher.GetPushAtLeastOne(strconv.Itoa(m.Sender.ID), []string{}))
}

func (TGBOT *TelegramBot) handle_controller(m *tb.Message) {
	handlers := map[string]func([]string, *tb.Message) string{
		"addchannel":  TGBOT.requireSuperAdmin(TGBOT.h_addchannel),
		"delchannel":  TGBOT.requireSuperAdmin(TGBOT.h_delchannel),
		"listchannel": TGBOT.requireSuperAdmin(TGBOT.h_listchannel),
		"addfollow":   TGBOT.h_addfollow,
		"delfollow":   TGBOT.h_delfollow,
		"listfollow":  TGBOT.h_listfollow,
		"addadmin":    TGBOT.requireSuperAdmin(TGBOT.h_addadmin),
		"deladmin":    TGBOT.requireSuperAdmin(TGBOT.h_deladmin),
		"listadmin":   TGBOT.requireSuperAdmin(TGBOT.h_listadmin),
		"setinterval": TGBOT.h_setinterval,
		"goback":      TGBOT.h_goback,
		"id":          TGBOT.h_getid,
	}

	var cmd string
	var params []string
	commands := strings.Fields(m.Text)
	if _, command_in := handlers[commands[0]]; command_in {
		cmd, params = commands[0], commands[1:]
		TGBOT.Send(m.Sender, f.ReplyMessage{Caption: handlers[cmd](params, m)})
	} else {
		available_commands := make([]string, 0, len(handlers))
		for c := range handlers {
			available_commands = append(available_commands, c)
		}
		reply := AboutMessage + "\nAlso, you can send /id to any chat to get chatid." + "\n\nUnrecognized command.\nAvailable commands: \n" + strings.Join(available_commands, "\n")
		TGBOT.Send(m.Sender, f.ReplyMessage{Caption: reply})
	}
}
