package main

import (
	"github.com/robfig/cron"
	"log"
)

type CronConfig struct {
}

func RunCron(TGBOT *TelegramBot) {
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		log.Println("test task")
	})
}
