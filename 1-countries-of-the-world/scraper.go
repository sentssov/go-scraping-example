package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

type Item struct {
	Name       string `json:"name"`
	Capital    string `json:"capital"`
	Population string `json:"population"`
	Area       string `json:"area"`
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}

func main() {
	defer timer("main")()
	c := colly.NewCollector()

	items := []Item{}

	c.OnHTML("div.col-md-4.country", func(h *colly.HTMLElement) {
		i := Item{
			Name:       h.ChildText("h3"),
			Capital:    h.ChildText("span.country-capital"),
			Population: h.ChildText("span.country-population"),
			Area:       h.ChildText("span.country-area"),
		}
		items = append(items, i)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		bytes, err := json.MarshalIndent(items, " ", "")
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile("items.json", bytes, 0664); err != nil {
			log.Fatal(err)
		}
	})

	c.Visit("http://www.scrapethissite.com/pages/simple/")
}
