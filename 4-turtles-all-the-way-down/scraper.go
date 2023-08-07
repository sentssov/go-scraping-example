package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	c := colly.NewCollector(
		colly.Async(false))

	items := []Item{}

	c.OnHTML("iframe", func(h *colly.HTMLElement) {
		link := h.Request.AbsoluteURL(h.Attr("src"))
		c.Request("GET", link, nil, nil, nil)
	})

	c.OnHTML("div.turtle-family-card", func(h *colly.HTMLElement) {
		link := h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
		c.Request("GET", link, nil, nil, nil)
	})

	c.OnHTML("div.turtle-family-detail", func(h *colly.HTMLElement) {
		i := Item{
			Name:        h.ChildText("h3"),
			Description: h.ChildText("p.lead"),
		}
		items = append(items, i)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://www.scrapethissite.com/pages/frames/")

	bytes, err := json.MarshalIndent(items, " ", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bytes))
}
