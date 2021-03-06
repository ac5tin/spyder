package crawler

import (
	"log"
	URL "net/url"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	collector *colly.Collector
}

type Results struct {
	RawHTML              string
	URL                  string
	Title                string
	Summary              string
	MainContent          string
	Author               string
	Timestamp            uint64
	Site                 string
	Country              string
	Lang                 string
	Type                 string
	RelatedInternalLinks []string
	RelatedExternalLinks []string
	Tokens               []string
}

func (c *Crawler) Init() {
	c.collector = colly.NewCollector(colly.Async(), colly.AllowURLRevisit())
	c.collector.Limit(&colly.LimitRule{
		Parallelism: 100,
	})

	extensions.RandomUserAgent(c.collector)
	extensions.Referer(c.collector)

	c.collector.IgnoreRobotsTxt = true
}

func (c *Crawler) Raw(url string) (string, error) {
	raw := ""
	c.collector.OnHTML("html", func(h *colly.HTMLElement) {
		html, err := h.DOM.Html()
		if err != nil {
			log.Println(err.Error())
		}
		raw = html
	})

	c.collector.Visit(url)
	c.collector.Wait()
	return raw, nil
}

func (c *Crawler) Full(url string, r *Results) error {
	// --- URL ---
	r.URL = url
	c.collector.OnHTML("html", func(h *colly.HTMLElement) {
		// --- RAW HTML ---
		{
			html, err := h.DOM.Html()
			if err != nil {
				log.Println(err.Error())
			}
			r.RawHTML = html
		}
		// --- TITLE ---
		{

			r.Title = h.DOM.Find("title").Text()
		}
		// --- Summary ---
		{

			if v, ok := h.DOM.Find("meta[itemprop=description][content]").Attr("content"); ok {
				r.Summary = v
			}
			if r.Summary == "" {
				if v, ok := h.DOM.Find("meta[name=description][content]").Attr("content"); ok {
					r.Summary = v
				}
			}

			if r.Summary == "" {
				if v, ok := h.DOM.Find("meta[name=cse_summary][content]").Attr("content"); ok {
					r.Summary = v
				}
			}
		}
		// --- MAIN CONTENT ---
		{
			v := h.DOM.Find("main")
			if len(v.Nodes) == 0 {
				v = h.DOM.Find("[id^=content")
			}
			if len(v.Nodes) == 0 {
				v = h.DOM.Find("[id^=main]")
			}
			// final fallback
			if len(v.Nodes) == 0 {
				v = h.DOM.Find("body")
				v.Find("header").Remove()
				v.Find("footer ~ *").Remove()
				v.Find("footer").Remove()
			}
			if len(v.Nodes) > 0 {
				vv := v.Clone()
				// --- CLEAN ---
				// styling
				vv.Find("style").Remove()
				// navigation
				vv.Find("nav").Remove()
				vv.Find("[role=navigation]").Remove()
				// scripts
				vv.Find("script").Remove()
				vv.Find("noscript").Remove()
				// ads
				vv.Find("ins").Remove()
				vv.Find("[data-ad-client]").Remove()
				vv.Find(".ads").Remove()
				vv.Find(".advert").Remove()

				r.MainContent = vv.Text()
			}

		}

		// --- AUTHOR ---
		{
			if v, ok := h.DOM.Find("meta[name=cse_author][content]").Attr("content"); ok {
				r.Author = v
			}
		}
		// --- TIMESTAMP ---
		{
			if v, ok := h.DOM.Find("meta[property$=updated_time][content]").Attr("content"); ok {
				t, err := time.Parse("2006-01-02T15:04:05+-07:00", v)
				if err == nil {
					r.Timestamp = uint64(t.Unix())
				}

			}
		}
		// --- SITE ---
		{
			u, err := URL.Parse(url)
			if err == nil {
				r.Site = u.Host
			}
		}
		// --- LANGUAGE ---
		{
			if v := h.Attr("lang"); v != "" {

				if strings.Contains(v, "en") {
					r.Lang = "en"
					return
				}
				if strings.Contains(v, "zh") {
					r.Lang = "zh"
					return
				}
			}
		}
		// Links
		{
			links := h.DOM.Find("a[href]")
			links.Each(func(_ int, s *goquery.Selection) {
				if href, ok := s.Attr("href"); ok {
					if u, err := URL.Parse(href); err == nil {
						if u.Host == r.Site || u.Host == "" {
							// internal
							r.RelatedInternalLinks = append(r.RelatedInternalLinks, href)
						} else {
							// external
							r.RelatedExternalLinks = append(r.RelatedExternalLinks, href)
						}
					}
				}
			})
		}
	})

	c.collector.Visit(url)
	c.collector.Wait()
	return nil
}

var C *Crawler
