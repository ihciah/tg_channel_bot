package fetchers

import (
	"github.com/dghubble/sling"
	"log"
	"net/http"
)

const (
	TIMAGE = iota
	TVIDEO
	TTEXT
	TERROR
)

type ReplyMessage struct {
	Resources interface{}
	T         int
}

type Fetcher interface {
	Init()
	Get() ReplyMessage
}

type BaseFetcher struct {
	UA     string
	sling  *sling.Sling
	client http.Client
}

func (f *BaseFetcher) Init() {
	f.UA = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	f.client = http.Client{}
	f.sling = sling.New().Client(&f.client).Set("User-Agent", f.UA)
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
