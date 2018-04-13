package fetchers

import (
	"errors"
	"github.com/dghubble/sling"
	"github.com/go-xmlpath/xmlpath"
	"log"
	"net/http"
)

type ExampleFetcher struct {
	BaseFetcher
}

func (f *ExampleFetcher) Init() {
	f.UA = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	f.client = http.Client{}
	f.sling = sling.New().Client(&f.client).Set("User-Agent", f.UA)
}
func (f *ExampleFetcher) Get() ReplyMessage {
	page_url := "https://www.v2ex.com/i/R7yApIA5.jpeg"
	response, err := f.HTTPGet(page_url)
	img_url, err := f.parse_img_page(response)
	if err != nil {
		log.Fatal("Cannot do parse", err)
		return ReplyMessage{err, TERROR}
	}
	log.Println("Image url get", img_url)
	return ReplyMessage{img_url, TIMAGE}
}
func (f *ExampleFetcher) parse_img_page(resp *http.Response) (string, error) {
	path := xmlpath.MustCompile("//input[@class='sls']/@value")
	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil {
		return "", err
	}
	results := make([]string, 0, 5)
	iter := path.Iter(root)
	for iter.Next() {
		results = append(results, iter.Node().String())
	}
	if len(results) >= 2 {
		return results[1], nil
	}
	return "", errors.New("Unable to parse html.")
}
