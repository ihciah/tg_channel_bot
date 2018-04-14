# tg_channel_bot

A telegram bot designed for syncing text/image/video from given sites to telegram channel.

Dependences:

```
"gopkg.in/tucnak/telebot.v2"
"github.com/ihciah/tg_channel_bot/fetchers"
"github.com/robfig/cron"
"github.com/dghubble/sling"
"github.com/go-xmlpath/xmlpath"
"github.com/ChimeraCoder/anaconda"
```




You can run the program through `./tg_channel_bot -c config.json`. Check `config_example.json` for example of `config.json`.

Tips for users in China Mainland to set socks proxy:
- Linux:
```
export http_proxy=https://127.0.0.1:1080
export https_proxy=https://127.0.0.1:1080
```

- Windows:
```
set http_proxy=https://127.0.0.1:1080
set https_proxy=https://127.0.0.1:1080
```
