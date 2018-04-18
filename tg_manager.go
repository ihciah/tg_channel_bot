package main

import (
	"fmt"
	tb "github.com/ihciah/telebot"
	"log"
	"strconv"
	"strings"
)

func (TGBOT *TelegramBot) h_addchannel(p string, m *tb.Message) string {
	c, err := AddChannelIfNotExists(TGBOT, p)
	if err != nil {
		return fmt.Sprintf("Channel %s cannot be added. %s", p, err)
	}
	*TGBOT.Channels = append(*TGBOT.Channels, c)
	go c.Push()
	return fmt.Sprintf("Channel %s added.", p)
}

func (TGBOT *TelegramBot) h_delchannel(p string, m *tb.Message) string {
	if err := DelChannelIfExists(TGBOT, p); err != nil {
		return fmt.Sprintf("Channel %s cannot be deleted.", p)
	}
	for i, v := range *TGBOT.Channels {
		if v.ID == p {
			go v.Exit()
		}
		*TGBOT.Channels = append((*TGBOT.Channels)[:i], (*TGBOT.Channels)[i+1:]...)
	}
	return fmt.Sprintf("Channel %s deleted.", p)
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
	if len(commands) != 3 {
		return "Usage: addfollow/delfollow @Channel site userid"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == commands[0] {
			if !auth_user(m.Sender, v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module := MakeModuleLabeler().Str2Module(commands[1])
			if module == -1 {
				return "Unsupported site."
			}
			if is_add {
				v.AddFollowing(ModuleUser{module, commands[2]})
				return "Following added."
			} else {
				v.DelFollowing(ModuleUser{module, commands[2]})
				return "Following deleted."
			}
		}
	}
	return "No such channel. Add it first."
}

func (TGBOT *TelegramBot) h_listfollow(p string, m *tb.Message) string {
	if p == "" {
		return "Usage: listfollow @Channel"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p {
			if !auth_user(m.Sender, v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			ret := make([]string, 0, len(*TGBOT.Channels))
			for module_id, names := range *v.Followings {
				ret = append(ret, fmt.Sprintf("Module: %s\nUpdateInterval: %d\nFollowings:\n%s", MakeModuleLabeler().Module2Str(module_id), (*v.PushIntervals)[module_id], strings.Join(names, "\n")))
			}
			if len(ret) == 0 {
				return "No followings."
			}
			return strings.Join(ret, "\n")
		}
	}
	return "No such channel"
}

func (TGBOT *TelegramBot) h_setinterval(p string, m *tb.Message) string {
	commands := strings.SplitN(p, " ", 3)
	if len(commands) != 3 {
		return "Usage: setinterval @Channel site N(second)"
	}
	interval, err := strconv.Atoi(commands[2])
	if err != nil || interval <= 0 {
		return "Usage: setinterval @Channel site N(second), N should be a positive number"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == commands[0] {
			if !auth_user(m.Sender, v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module_id := MakeModuleLabeler().Str2Module(commands[1])
			if module_id < 0 {
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
	for _, v := range *TGBOT.Channels {
		names = append(names, v.ID)
	}
	return "Channels:\n" + strings.Join(names, "\n")
}

func (TGBOT *TelegramBot) h_goback(p string, m *tb.Message) string {
	commands := strings.SplitN(p, " ", 3)
	if len(commands) != 3 {
		return "Usage: goback @Channel site N(second), N=0 means reset to Now."
	}
	back, err := strconv.ParseInt(commands[2], 10, 64)
	if err != nil || back < 0 {
		return "Usage: goback @Channel site N(second), N >= 0"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == commands[0] {
			if !auth_user(m.Sender, v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module_id := MakeModuleLabeler().Str2Module(commands[1])
			if module_id < 0 {
				return "Unsupported site."
			}
			fetcher := TGBOT.CreateModule(module_id)
			if err := fetcher.GoBack(v.ID, back); err != nil {
				return fmt.Sprintf("Error when go back. %s", err)
			}
			return fmt.Sprintf("Site %s for channel %s has been set to %d seconds before.", commands[1], v.ID, back)
		}
	}
	return "No such channel"
}

func (TGBOT *TelegramBot) requireSuperAdmin(f func(string, *tb.Message) string) func(string, *tb.Message) string {
	return func(p string, m *tb.Message) string {
		if auth_user(m.Sender, []int{}, TGBOT.Admins) {
			log.Println("Authorized.")
			return f(p, m)
		}
		log.Println("Unauthorized", m.Sender.Username)
		return "Unauthorized user."
	}
}

func auth_user(user *tb.User, admin_list []int, super_admin_list []string) bool {
	for _, u := range admin_list {
		if user.ID == u {
			return true
		}
	}
	for _, su := range super_admin_list {
		if user.Username == su {
			return true
		}
	}
	return false
}
