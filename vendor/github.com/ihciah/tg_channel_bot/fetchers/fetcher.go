package fetchers

import (
	"errors"
	"github.com/asdine/storm"
	"github.com/dghubble/sling"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	TIMAGE = iota
	TVIDEO
)

const (
	cacheExp   = 168
	cachePurge = 336
)

type Resource struct {
	URL     string
	T       int
	Caption string
}

type ReplyMessage struct {
	Resources []Resource
	Caption   string
	Err       error
}

type Fetcher interface {
	Init(*storm.DB, string) error                      // Initializing
	GetPush(string, []string) []ReplyMessage           // For channel message
	GetPushAtLeastOne(string, []string) []ReplyMessage // For user message
	GoBack(string, int64) error                        // Set last update time to N seconds before
	Block(string) string
}

type BaseFetcher struct {
	UA     string
	DB     storm.Node
	sling  *sling.Sling
	client http.Client
}

// Initialize
func (f *BaseFetcher) Init(db *storm.DB, _ string) error {
	f.UA = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	f.client = http.Client{Timeout: time.Duration(30) * time.Second}
	return nil
}

func (f *BaseFetcher) HTTPGet(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Cannot create request", err)
		return []byte{}, err
	}
	request.Close = true
	request.Header.Set("User-Agent", f.UA)
	response, err := f.client.Do(request)
	if err != nil {
		log.Println("Cannot do request", err)
		return []byte{}, err
	}
	defer response.Body.Close()
	resp_content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Unable to read response", err)
		return []byte{}, err
	}
	return resp_content, nil
}

// For channel update
func (f *BaseFetcher) GetPush(string, []string) []ReplyMessage {
	return []ReplyMessage{{Caption: "Unsupported. You should define GetPush function first."}}
}

// For user request update
func (f *BaseFetcher) GetPushAtLeastOne(userid string, following []string) (ret []ReplyMessage) {
	ret = f.GetPush(userid, following)
	if len(ret) == 0 {
		ret = []ReplyMessage{{Caption: "No new updates."}}
	}
	return
}

// Set last update time to several seconds before
func (f *BaseFetcher) GoBack(string, int64) error {
	return errors.New("Time machine unsupported for this site.")
}

func (f *BaseFetcher) Block(string) string {
	return "Unimplement."
}
