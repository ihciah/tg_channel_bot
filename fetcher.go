package main

import (
	"net/http"
	"github.com/dghubble/sling"
	"github.com/go-xmlpath/xmlpath"
	"log"
	"github.com/pkg/errors"
)

const (
	tIMAGE = iota
	tVIDEO
	tTEXT
	tERROR
)

type ReplyMessage struct{
	resources interface{}
	t int
}

type Fetcher interface{
	Init()
	Get() ReplyMessage
}

func CreateFetcher(fetcher Fetcher) *Fetcher{
	fetcher.Init()
	return &fetcher
}


type ExampleFetcher struct{
	UA string
	sling *sling.Sling
	client http.Client
}

func (f *ExampleFetcher)Init(){
	f.UA = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"
	f.client = http.Client{}
	f.sling = sling.New().Client(&f.client).Set("User-Agent", f.UA)
}
func (f *ExampleFetcher)Get() ReplyMessage{
	page_url := "https://www.v2ex.com/i/R7yApIA5.jpeg"
	request, err := f.sling.Get(page_url).Request()
	if err != nil{
		log.Fatal("Cannot create request", err)
		return ReplyMessage{err, tERROR}
	}
	response, err := f.client.Do(request)
	if err != nil{
		log.Fatal("Cannot do request", err)
		return ReplyMessage{err, tERROR}
	}
	img_url, err := f.parse_img_page(response)
	if err != nil{
		log.Fatal("Cannot do parse", err)
		return ReplyMessage{err, tERROR}
	}
	log.Println("Image url get", img_url)
	return ReplyMessage{img_url, tIMAGE}
}
func (f *ExampleFetcher)parse_img_page(resp *http.Response) (string, error){
	path := xmlpath.MustCompile("//input[@class='sls']/@value")
	root, err := xmlpath.ParseHTML(resp.Body)
	if err != nil{
		return "", err
	}
	results := make([]string, 0, 5)
	iter := path.Iter(root)
	for iter.Next(){
		results = append(results, iter.Node().String())
	}
	if len(results) >= 2{
		return results[1], nil
	}
	return "", errors.New("Unable to parse html.")
}
