package main

import (
	f "./fetchers"
	"github.com/asdine/storm"
	"github.com/ihciah/telebot"
	"log"
	"time"
	"strings"
)

const (
	ChannelActionEnable = iota
	ChannelActionDisable
	ChannelActionDelete
	ChannelActionAddAdmin
	ChannelActionDelAdmin
	ChannelActionAddFollow
	ChannelActionDelFollow
	ChannelActionUpdatePushInterval
)

const (
	MTwitter = iota
	MTumblr
	MV2EX
)

const (
	DefaultInterval = 30
)

const (
	SignalExit = iota
	SignalReload
)


type ModuleUser struct {
	Module   int
	Username string
}

type ModuleInterval struct {
	Module       int
	PushInterval int
}

type ChannelSetting struct {
	ID            string `storm:"id"`
	Enabled       bool   `storm:"index"`
	AdminUserIDs  []int
	Followings    *map[int][]string
	PushIntervals *map[int]int
}

func (cset *ChannelSetting) update(action int, param interface{}) {
	pint, iok := param.(ModuleInterval)
	puser, uok := param.(ModuleUser)
	switch action {
	case ChannelActionEnable:
		cset.Enabled = true
	case ChannelActionDisable:
		cset.Enabled = false
	case ChannelActionAddFollow:
		if uok {
			if cset.Followings == nil{
				followings := make(map[int][]string)
				cset.Followings = &followings
			}
			_, ok := (*cset.Followings)[puser.Module]
			if !ok{
				(*cset.Followings)[puser.Module] = make([]string, 0, 0)
			}
			(*cset.Followings)[puser.Module] = append((*cset.Followings)[puser.Module], puser.Username)
			_, ok = (*cset.PushIntervals)[puser.Module]
			if !ok{
				(*cset.PushIntervals)[puser.Module] = DefaultInterval
			}
		}
	case ChannelActionDelFollow:
		if uok {
			for i, u := range (*cset.Followings)[puser.Module] {
				if u == puser.Username {
					(*cset.Followings)[puser.Module] = append((*cset.Followings)[puser.Module][:i], (*cset.Followings)[puser.Module][i+1:]...)
					if len((*cset.Followings)[puser.Module]) == 0{
						delete(*cset.Followings, puser.Module)
					}
					break
				}
			}
		}
	case ChannelActionUpdatePushInterval:
		if cset.PushIntervals == nil{
			intervals := make(map[int]int)
			cset.PushIntervals = &intervals
		}
		if iok {
			(*cset.PushIntervals)[pint.Module] = pint.PushInterval
		}
	}
}

type Channel struct {
	*ChannelSetting
	DB    *storm.DB
	TGBOT *TelegramBot
	Control chan int
	Chat  *telebot.Chat
}

func (c *Channel) UpdateSettings(action int, param interface{}) {
	c.update(action, param)
	c.DB.Save(c.ChannelSetting)
}

func (c *Channel) PushModule(control chan int, chat *telebot.Chat, module_id int, followings []string, wait_time time.Duration) {
	var fetcher f.Fetcher
	switch module_id {
	case MTwitter:
		fetcher = &c.TGBOT.FetcherConfigs.Twitter
		fetcher.Init(c.TGBOT.Database)
	}
	for {
		log.Printf("Will check for update for module %s:%s", c.ID, strings.Join(followings, ","))
		next_start := time.After(wait_time)
		c.TGBOT.SendAll(chat, fetcher.GetPush(c.ID, followings))
		select {
		case <-control:
			log.Println("Received exit signal")
			return
		case <-next_start:
			log.Println("Sleeping")
			continue
		}
	}
}

func (c *Channel) Push() {
	for {
		controlers := make([]chan int, 0, len(*c.Followings))
		for module_id, followings := range *c.Followings {
			signal := make(chan int)
			controlers = append(controlers, signal)
			if len(followings) == 0{
				log.Printf("Module %d started but there's no followings.", module_id)
			}else{
				go c.PushModule(signal, c.Chat, module_id, followings, time.Duration((*c.PushIntervals)[module_id])*time.Second)
				log.Printf("Module %d started:%s.", module_id, strings.Join(followings, ","))
			}
		}
		select{
		case t := <- c.Control:
			log.Printf("Receive signal %d.", t)
			for _, cc := range controlers{
				cc <- 1
			}
			if t == SignalExit{
				return
			}else if t == SignalReload{
				continue
			}
		}
	}
}

func (c *Channel) Reload() {
	c.Control <- SignalReload
}

func (c *Channel) Exit(){
	c.Control <- SignalExit
}

func (c *Channel) Enable() {
	c.Enabled = true
	c.Reload()
}

func (c *Channel) Disable(){
	c.Enabled = false
	c.Reload()
}

func (c *Channel) AddFollowing(user ModuleUser){
	c.UpdateSettings(ChannelActionAddFollow, user)
	c.Reload()
}

func (c *Channel) DelFollowing(user ModuleUser){
	c.UpdateSettings(ChannelActionDelFollow, user)
	c.Reload()
}

func (c *Channel) UpdateInterval(interval ModuleInterval){
	c.UpdateSettings(ChannelActionUpdatePushInterval, interval)
	c.Reload()
}

func Str2Module(s string) int{
	switch s{
	case "twitter":
		return MTwitter
	default:
		return -1
	}
}

func Module2Str(i int) string{
	switch i{
	case MTwitter:
		return "twitter"
	default:
		return ""
	}
}

func MakeChannels(TGBOT *TelegramBot) []*Channel {
	db := TGBOT.Database
	var channel_settings []ChannelSetting
	err := db.All(&channel_settings)
	if err != nil {
		log.Fatal("Cannot read channel settings.", err)
		return []*Channel{}
	}
	var channels []*Channel
	for i := range channel_settings {
		chat, err := TGBOT.Bot.ChatByID(channel_settings[i].ID)
		if err != nil {
			log.Fatal("Error when start chat.", err)
		}
		channels = append(channels, &Channel{&channel_settings[i], db, TGBOT, make(chan int), chat})
	}
	return channels
}

func AddChannelIfNotExists(TGBOT *TelegramBot, channel_id string) *Channel{
	db := TGBOT.Database
	followings := make(map[int][]string)
	intervals := make(map[int]int)
	channel_setting := ChannelSetting{ID: channel_id, Enabled:true, AdminUserIDs:make([]int, 0, 0),
	Followings: &followings, PushIntervals: &intervals}
	// Check if channel exists and add it.
	// Use tx to make it safe.
	db.Save(&channel_setting)
	chat, err := TGBOT.Bot.ChatByID(channel_id)
	if err != nil {
		log.Println("Error when start chat.", err)
	}
	log.Println("Channel added.")
	return &Channel{&channel_setting, db, TGBOT, make(chan int), chat}
}

func RunPusher(TGBOT *TelegramBot){
	channels := MakeChannels(TGBOT)
	TGBOT.Channels = &channels
	for _, c := range channels{
		go c.Push()
	}
}
