package crawler

import (
	"log"
	"testing"
)

func TestFullScrape(t *testing.T) {
	c := new(Crawler)
	c.Init()

	urls := []string{
		"https://stackoverflow.com/questions/68230031/cannot-create-a-string-longer-than-0x1fffffe8-characters-in-json-parse",
		"https://www.epochconverter.com/",
		"https://github.com/trending",
	}
	for _, url := range urls {
		log.Println("=================================================")
		log.Printf("URL: %s", url)

		r := new(Results)
		c.Full(url, r)

		if r == nil {
			t.Errorf("Failed to scrape anything")
		}

		// --- TITLE ---
		log.Printf("Title: %s", r.Title)
		if r.Title == "" {
			t.Errorf("Failed to scrape title")
		}
		// --- SUMMARY ---
		log.Printf("Summary: %s", r.Summary)
		if r.Title == "" {
			t.Errorf("Failed to scrape summary")
		}
		log.Println("=================================================")
	}
	t.Log("successfully scraped data")
}
