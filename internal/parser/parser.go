package parser

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"opengraphy/internal/fetcher"
	"opengraphy/internal/model"
)

type Parser struct{}

func New() *Parser { return &Parser{} }

func (p *Parser) Parse(html []byte, finalURL string) (model.Metadata, map[string]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return model.Metadata{}, nil, err
	}
	raw := map[string]string{}
	meta := model.Metadata{}
	meta.Title = strings.TrimSpace(doc.Find("title").First().Text())

	doc.Find("meta").Each(func(_ int, s *goquery.Selection) {
		key := attr(s, "property")
		if key == "" {
			key = attr(s, "name")
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return
		}
		value := strings.TrimSpace(attr(s, "content"))
		if value == "" {
			return
		}
		lower := strings.ToLower(key)
		if _, exists := raw[lower]; !exists {
			raw[lower] = value
		}
		switch lower {
		case "description":
			meta.Description = value
		case "robots":
			meta.Robots = value
		case "theme-color":
			meta.ThemeColor = value
		case "og:title":
			meta.OG.Title = value
		case "og:description":
			meta.OG.Description = value
		case "og:image", "og:image:secure_url":
			if meta.OG.Image == "" {
				meta.OG.Image = fetcher.ResolveURL(finalURL, value)
			}
			meta.Images = append(meta.Images, fetcher.ResolveURL(finalURL, value))
		case "og:image:width":
			meta.OG.ImageWidth = value
		case "og:image:height":
			meta.OG.ImageHeight = value
		case "og:url":
			meta.OG.URL = fetcher.ResolveURL(finalURL, value)
		case "og:type":
			meta.OG.Type = value
		case "og:site_name":
			meta.OG.SiteName = value
		case "og:locale":
			meta.OG.Locale = value
		case "twitter:card":
			meta.Twitter.Card = value
		case "twitter:title":
			meta.Twitter.Title = value
		case "twitter:description":
			meta.Twitter.Description = value
		case "twitter:image":
			meta.Twitter.Image = fetcher.ResolveURL(finalURL, value)
		case "twitter:url":
			meta.Twitter.URL = fetcher.ResolveURL(finalURL, value)
		case "twitter:site":
			meta.Twitter.Site = value
		}
	})

	doc.Find("link").Each(func(_ int, s *goquery.Selection) {
		rel := strings.ToLower(attr(s, "rel"))
		href := attr(s, "href")
		if href == "" {
			return
		}
		switch {
		case rel == "canonical":
			meta.Canonical = fetcher.ResolveURL(finalURL, href)
		case strings.Contains(rel, "icon") && meta.Favicon == "":
			meta.Favicon = fetcher.ResolveURL(finalURL, href)
		}
	})
	return meta, raw, nil
}

func attr(s *goquery.Selection, name string) string {
	v, _ := s.Attr(name)
	return v
}
