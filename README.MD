# tg_channel_bot

A telegram bot designed for syncing text/image/video from given sites to telegram channel.

Dependences:

```
"gopkg.in/tucnak/telebot.v2"
"github.com/ihciah/tg_channel_bot/fetchers"
"github.com/robfig/cron"
"github.com/dghubble/sling"
"github.com/go-xmlpath/xmlpath"
```



The current configure files are like the following:

```json
{
  "token": "YOUR_TELEGRAM_BOT_TOKEN",
  "timeout": 10
}
```

You can run the program through `./tg_channel_bot -c config.json`

Check `fetchers/example.go` and `fetchers/v2ex.go` if you want to write your own site crawlers.