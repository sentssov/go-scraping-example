package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"net/url"
	"os"
	"time"
)

type Item struct {
	TeamName string `json:"team-name"`
	Year     string `json:"year"`
	Wins     string `json:"wins"`
	Losses   string `json:"losses"`
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

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("form[action='/pages/forms/']", func(h *colly.HTMLElement) {
		query := fmt.Sprintf("%s?q=%s", h.Attr("action"), url.QueryEscape(os.Args[1]))
		h.Request.Visit(query)
	})

	c.OnHTML("ul.pagination", func(h *colly.HTMLElement) {
		h.ForEach("li a", func(_ int, a *colly.HTMLElement) {
			c.Visit(h.Request.AbsoluteURL(a.Attr("href")))
		})
	})

	c.OnHTML("table tr.team", func(h *colly.HTMLElement) {
		i := Item{
			TeamName: h.ChildText("td.name"),
			Year:     h.ChildText("td.year"),
			Wins:     h.ChildText("td.wins"),
			Losses:   h.ChildText("td.losses"),
		}
		items = append(items, i)
	})

	c.Visit("http://www.scrapethissite.com/pages/forms/")

	bytes, err := json.MarshalIndent(items, " ", "")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(GenerateString(), bytes, 0664); err != nil {
		log.Fatal(err)
	}
}

func GenerateString() string {
	rBytes := make([]byte, 16)
	_, err := rand.Read(rBytes)
	if err != nil {
		log.Fatal(err)
	}

	rStr := hex.EncodeToString(rBytes)
	return "file_" + rStr + ".json"
}
