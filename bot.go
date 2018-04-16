package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

func parseCli() (string, bool) {
	var config_file string
	var verbose bool
	flag.BoolVar(&verbose, "v", true, "verbose mode")
	flag.StringVar(&config_file, "c", "config.json", "config file path")
	flag.Parse()
	return config_file, verbose
}

func ListenExit(TGBOT *TelegramBot) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	_ = <-c
	TGBOT.Bot.Stop()
	TGBOT.Database.Close()
	log.Println("Exit.")
}

func main() {
	config_path, _ := parseCli()
	t := TelegramBot{}
	if err := t.LoadConfig(config_path); err != nil {
		log.Fatal("Cannot load config", err)
		return
	}
	RunPusher(&t)
	go ListenExit(&t)
	t.Serve()
}
