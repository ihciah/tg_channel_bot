package fetchers

import (
	"github.com/asdine/storm"
	"time"
	"errors"
)

type TumblrFetcher struct {
	BaseFetcher
}

func (f *TumblrFetcher) Init(db *storm.DB) (err error) {
	f.DB = db.From("tumblr")
	return
}

func (f *TumblrFetcher) GoBack(userid string, back int64) error {
	now := time.Now().Unix()
	if back > now {
		return errors.New("Back too long!")
	}
	return f.DB.Set("last_update", userid, now-back)
}

func (f *TumblrFetcher) GetPush(userid string, followings []string) []ReplyMessage {
	var last_update int64
	if err := f.DB.Get("last_update", userid, &last_update); err != nil {
		last_update = 0
	}
	ret := make([]ReplyMessage, 0, 0)
	// Fetch media here.
	return ret
}
