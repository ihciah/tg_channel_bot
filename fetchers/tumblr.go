package fetchers

import (
	"github.com/asdine/storm"
	"time"
	"errors"
	"log"
	"encoding/json"
	"fmt"
)

type TumblrPosts struct {
	Meta struct {
		Status               int    `json:"status"`
		Msg                  string `json:"msg"`
		XTumblrContentRating string `json:"x_tumblr_content_rating"`
	} `json:"meta"`
	Response struct {
		Blog struct {
			Ask                 bool   `json:"ask"`
			AskAnon             bool   `json:"ask_anon"`
			AskPageTitle        string `json:"ask_page_title"`
			CanSubscribe        bool   `json:"can_subscribe"`
			Description         string `json:"description"`
			IsAdult             bool   `json:"is_adult"`
			IsNsfw              bool   `json:"is_nsfw"`
			Name                string `json:"name"`
			Posts               int    `json:"posts"`
			ReplyConditions     string `json:"reply_conditions"`
			ShareLikes          bool   `json:"share_likes"`
			SubmissionPageTitle string `json:"submission_page_title"`
			Subscribed          bool   `json:"subscribed"`
			Title               string `json:"title"`
			TotalPosts          int    `json:"total_posts"`
			Updated             int    `json:"updated"`
			URL                 string `json:"url"`
			IsOptoutAds         bool   `json:"is_optout_ads"`
		} `json:"blog"`
		Posts []struct {
			Type               string        `json:"type"`
			BlogName           string        `json:"blog_name"`
			PostURL            string        `json:"post_url"`
			Slug               string        `json:"slug"`
			Date               string        `json:"date"`
			Timestamp          int           `json:"timestamp"`
			State              string        `json:"state"`
			Format             string        `json:"format"`
			ReblogKey          string        `json:"reblog_key"`
			Tags               []interface{} `json:"tags"`
			ShortURL           string        `json:"short_url"`
			Summary            string        `json:"summary"`
			IsBlocksPostFormat bool          `json:"is_blocks_post_format"`
			RecommendedSource  interface{}   `json:"recommended_source"`
			RecommendedColor   interface{}   `json:"recommended_color"`
			NoteCount          int           `json:"note_count"`
			SourceURL          string        `json:"source_url"`
			SourceTitle        string        `json:"source_title"`
			Caption            string        `json:"caption"`
			Reblog             struct {
				Comment  string `json:"comment"`
				TreeHTML string `json:"tree_html"`
			} `json:"reblog"`
			Trail []struct {
				Blog struct {
					Name   string `json:"name"`
					Active bool   `json:"active"`
					ShareLikes     bool `json:"share_likes"`
					ShareFollowing bool `json:"share_following"`
					CanBeFollowed  bool `json:"can_be_followed"`
				} `json:"blog"`
				ContentRaw    string `json:"content_raw"`
				Content       string `json:"content"`
				IsCurrentItem bool   `json:"is_current_item"`
			} `json:"trail"`
			VideoURL        string `json:"video_url"`
			HTML5Capable    bool   `json:"html5_capable"`
			ThumbnailURL    string `json:"thumbnail_url"`
			ThumbnailWidth  int    `json:"thumbnail_width"`
			ThumbnailHeight int    `json:"thumbnail_height"`
			Duration        int    `json:"duration"`
			Player          []struct {
				Width     int    `json:"width"`
				EmbedCode string `json:"embed_code"`
			} `json:"player"`
			VideoType        string `json:"video_type"`
			CanLike          bool   `json:"can_like"`
			CanReblog        bool   `json:"can_reblog"`
			CanSendInMessage bool   `json:"can_send_in_message"`
			CanReply         bool   `json:"can_reply"`
			DisplayAvatar    bool   `json:"display_avatar"`
			PhotosetLayout   string `json:"photoset_layout"`
			Photos           []struct {
				Caption      string `json:"caption"`
				OriginalSize struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"original_size"`
				AltSizes []struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"alt_sizes"`
			} `json:"photos"`
			ImagePermalink string `json:"image_permalink"`
			Title          string `json:"title"`
			Body           string `json:"body"`
		} `json:"posts"`
		TotalPosts int `json:"total_posts"`
	} `json:"response"`
}

type TumblrFetcher struct {
	BaseFetcher
	OAuthConsumerKey string  `json:"oauth_consumer_key"`
}

func (f *TumblrFetcher) Init(db *storm.DB) (err error) {
	f.DB = db.From("tumblr")
	return
}

func (f *TumblrFetcher) getUserTimeline(user string, time int64) ([]ReplyMessage, error){
	if f.OAuthConsumerKey == ""{
		return []ReplyMessage{}, errors.New("Need API key.")
	}
	api_url := fmt.Sprintf("https://api.tumblr.com/v2/blog/%s.tumblr.com/posts?api_key=%s", user, f.OAuthConsumerKey)
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
	if posts.Meta.Status != 200{
		log.Println("Tumblr return err. Code", posts.Meta.Status)
		return []ReplyMessage{}, errors.New("Tumblr api error.")
	}
	ret := make([]ReplyMessage, 0, len(posts.Response.Posts))
	for _, p := range posts.Response.Posts{
		if p.Type != "photo" && p.Type != "video"{
			continue
		}
		if int64(p.Timestamp) < time{
			break
		}
		res := make([]Resource, 0, len(p.Photos))
		for _, photo := range p.Photos{
			res = append(res, Resource{photo.OriginalSize.URL, TIMAGE})
		}
		if p.VideoURL != ""{
			res = append(res, Resource{p.VideoURL, TVIDEO})
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
