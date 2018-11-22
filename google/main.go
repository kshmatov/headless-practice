package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/raff/godet"
)

type Playlist map[string]string
type PLItem [2]string

const qString = "https://www.google.com/search?q=\"%v\"+смотреть+site:vk.com&tbm=vid&source=lnt&tbs=dur:l&sa=X"

func getURL(movie string) string {
	m := strings.Replace(movie, " ", "+", -1)
	return fmt.Sprintf(qString, m)
}

func handleChrome(url string, ch chan<- PLItem) {
	defer close(ch)

	remote, err := godet.Connect("localhost:9222", true)
	if err != nil {
		return
	}

	remote.NetworkEvents(true)

	remote.CallbackEvent("Network.loadingFinished", func(params godet.Params) {
		rID := params["requestId"].(string)
		if rID != "" {
			res, err := remote.SendRequest("Network.getResponseBody", godet.Params{
				"requestId": rID,
			})

			if err != nil {
				fmt.Printf("get response body: %v\n", err)
				panic(err)
			}
			// bbb := res["body"].(string)
			// fmt.Printf("Resp: %v\n", string([]rune(bbb)[0:30]))
			bd := res["body"]
			if bd == nil {
				panic("Ups!!!")
			}
			body := ""
			if b, ok := res["base64Encoded"]; ok && b.(bool) {
				r, _ := base64.StdEncoding.DecodeString(res["body"].(string))
				body = string(r)
			} else {
				body = res["body"].(string)
			}
			if strings.Index(body, "<html") >= 0 {
				ioutil.WriteFile("respond.html", []byte(body), 0644)
				html, err := goquery.NewDocumentFromReader(strings.NewReader(body))
				if err != nil {
					fmt.Printf("On parse html: %v", err)
					return
				}
				html.Find("div.rc").Each(func(i int, s *goquery.Selection) {
					a := s.Find("a").First()
					videoURL, _ := a.Attr("href")
					comment := s.Find("span.st").First().Text()
					ch <- PLItem{videoURL, comment}
				})
				// if tmpURL, ok := html.Find("#pnnext").First().Attr("href"); ok {
				// 	url = tmpURL
				// 	fmt.Printf("------> new url: %v\n", url)
				// }
			}
		}
	})

	for i := 0; i < 3; i++ {
		nurl := url
		if i > 0 {
			nurl += "&start=" + strconv.Itoa(i*10)
		}
		fmt.Printf("------> url: %v\n", nurl)
		s, err := remote.Navigate(nurl)
		fmt.Printf("%v - %v\n", s, err)
		time.Sleep(time.Second * time.Duration((rand.Intn(15) + 20)))
	}
}

func getPlaylist(url string) (Playlist, error) {
	result := make(chan PLItem)
	pl := Playlist{}

	go handleChrome(url, result)

	for r := range result {
		pl[r[0]] = r[1]
	}
	return pl, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n$ google <movie name>")
		return
	}
	mName := strings.Join(os.Args[1:], " ")
	mSearchURL := getURL(mName)
	fmt.Printf("For search: %v\n", mSearchURL)

	r, err := getPlaylist(mSearchURL)
	if err != nil {
		fmt.Println("Error :" + err.Error())
		return
	}
	for i := range r {
		fmt.Println(i)
	}
}
