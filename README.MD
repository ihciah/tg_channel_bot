# tg_channel_bot

A telegram bot designed for syncing text/image/video from given sites to telegram channel.

Dependences:

```
"github.com/ihciah/telebot"
"github.com/ihciah/tg_channel_bot/fetchers"
"github.com/dghubble/sling"
"github.com/go-xmlpath/xmlpath"
"github.com/ChimeraCoder/anaconda"
"github.com/asdine/storm"
"github.com/coreos/bbolt"
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
