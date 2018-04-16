package fetchers

import (
	"github.com/asdine/storm"
	"github.com/dghubble/sling"
	"log"
	"net/http"
)

const (
	TIMAGE = iota
	TVIDEO
)

type Resource struct{
	URL string
	T int
}

type ReplyMessage struct {
	Resources []Resource
	Caption string
	Err error
}

type Fetcher interface {
	Init(*storm.DB) error // Initializing
	GetPush(string, []string) []ReplyMessage
	GetPushAtLeastOne(string, []string) []ReplyMessage
}

type BaseFetcher struct {
	UA     string
	DB     storm.Node
	sling  *sling.Sling
	client http.Client
}

func (f *BaseFetcher) Init(db *storm.DB) error {
	f.UA = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	f.client = http.Client{}
	f.sling = sling.New().Client(&f.client).Set("User-Agent", f.UA)
	return nil
}

func (f *BaseFetcher) HTTPGet(url string) (*http.Response, error) {
	var resp *http.Response
	request, err := f.sling.Get(url).Request()
	if err != nil {
		log.Fatal("Cannot create request", err)
		return resp, err
	}
	response, err := f.client.Do(request)
	if err != nil {
		log.Fatal("Cannot do request", err)
		return resp, err
	}
	return response, nil
}

func (f *BaseFetcher) GetPush(string, []string) []ReplyMessage {
	return []ReplyMessage{}
}

func (f *BaseFetcher) GetPushAtLeastOne(string, []string) []ReplyMessage {
	return []ReplyMessage{}
}
