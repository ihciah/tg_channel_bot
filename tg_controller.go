package main

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	f "github.com/ihciah/tg_channel_bot/fetchers"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"time"
)

type TelegramBot struct {
	Bot            *tb.Bot
	Database       *bolt.DB
	Token          string        `json:"token"`
	Timeout        int           `json:"timeout"`
	DatabasePath   string        `json:"database"`
	FetcherConfigs FetcherConfig `json:"fetcher_config"`
}

func (TGBOT *TelegramBot) LoadConfig(json_path string) (err error) {
	data, err := ioutil.ReadFile(json_path)
	if err != nil {
		log.Fatal("[Cannot read telegram config]", err)
		return err
	}
	if err := json.Unmarshal(data, TGBOT); err != nil {
		log.Fatal("[Cannot parse telegram config]", err)
		return err
	}
	TGBOT.Bot, err = tb.NewBot(tb.Settings{
		Token:  TGBOT.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(TGBOT.Timeout) * time.Second},
	})
	if err != nil {
		log.Fatal("[Cannot initialize telegram Bot]", err)
		return err
	}
	TGBOT.Database, err = bolt.Open(TGBOT.DatabasePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	log.Printf("[Bot initialized]Token: %s\nTimeout: %d\n", TGBOT.Token, TGBOT.Timeout)
	return
}

func (TGBOT *TelegramBot) Serve() {
	TGBOT.RegisterHandler()
	TGBOT.Bot.Start()
}

func (TGBOT *TelegramBot) Send(to tb.Recipient, message f.ReplyMessage) error {
	switch message.T {
	case f.TERROR:
		err, ok := message.Resources.(error)
		if ok {
			return err
		}
		return errors.New("[Unknown error] cannot convert types")
	case f.TIMAGE:
		switch v := message.Resources.(type) {
		case string:
			if _, err := TGBOT.Bot.Send(to, &tb.Photo{File: tb.FromURL(v)}); err != nil {
				log.Println("Unable to send image with URL ", v)
				return err
			} else {
				log.Println("Sent image with URL ", v)
			}
		default:
			log.Println("Unable to convert image")
		}
	case f.TTEXT:
		text, ok := message.Resources.(string)
		if ok {
			if _, err := TGBOT.Bot.Send(to, text); err != nil {
				log.Println("Unable to send text")
				return err
			} else {
				log.Println("Sent text ", text)
			}
		} else {
			return errors.New("[Unknown error] cannot convert types")
		}
	case f.TVIDEO:
		switch v := message.Resources.(type) {
		case string:
			if _, err := TGBOT.Bot.Send(to, &tb.Video{File: tb.FromURL(v)}); err != nil {
				log.Println("Unable to send video with URL ", v)
				return err
			} else {
				log.Println("Sent video with URL ", v)
			}
		default:
			log.Println("Unable to convert video")
		}
	}
	return nil
}
