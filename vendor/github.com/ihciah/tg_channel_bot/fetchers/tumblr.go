package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asdine/storm"
	"github.com/dghubble/oauth1"
	"github.com/patrickmn/go-cache"
	"log"
	"strings"
	"time"
)

type TumblrPosts struct {
	Meta struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
	Response struct {
		Posts []struct {
			Type               string `json:"type"`
			BlogName           string `json:"blog_name"`
			ID                 int64  `json:"id"`
			PostURL            string `json:"post_url"`
			Slug               string `json:"slug"`
			Date               string `json:"date"`
			Timestamp          int    `json:"timestamp"`
			State              string `json:"state"`
			Format             string `json:"format"`
			ShortURL           string `json:"short_url"`
			IsBlocksPostFormat bool   `json:"is_blocks_post_format"`
			SourceURL          string `json:"source_url,omitempty"`
			SourceTitle        string `json:"source_title,omitempty"`
			Caption            string `json:"caption,omitempty"`
			Reblog             struct {
				Comment  string `json:"comment"`
				TreeHTML string `json:"tree_html"`
			} `json:"reblog"`
			Trail []struct {
				Post struct {
					ID interface{} `json:"id"`
				} `json:"post"`
				ContentRaw string `json:"content_raw"`
				Content    string `json:"content"`
			} `json:"trail"`
			VideoURL        string `json:"video_url,omitempty"`
			ThumbnailURL    string `json:"thumbnail_url,omitempty"`
			ThumbnailWidth  int    `json:"thumbnail_width,omitempty"`
			ThumbnailHeight int    `json:"thumbnail_height,omitempty"`
			Duration        int    `json:"duration,omitempty"`
			VideoType       string `json:"video_type,omitempty"`
			DisplayAvatar   bool   `json:"display_avatar"`
			PhotosetLayout  string `json:"photoset_layout,omitempty"`
			Photos          []struct {
				Caption      string `json:"caption"`
				OriginalSize struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"original_size"`
			} `json:"photos,omitempty"`
			ImagePermalink string `json:"image_permalink,omitempty"`
			Title          string `json:"title,omitempty"`
			Body           string `json:"body,omitempty"`
		} `json:"posts"`
	} `json:"response"`
}

type TumblrFetcher struct {
	BaseFetcher
	OAuthConsumerKey    string `json:"consumer_key"`
	OAuthConsumerSecret string `json:"consumer_secret"`
	OAuthToken          string `json:"access_token"`
	OAuthTokenSecret    string `json:"access_token_secret"`
	cache               *cache.Cache
	channel_id          string
}

func (f *TumblrFetcher) Init(db *storm.DB, channel_id string) (err error) {
	f.DB = db.From("tumblr")
	f.cache = cache.New(cacheExp*time.Hour, cachePurge*time.Hour)
	f.channel_id = channel_id
	config := oauth1.NewConfig(f.OAuthConsumerKey, f.OAuthConsumerSecret)
	token := oauth1.NewToken(f.OAuthToken, f.OAuthTokenSecret)
	f.client = *config.Client(oauth1.NoContext, token)
	return
}

func (f *TumblrFetcher) getUserTimeline(user string, time int64) ([]ReplyMessage, error) {
	if f.OAuthConsumerKey == "" {
		return []ReplyMessage{}, errors.New("Need API key.")
	}
	api_url := fmt.Sprintf("https://api.tumblr.com/v2/blog/%s.tumblr.com/posts", user)
	resp_content, err := f.HTTPGet(api_url)
	if err != nil {
		log.Println("Unable to request tumblr api", err)
		return []ReplyMessage{}, err
	}
	posts := TumblrPosts{}
	if err := json.Unmarshal(resp_content, &posts); err != nil {
		log.Println("Unable to load json", err)
		return []ReplyMessage{}, err
	}
	if posts.Meta.Status != 200 {
		log.Println("Tumblr return err. Code", posts.Meta.Status)
		return []ReplyMessage{}, errors.New("Tumblr api error.")
	}
	ret := make([]ReplyMessage, 0, len(posts.Response.Posts))
	for _, p := range posts.Response.Posts {
		if p.Type != "photo" && p.Type != "video" {
			continue
		}
		if int64(p.Timestamp) < time {
			break
		}

		res := make([]Resource, 0, len(p.Photos))
		for _, photo := range p.Photos {
			tType := TIMAGE
			if strings.HasSuffix(strings.ToLower(photo.OriginalSize.URL), ".gif") {
				tType = TVIDEO
			}
			// Duplicate
			strsplit := strings.Split(photo.OriginalSize.URL, "/")
			if len(strsplit) < 4 {
				continue
			}
			imghash := fmt.Sprintf("%s@%s", f.channel_id, strsplit[3])
			_, found := f.cache.Get(imghash)
			f.cache.Set(imghash, true, cache.DefaultExpiration)
			if found {
				continue
			}

			// Blacklist
			is_blocked := false
			if err := f.DB.Get("block", imghash, &is_blocked); err == nil && is_blocked {
				continue
			}

			res = append(res, Resource{photo.OriginalSize.URL, tType, photo.OriginalSize.URL})
		}
		if p.VideoURL != "" {
			urlpath := strings.Split(p.VideoURL, "/")
			videopath := urlpath[len(urlpath)-1]
			if strings.Contains(videopath, ".") {
				videohash := fmt.Sprintf("%s@%s", f.channel_id, urlpath[len(urlpath)-1])
				_, found := f.cache.Get(videohash)
				f.cache.Set(videohash, true, cache.DefaultExpiration)
				if !found {
					res = append(res, Resource{p.VideoURL, TVIDEO, p.VideoURL})
				}
			} else {
				res = append(res, Resource{p.VideoURL, TVIDEO, p.VideoURL})
			}
		}
		if len(res) > 0 {
			ret = append(ret, ReplyMessage{res, p.ShortURL, nil})
		}
	}
	return ret, nil
}

func (f *TumblrFetcher) GetPush(userid string, followings []string) []ReplyMessage {
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

func (f *TumblrFetcher) GoBack(userid string, back int64) error {
	now := time.Now().Unix()
	if back > now {
		return errors.New("Back too long!")
	}
	return f.DB.Set("last_update", userid, now-back)
}

func (f *TumblrFetcher) Block(caption string) string {
	strsplit := strings.Split(caption, "/")
	if len(strsplit) >= 4 {
		imghash := fmt.Sprintf("%s@%s", f.channel_id, strsplit[3])
		f.DB.Set("block", imghash, true)
		return fmt.Sprintf("%s blocked.", imghash)
	}
	return "Unrecognized image caption."
}
