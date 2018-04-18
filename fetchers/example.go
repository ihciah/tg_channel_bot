package fetchers

import (
	"errors"
	"github.com/go-xmlpath/xmlpath"
	"log"
	"bytes"
)

type ExampleFetcher struct {
	BaseFetcher
}

func (f *ExampleFetcher) GetPush(string, []string) []ReplyMessage {
	page_url := "https://www.v2ex.com/i/R7yApIA5.jpeg"
	response, err := f.HTTPGet(page_url)
	img_url, err := f.parse_img_page(response)
	if err != nil {
		log.Println("Cannot do parse", err)
		return []ReplyMessage{{Err: err}}
	}
	log.Println("Image url get", img_url)
	reply := ReplyMessage{[]Resource{{URL: img_url, T: TIMAGE}}, "", nil}
	return []ReplyMessage{reply}
}

func (f *ExampleFetcher) parse_img_page(resp []byte) (string, error) {
	path := xmlpath.MustCompile("//input[@class='sls']/@value")
	root, err := xmlpath.ParseHTML(bytes.NewBuffer(resp))
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
