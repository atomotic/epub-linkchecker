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

			client := http.Client{
				Timeout: time.Duration(30 * time.Second),
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

	if len(os.Args) < 2 {
		fmt.Println("USAGE: epub-linkchecker {file.epub}")
		os.Exit(1)
	}

	r, err := zip.OpenReader(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "html") {

			page, err := f.Open()
			defer page.Close()

			doc, err := goquery.NewDocumentFromReader(page)
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
