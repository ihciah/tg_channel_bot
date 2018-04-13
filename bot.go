package main

import "flag"

func parseCli() (string, bool){
	var config_file string
	var verbose bool
	flag.BoolVar(&verbose, "v", true, "verbose mode")
	flag.StringVar(&config_file, "c", "config.json", "config file path")
	flag.Parse()
	return config_file, verbose
}

func main(){
	config_path, _ := parseCli()
	t := TelegramBot{}
	t.LoadConfig(config_path)
	t.Serve()
}

