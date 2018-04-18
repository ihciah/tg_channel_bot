package fetchers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type V2EXFetcher struct {
	BaseFetcher
}

type V2EXHot []struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	Content         string `json:"content"`
	ContentRendered string `json:"content_rendered"`
	Replies         int    `json:"replies"`
	Member          struct {
		ID           int    `json:"id"`
		Username     string `json:"username"`
		Tagline      string `json:"tagline"`
		AvatarMini   string `json:"avatar_mini"`
		AvatarNormal string `json:"avatar_normal"`
		AvatarLarge  string `json:"avatar_large"`
	} `json:"member"`
	Node struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		Title            string `json:"title"`
		TitleAlternative string `json:"title_alternative"`
		URL              string `json:"url"`
		Topics           int    `json:"topics"`
		AvatarMini       string `json:"avatar_mini"`
		AvatarNormal     string `json:"avatar_normal"`
		AvatarLarge      string `json:"avatar_large"`
	} `json:"node"`
	Created      int `json:"created"`
	LastModified int `json:"last_modified"`
	LastTouched  int `json:"last_touched"`
}

func (f *V2EXFetcher) GetPush(string, []string) []ReplyMessage {
	api_url := "https://www.v2ex.com/api/topics/hot.json"
	resp, err := f.HTTPGet(api_url)
	if err != nil {
		log.Println("Unable to crawl v2ex api", err)
		return []ReplyMessage{{Err: err}}
	}
	resp_content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Unable to read response", err)
		return []ReplyMessage{{Err: err}}
	}
	hot := V2EXHot{}
	if err := json.Unmarshal(resp_content, &hot); err != nil {
		log.Println("Unable to load json", err)
		return []ReplyMessage{{Err: err}}
	}
	titles := make([]string, 0, 10)
	for _, v := range hot {
		titles = append(titles, v.Title)
	}
	return []ReplyMessage{{Caption: strings.Join(titles, "\n")}}
}
