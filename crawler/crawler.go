package crawler

import (
	"log"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

type Crawler struct {
	Collector *colly.Collector
}

type Results struct {
	RawHTML              string
	URL                  string
	Title                string
	Summary              string
	Author               string
	Timestamp            uint64
	Domain               string
	Country              string
	Lang                 string
	Type                 string
	RelatedInternalLinks []string
	RelatedExternalLinks []string
	Tokens               []string
}

func (c *Crawler) Init() {
	c.Collector = colly.NewCollector(colly.Async(), colly.AllowURLRevisit())
	c.Collector.Limit(&colly.LimitRule{
		Parallelism: 100,
	})

	extensions.RandomUserAgent(c.Collector)
	extensions.Referer(c.Collector)
}

func (c *Crawler) Raw(url string) (string, error) {
	raw := ""
	c.Collector.OnHTML("html", func(h *colly.HTMLElement) {
		html, err := h.DOM.Html()
		if err != nil {
			log.Println(err.Error())
		}
		raw = html
	})

	c.Collector.Visit(url)
	c.Collector.Wait()
	return raw, nil
}

func (c *Crawler) Full(url string, r *Results) error {
	// --- URL ---
	r.URL = url
	c.Collector.OnHTML("html", func(h *colly.HTMLElement) {
		// --- RAW HTML ---
		html, err := h.DOM.Html()
		if err != nil {
			log.Println(err.Error())
		}
		r.RawHTML = html
		// --- TITLE ---
		r.Title = h.DOM.Clone().Find("title").Text()
		// --- Summary ---
		h.DOM.Clone().Find("meta[itemprop=description][content]").Attr("content")
	})

	c.Collector.Visit(url)
	c.Collector.Wait()
	return nil
}

var C *Crawler
