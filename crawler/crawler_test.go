package crawler

import (
	"log"
	"testing"
	"time"
)

func TestFullScrape(t *testing.T) {
	c := new(Crawler)
	c.Init()

	urls := []string{
		"https://stackoverflow.com/questions/68230031/cannot-create-a-string-longer-than-0x1fffffe8-characters-in-json-parse",
		"https://www.epochconverter.com/",
		"https://github.com/trending",
		"https://m3o.com/db",
		"https://news.ycombinator.com/",
		"https://google.co.uk",
		"https://facebook.com",
		"https://www.reddit.com/",
		"https://yahoo.com",
		"https://youtube.com",
		"https://www.scmp.com/news/hong-kong/health-environment/article/3154589/hong-kong-students-showing-signs-depression-new?module=lead_hero_story&pgtype=homepage",
		"https://ai.facebook.com/blog/reskin-a-versatile-replaceable-low-cost-skin-for-ai-research-on-tactile-perception/",
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
		if r.Summary == "" {
			t.Errorf("Failed to scrape summary")
		}

		// --- MAIN CONTENT ---
		if len(r.MainContent) > 500 {
			r.MainContent = r.MainContent[:500]
		}
		log.Printf("Main content: %s", r.MainContent)
		if r.MainContent == "" {
			t.Errorf("Failed to scrape summary")
		}

		// --- TIMESTAMP ---
		log.Printf("Raw Timestamp: %d", r.Timestamp)
		log.Printf("Timestamp: %s", time.Unix(int64(r.Timestamp), 0).Format("02-01-2006 15:04"))
		if r.Timestamp == 0 {
			t.Errorf("Failed to scrape timestamp")
		}

		// --- SITE ---
		log.Printf("Site : %s", r.Site)
		if r.Site == "" {
			t.Errorf("Failed to scrape site")
		}

		// --- COUNTRY ---
		log.Printf("Country : %s", r.Country)
		if r.Country == "" {
			t.Errorf("Failed to scrape country")
		}

		// --- LANGUAGE ---
		log.Printf("Language : %s", r.Lang)
		if r.Lang == "" {
			t.Errorf("Failed to scrape language")
		}

		// --- Links ---
		log.Printf("Internal Links : %v", r.RelatedInternalLinks)
		log.Printf("External Links : %v", r.RelatedExternalLinks)
		if len(r.RelatedExternalLinks) == 0 {
			t.Errorf("Failed to scrape related external links")
		}

		if len(r.RelatedInternalLinks) == 0 {
			t.Errorf("Failed to scrape related internal links")
		}
		log.Println("=================================================")
	}
	t.Log("successfully scraped data")
}
