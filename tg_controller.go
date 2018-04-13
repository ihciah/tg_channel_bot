package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"encoding/json"
	"io/ioutil"
	"time"
	"errors"
	f "github.com/ihciah/tg_channel_bot/fetchers"
)

type TelegramBot struct {
	bot *tb.Bot
	Token     string `json:"token"`
	Timeout   int `json:"timeout"`
}

func (TGBOT *TelegramBot) LoadConfig(json_path string) (err error){
	data, err := ioutil.ReadFile(json_path)
	if err != nil{
		log.Fatal("[Cannot read telegram config]", err)
		return err
	}
	if err := json.Unmarshal(data, TGBOT); err!=nil{
		log.Fatal("[Cannot parse telegram config]", err)
		return err
	}
	TGBOT.bot, err = tb.NewBot(tb.Settings{
		Token:  TGBOT.Token,
		Poller: &tb.LongPoller{Timeout: time.Duration(TGBOT.Timeout) * time.Second},
	})
	if err != nil{
		log.Fatal("[Cannot initialize telegram bot]", err)
		return err
	}
	log.Printf("[Bot initialized]Token: %s\nTimeout: %d\n", TGBOT.Token, TGBOT.Timeout)
	return
}

func (TGBOT *TelegramBot) Serve(){
	TGBOT.RegisterHandler()
	TGBOT.bot.Start()
}

func (TGBOT *TelegramBot)Send(to tb.Recipient, message f.ReplyMessage) error{
	switch message.T{
	case f.TERROR:
		err, ok := message.Resources.(error)
		if ok{
			return err
		}
		return errors.New("[Unknown error] cannot convert types")
	case f.TIMAGE:
		switch v := message.Resources.(type){
		case string:
			if _, err := TGBOT.bot.Send(to, &tb.Photo{File: tb.FromURL(v)}); err != nil{
				log.Println("Sent image with URL ", v)
			}else {
				log.Println("Unable to send image with URL ", v)
				return err
			}
		default:
			log.Println("Unable to convert image")
		}
	case f.TTEXT:
		text, ok := message.Resources.(string)
		if ok{
			if _, err := TGBOT.bot.Send(to, text); err!=nil{
				log.Println("Sent text ", text)
			}else{
				log.Println("Unable to send text")
				return err
			}
		}else {
			return errors.New("[Unknown error] cannot convert types")
		}
	case f.TVIDEO:
		switch v := message.Resources.(type){
		case string:
			if _, err := TGBOT.bot.Send(to, &tb.Video{File: tb.FromURL(v)}); err != nil{
				log.Println("Sent video with URL ", v)
			}else {
				log.Println("Unable to send video with URL ", v)
				return err
			}
		default:
			log.Println("Unable to convert video")
		}
	}
	return nil
}

