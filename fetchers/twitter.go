package fetchers

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/boltdb/bolt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type TwitterFetcher struct {
	BaseFetcher
	api               *anaconda.TwitterApi
	AccessToken       string `json:"access_token"`
	AccessToeknSecret string `json:"access_token_secret"`
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	GetUser           string `json:"get_user_id"`
	GetCount          int    `json:"get_count"`
}

func (f *TwitterFetcher) Init(db *bolt.DB) {
	f.DB = db
	f.api = anaconda.NewTwitterApiWithCredentials("340822756-SxzEX4nN4I5OMEse2DalE5LPXH44nkv7eiABKSWg",
		"GeyR5WA05HemvAryOu2nuyNym4Zgbz9SCnwozKSQ0JMXc",
		"ZaRyEEFjuft0iXZDBmydsdbCX",
		"ApGOwDwC2Mf7tCQ1OS9rOBe4wap1ZTNq5l7nV0ajJTDp2CKh9I")
}

func (f *TwitterFetcher) Get() ReplyMessage {
	v := url.Values{}
	v.Set("count", strconv.Itoa(f.GetCount))
	v.Set("screen_name", f.GetUser)
	results, err := f.api.GetUserTimeline(v)
	if err != nil {
		log.Fatal(err)
		return ReplyMessage{err, TERROR}
	}
	tweets := make([]string, 0, len(results))
	for _, tweet := range results {
		tweets = append(tweets, fmt.Sprintf("%s(%s): \n%s", tweet.User.Name, tweet.User.ScreenName, tweet.FullText))
	}
	return ReplyMessage{strings.Join(tweets, "\n"), TTEXT}
}

func (f *TwitterFetcher) GetPush() []ReplyMessage {
	return make([]ReplyMessage, 0, 0)
}
