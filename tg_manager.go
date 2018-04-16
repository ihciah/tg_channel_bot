package main

import (
	tb "github.com/ihciah/telebot"
	"strconv"
	"strings"
	"fmt"
)

func (TGBOT *TelegramBot) h_addchannel(p string, m *tb.Message) string {
	// Admin check
	c := AddChannelIfNotExists(TGBOT, p)
	*TGBOT.Channels = append(*TGBOT.Channels, c)
	go c.Push()
	return fmt.Sprintf("Channel %s added.", p)
}

func (TGBOT *TelegramBot) h_addfollow(p string, m *tb.Message) string {
	return TGBOT.h_user(p, m, true)
}

func (TGBOT *TelegramBot) h_delfollow(p string, m *tb.Message) string {
	return TGBOT.h_user(p, m, false)
}

func (TGBOT *TelegramBot) h_user(p string, m *tb.Message, is_add bool) string {
	// Admin check
	commands := strings.SplitN(p, " ", 3)
	if len(commands) != 3{
		return "Usage: addfollow/delfollow site channel_id userid"
	}
	for _, v := range *TGBOT.Channels{
		if v.ID == commands[1]{
			module := Str2Module(commands[0])
			if module == -1{
				return "Unsupported site."
			}
			if is_add {
				v.AddFollowing(ModuleUser{module, commands[2]})
				return "Following added."
			}else{
				v.DelFollowing(ModuleUser{module, commands[2]})
				return "Following deleted."
			}
		}
	}
	return "No such channel. Add it first."
}

func (TGBOT *TelegramBot) h_listfollow(p string, m *tb.Message) string {
	for _, v := range *TGBOT.Channels{
		if v.ID == p{
			ret := make([]string, 0, len(*TGBOT.Channels))
			for module_id, names := range *v.Followings{
				ret = append(ret, fmt.Sprintf("Module: %s\nUpdateInterval: %d\nFollowings:\n%s", Module2Str(module_id), (*v.PushIntervals)[module_id], strings.Join(names, "\n")))
			}
			if len(ret) == 0{
				return "No followings."
			}
			return strings.Join(ret, "\n")
		}
	}
	return "No such channel"
}

func (TGBOT *TelegramBot) h_setinterval(p string, m *tb.Message) string {
	commands := strings.SplitN(p, " ", 3)
	if len(commands) != 3{
		return "Usage: setinterval channel_id site N(second)"
	}
	interval, err := strconv.Atoi(commands[2])
	if err != nil || interval <= 0{
		return "Usage: setinterval channel_id site N(second), N should be a positive number"
	}
	for _, v := range *TGBOT.Channels{
		if v.ID == commands[0]{
			module_id := Str2Module(commands[1])
			if module_id < 0{
				return "Unsupported site."
			}
			(*v.PushIntervals)[module_id] = interval
			return "Push Interval Updated."
		}
	}
	return "No such channel"
}

func (TGBOT *TelegramBot) h_listchannel(p string, m *tb.Message) string {
	names := make([]string, 0, len(*TGBOT.Channels))
	for _, v := range *TGBOT.Channels{
		names = append(names, v.ID)
	}
	return "Channels:\n" + strings.Join(names, "\n")
}
