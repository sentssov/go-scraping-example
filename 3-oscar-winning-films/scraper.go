package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

type Item struct {
	Title       string `json:"title"`
	Nominations int    `json:"nominations"`
	Awards      int    `json:"awards"`
	BestPicture bool   `json:"best_picture"`
}

func main() {
	c := colly.NewCollector(
		colly.Async(false))

	var items []Item

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		if r.Headers.Get("Content-Type") == "application/json" {
			var el []Item
			if err := json.Unmarshal(r.Body, &el); err != nil {
				log.Fatal(err)
			}
			for _, v := range el {
				items = append(items, v)
			}
		}
	})

	c.OnHTML("a.year-link", func(h *colly.HTMLElement) {
		link := fmt.Sprintf("http://www.scrapethissite.com/pages/ajax-javascript/?ajax=true&year=%s", h.Attr("id"))
		err := c.Post(link, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	c.Visit("http://www.scrapethissite.com/pages/ajax-javascript/")

	bytes, err := json.MarshalIndent(items, " ", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
