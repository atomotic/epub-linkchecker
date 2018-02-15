package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/campoy/unique"
)

var urls []string

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

	for _, url := range urls {
		fmt.Println(url)
	}

}
