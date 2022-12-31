package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	fmt.Println(getVersion())
}

func getVersion() string {

	resp, err := http.Get("https://line.ar.uptodown.com/android/download")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var version string
	doc.Find(".version").Each(func(i int, s *goquery.Selection) {
		version = s.Text()

	})
	return version

}
