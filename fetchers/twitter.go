package fetchers

import (
	"errors"
	"github.com/ChimeraCoder/anaconda"
	"github.com/asdine/storm"
	"net/url"
	"time"
)

type TwitterFetcher struct {
	BaseFetcher
	api               *anaconda.TwitterApi
	AccessToken       string `json:"access_token"`
	AccessToeknSecret string `json:"access_token_secret"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
}

const (
	MaxTweetCount = "10"
)

func (f *TwitterFetcher) Init(db *storm.DB) (err error) {
	f.DB = db.From("twitter")
	f.api = anaconda.NewTwitterApiWithCredentials(f.AccessToken, f.AccessToeknSecret, f.ConsumerKey, f.ConsumerSecret)
	return
}

func (f *TwitterFetcher) getUserTimeline(user string, time int64) ([]ReplyMessage, error) {
	v := url.Values{}
	v.Set("count", MaxTweetCount)
	v.Set("screen_name", user)
	results, err := f.api.GetUserTimeline(v)
	if err != nil {
		return []ReplyMessage{}, err
	}
	ret := make([]ReplyMessage, 0, len(results))
	for _, tweet := range results {
		t, err := tweet.CreatedAtTime()
		if err != nil {
			continue
		}
		tweet_time := t.Unix()
		if tweet_time < time {
			break
		}
		resources := make([]Resource, 0, len(tweet.ExtendedEntities.Media))
		for _, media := range tweet.ExtendedEntities.Media {
			var rType int
			switch media.Type {
			case "photo":
				rType = TIMAGE
			case "video":
				rType = TVIDEO
			}
			resources = append(resources, Resource{media.Media_url_https, rType})
		}
		ret = append(ret, ReplyMessage{resources, tweet.FullText, nil})
	}
	return ret, nil
}

func (f *TwitterFetcher) GetPush(userid string, followings []string) []ReplyMessage {
	var last_update int64
	if err := f.DB.Get("last_update", userid, &last_update); err != nil {
		last_update = 0
	}
	ret := make([]ReplyMessage, 0, 0)
	for _, follow := range followings {
		single, err := f.getUserTimeline(follow, last_update)
		if err == nil {
			ret = append(ret, single...)
		}
	}
	if len(ret) != 0 {
		f.DB.Set("last_update", userid, time.Now().Unix())
	}
	return ret
}

func (f *TwitterFetcher) GoBack(userid string, back int64) error {
	now := time.Now().Unix()
	if back > now {
		return errors.New("Back too long!")
	}
	return f.DB.Set("last_update", userid, now-back)
}
