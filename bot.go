package main

import (
	"flag"
	"os"
	"os/signal"
	"log"
)

func parseCli() (string, bool) {
	var config_file string
	var verbose bool
	flag.BoolVar(&verbose, "v", true, "verbose mode")
	flag.StringVar(&config_file, "c", "config.json", "config file path")
	flag.Parse()
	return config_file, verbose
}

func ListenExit(TGBOT *TelegramBot){
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	_ = <-c
	TGBOT.Bot.Stop()
	log.Println("Exit.")
}

func main() {
	config_path, _ := parseCli()
	t := TelegramBot{}
	t.LoadConfig(config_path)
	RunCron(&t)
	ListenExit(&t)
	t.Serve()
}
