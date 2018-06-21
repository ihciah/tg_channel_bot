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
	flag.StringVar(&config_file, "c", "", "config file path")
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
	if config_path != "" {
		t.LoadConfig(config_path)
	} else {
		t.LoadConfigFromEnv()
	}
	RunPusher(&t)
	go ListenExit(&t)
	t.Serve()
}
