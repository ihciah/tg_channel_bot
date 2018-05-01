package main

import (
	"fmt"
	tb "github.com/ihciah/telebot"
	"log"
	"strconv"
	"strings"
)

func (TGBOT *TelegramBot) h_addchannel(p []string, m *tb.Message) string {
	if len(p) != 1 {
		return "Usage: addchannel @channel_id/chat_id"
	}
	c, err := AddChannelIfNotExists(TGBOT, p[0])
	if err != nil {
		return fmt.Sprintf("Channel/Chat %s cannot be added. %s", p, err)
	}
	*TGBOT.Channels = append(*TGBOT.Channels, c)
	go c.Push()
	return fmt.Sprintf("Channel/Chat %s added.", p)
}

func (TGBOT *TelegramBot) h_delchannel(p []string, m *tb.Message) string {
	if len(p) != 1 {
		return "Usage: delchannel @channel_id/chat_id"
	}
	if err := DelChannelIfExists(TGBOT, p[0]); err != nil {
		return fmt.Sprintf("Channel/Chat %s cannot be deleted.", p)
	}
	for i, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			go v.Exit()
			*TGBOT.Channels = append((*TGBOT.Channels)[:i], (*TGBOT.Channels)[i+1:]...)
			break
		}
	}
	return fmt.Sprintf("Channel/Chat %s deleted.", p)
}

func (TGBOT *TelegramBot) h_addfollow(p []string, m *tb.Message) string {
	return TGBOT.h_user(p, m, true)
}

func (TGBOT *TelegramBot) h_delfollow(p []string, m *tb.Message) string {
	return TGBOT.h_user(p, m, false)
}

func (TGBOT *TelegramBot) h_user(p []string, m *tb.Message, is_add bool) string {
	if len(p) != 3 {
		return "Usage: addfollow/delfollow @channel_id/chat_id site(twitter/tumblr) userid"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if !auth_user(m.Sender, *v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module := MakeModuleLabeler().Str2Module(p[1])
			if module == -1 {
				return "Unsupported site."
			}
			if is_add {
				v.AddFollowing(ModuleUser{module, p[2]})
				return "Following added."
			} else {
				v.DelFollowing(ModuleUser{module, p[2]})
				return "Following deleted."
			}
		}
	}
	return "No such channel."
}

func (TGBOT *TelegramBot) h_addadmin(p []string, m *tb.Message) string {
	return TGBOT.h_admin(p, m, true)
}

func (TGBOT *TelegramBot) h_deladmin(p []string, m *tb.Message) string {
	return TGBOT.h_admin(p, m, false)
}

func (TGBOT *TelegramBot) h_listadmin(p []string, m *tb.Message) string {
	if len(p) != 1 {
		return "Usage: listadmin @channel_id/chat_id"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if len(*v.AdminUserIDs) == 0 {
				return "No admin."
			}
			return fmt.Sprintf("Admins for %s:\n%s", v.ID, strings.Join(*v.AdminUserIDs, "\n"))
		}
	}
	return "No such channel."
}

func (TGBOT *TelegramBot) h_admin(p []string, m *tb.Message, is_add bool) string {
	if len(p) != 2 {
		return "Usage: addadmin/deladmin @channel_id/chat_id userid"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if is_add {
				v.AddAdmin(p[1])
				return "Admin added."
			} else {
				v.DelAdmin(p[1])
				return "Admin deleted."
			}
		}
	}
	return "No such channel/chat."
}

func (TGBOT *TelegramBot) h_listfollow(p []string, m *tb.Message) string {
	if len(p) != 1 {
		return "Usage: listfollow @channel_id/chat_id"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if !auth_user(m.Sender, *v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			ret := make([]string, 0, len(*TGBOT.Channels))
			for module_id, names := range *v.Followings {
				ret = append(ret, fmt.Sprintf("Module: %s\nUpdateInterval: %d\nFollowings:\n%s", MakeModuleLabeler().Module2Str(module_id), (*v.PushIntervals)[module_id], strings.Join(names, "\n")))
			}
			if len(ret) == 0 {
				return "No followings."
			}
			return strings.Join(ret, "\n\n")
		}
	}
	return "No such channel"
}

func (TGBOT *TelegramBot) h_setinterval(p []string, m *tb.Message) string {
	if len(p) != 3 {
		return "Usage: setinterval @channel_id/chat_id site N(second)"
	}
	interval, err := strconv.Atoi(p[2])
	if err != nil || interval <= 0 {
		return "Usage: setinterval @Channel/chat_id site N(second), N should be a positive number"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if !auth_user(m.Sender, *v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module_id := MakeModuleLabeler().Str2Module(p[1])
			if module_id < 0 {
				return "Unsupported site."
			}
			v.UpdateInterval(ModuleInterval{module_id, interval})
			return "Push Interval Updated."
		}
	}
	return "No such channel/chat"
}

func (TGBOT *TelegramBot) h_listchannel(p []string, m *tb.Message) string {
	names := make([]string, 0, len(*TGBOT.Channels))
	for _, v := range *TGBOT.Channels {
		names = append(names, v.ID)
	}
	return "Channels:\n" + strings.Join(names, "\n")
}

func (TGBOT *TelegramBot) h_goback(p []string, m *tb.Message) string {
	if len(p) != 3 {
		return "Usage: goback @channel_id/chat_id site N(second), N=0 means reset to Now."
	}
	back, err := strconv.ParseInt(p[2], 10, 64)
	if err != nil || back < 0 {
		return "Usage: goback @Channel/chat_id site N(second), N >= 0"
	}
	for _, v := range *TGBOT.Channels {
		if v.ID == p[0] {
			if !auth_user(m.Sender, *v.AdminUserIDs, TGBOT.Admins) {
				return "Unauthorized."
			}
			module_id := MakeModuleLabeler().Str2Module(p[1])
			if module_id < 0 {
				return "Unsupported site."
			}
			fetcher := TGBOT.CreateModule(module_id, v.ID)
			if err := fetcher.GoBack(v.ID, back); err != nil {
				return fmt.Sprintf("Error when go back. %s", err)
			}
			return fmt.Sprintf("Site %s for channel/chat %s has been set to %d seconds before.", p[1], v.ID, back)
		}
	}
	return "No such channel/chat"
}

func (TGBOT *TelegramBot) h_getid(p []string, m *tb.Message) string {
	chat_id := m.Chat.ID
	chat_title := m.Chat.Title
	user_id := m.Sender.ID
	first_name := m.Sender.FirstName
	last_name := m.Sender.LastName
	username := m.Sender.Username
	return fmt.Sprintf("Hi %s %s(%s) !\nYour ID: %d\n\nChat: %s\nChatID: %d", last_name, first_name, username, user_id, chat_title, chat_id)
}

func (TGBOT *TelegramBot) requireSuperAdmin(f func([]string, *tb.Message) string) func([]string, *tb.Message) string {
	return func(p []string, m *tb.Message) string {
		if auth_user(m.Sender, []string{}, TGBOT.Admins) {
			log.Println("Authorized.")
			return f(p, m)
		}
		log.Println("Unauthorized", m.Sender.Username)
		return "Unauthorized user. Superadmin needed."
	}
}

func auth_user(user *tb.User, admin_list []string, super_admin_list []string) bool {
	for _, u := range admin_list {
		if user.Username == u {
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
