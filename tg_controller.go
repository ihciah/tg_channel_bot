package main

import (
	"encoding/json"
	"errors"
	"github.com/asdine/storm"
	"github.com/coreos/bbolt"
	tb "github.com/ihciah/telebot"
	f "github.com/ihciah/tg_channel_bot/fetchers"
	"io/ioutil"
	"log"
	"time"
	"os"
)

const MaxAlbumSize = 10

type TelegramBot struct {
	Bot            *tb.Bot
	Database       *storm.DB
	Token          string        `json:"token"`
	Timeout        int           `json:"timeout"`
	DatabasePath   string        `json:"database"`
	FetcherConfigs FetcherConfig `json:"fetcher_config"`
	Channels       *[]*Channel
	Admins         []string `json:"admins"`
}

func (TGBOT *TelegramBot) LoadConfigFromEnv() {
	TGBOT.Token = os.Getenv("BOT_TOKEN")
	TGBOT.DatabasePath = "database.db"
	TGBOT.Timeout = 120

	admin := []string{os.Getenv("ADMIN_NAME")}
	TGBOT.Admins = admin

	var tumblrfetcher f.TumblrFetcher
	tumblrfetcher.OAuthConsumerKey = os.Getenv("TUMBLR_KEY")
	tumblrfetcher.OAuthConsumerSecret = os.Getenv("TUMBLR_SECRET")
	tumblrfetcher.OAuthToken = os.Getenv("TUMBLR_TOKEN")
	tumblrfetcher.OAuthTokenSecret = os.Getenv("TUMBLR_TOKEN_SECRET")

	var twitterfetcher f.TwitterFetcher
	twitterfetcher.AccessTokenSecret = os.Getenv("TWITTER_TOKEN_SECRET")
	twitterfetcher.AccessToken = os.Getenv("TWITTER_TOKEN")
	twitterfetcher.ConsumerKey = os.Getenv("TWITTER_KEY")
	twitterfetcher.ConsumerSecret = os.Getenv("TWITTER_SECRET")

	var fetcherconfig FetcherConfig
	fetcherconfig.Tumblr = tumblrfetcher
	fetcherconfig.Twitter = twitterfetcher

	TGBOT.FetcherConfigs = fetcherconfig

	var err error
	TGBOT.Bot, err = tb.NewBot(tb.Settings{
		Token:       TGBOT.Token,
		Poller:      &tb.LongPoller{Timeout: time.Duration(TGBOT.Timeout) * time.Second},
		HTTPTimeout: TGBOT.Timeout,
	})
	if err != nil {
		log.Fatal("[Cannot initialize telegram Bot]", err)
		return
	}

	TGBOT.Database, err = storm.Open(TGBOT.DatabasePath, storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		log.Fatal("[Cannot initialize database]", err)
	}
	log.Printf("[Bot initialized]Token: %s\nTimeout: %d\n", TGBOT.Token, TGBOT.Timeout)
}

func (TGBOT *TelegramBot) LoadConfig(json_path string) {
	data, err := ioutil.ReadFile(json_path)
	if err != nil {
		log.Fatal("[Cannot read telegram config]", err)
		return
	}
	if err := json.Unmarshal(data, TGBOT); err != nil {
		log.Fatal("[Cannot parse telegram config]", err)
		return
	}
	TGBOT.Bot, err = tb.NewBot(tb.Settings{
		Token:       TGBOT.Token,
		Poller:      &tb.LongPoller{Timeout: time.Duration(TGBOT.Timeout) * time.Second},
		HTTPTimeout: TGBOT.Timeout,
	})
	if err != nil {
		log.Fatal("[Cannot initialize telegram Bot]", err)
		return
	}

	TGBOT.Database, err = storm.Open(TGBOT.DatabasePath, storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}))
	if err != nil {
		log.Fatal("[Cannot initialize database]", err)
	}
	log.Printf("[Bot initialized]Token: %s\nTimeout: %d\n", TGBOT.Token, TGBOT.Timeout)
}

func (TGBOT *TelegramBot) Serve() {
	TGBOT.RegisterHandler()
	TGBOT.Bot.Start()
}

func (TGBOT *TelegramBot) Send(to tb.Recipient, message f.ReplyMessage) error {
	if message.Err != nil {
		return message.Err
	}

	if len(message.Resources) == 1 {
		var err error
		var mediaFile tb.InputMedia
		if message.Resources[0].T == f.TIMAGE {
			mediaFile = &tb.Photo{File: tb.FromURL(message.Resources[0].URL), Caption: message.Resources[0].Caption}
		} else if message.Resources[0].T == f.TVIDEO {
			mediaFile = &tb.Video{File: tb.FromURL(message.Resources[0].URL), Caption: message.Resources[0].Caption}
		} else {
			return errors.New("Undefined message type.")
		}
		_, err = TGBOT.Bot.Send(to, mediaFile)
		return err
	}

	if len(message.Resources) == 0 {
		if _, err := TGBOT.Bot.Send(to, message.Caption); err != nil {
			log.Println("Unable to send text:", message.Caption)
			return err
		} else {
			log.Println("Sent text:", message.Caption)
		}
	}

	var ret error
	for i := 0; i < len(message.Resources); i += MaxAlbumSize {
		end := i + MaxAlbumSize
		if end > len(message.Resources) {
			end = len(message.Resources)
		}
		mediaFiles := make(tb.Album, 0, MaxAlbumSize)
		for _, r := range message.Resources[i:end] {
			if r.T == f.TIMAGE {
				mediaFiles = append(mediaFiles, &tb.Photo{File: tb.FromURL(r.URL), Caption: r.Caption})
			} else if r.T == f.TVIDEO {
				mediaFiles = append(mediaFiles, &tb.Video{File: tb.FromURL(r.URL), Caption: r.Caption})
			} else {
				continue
			}
		}
		if _, err := TGBOT.Bot.SendAlbum(to, mediaFiles); err != nil {
			log.Println("Unable to send album", err)
			ret = err
		} else {
			log.Println("Sent album")
		}
	}
	return ret
}

func (TGBOT *TelegramBot) SendAll(to tb.Recipient, messages []f.ReplyMessage) (err error) {
	err = nil
	for _, msg := range messages {
		go TGBOT.Send(to, msg)
	}
	return
}
