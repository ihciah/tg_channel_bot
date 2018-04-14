package main

import (
	"github.com/robfig/cron"
	"log"
)

type CronTab struct {
	cron *cron.Cron
}

func RunCron(TGBOT *TelegramBot) {
	c := cron.New()
	c.AddFunc("* * * * *", func() {
		log.Println("test task")
	})
}
