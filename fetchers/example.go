package fetchers

import (
	"errors"
	"github.com/go-xmlpath/xmlpath"
	"log"
	"net/http"
)

type ExampleFetcher struct {
	BaseFetcher
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
