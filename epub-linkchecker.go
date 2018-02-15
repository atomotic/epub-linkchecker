package main

import (
	"archive/zip"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/campoy/unique"
)

var urls []string
var wg sync.WaitGroup

func URLCheck(urls []string) {
	statuses := make(chan string)
	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()

			timeout := time.Duration(10 * time.Second)
			client := http.Client{
				Timeout: timeout,
			}
			resp, err := client.Get(url)

			if err != nil {
				statuses <- fmt.Sprintf("%s\tNetwork Error", url)
			} else {
				statuses <- fmt.Sprintf("%s\t%s", url, resp.Status)
			}
		}(url)

	}
	go func() {
		for status := range statuses {
			fmt.Println(status)
		}
	}()

	wg.Wait()

}

func main() {

	r, err := zip.OpenReader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "html") {

			rc, err := f.Open()
			defer rc.Close()

			doc, err := goquery.NewDocumentFromReader(rc)
			if err != nil {
				log.Fatal(err)
			}
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				if strings.HasPrefix(href, "http") {
					urls = append(urls, href)
				}
			})
		}
	}

	less := func(i, j int) bool { return urls[i] != urls[j] }
	unique.Slice(&urls, less)

	URLCheck(urls)
}
