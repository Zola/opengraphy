package model

import "time"

type CheckResult struct {
	URL         string            `json:"url"`
	FinalURL    string            `json:"final_url"`
	Cache       CacheInfo         `json:"cache"`
	Fetch       FetchInfo         `json:"fetch"`
	Meta        Metadata          `json:"meta"`
	Image       ImageInfo         `json:"image"`
	Diagnostics []Diagnostic      `json:"diagnostics"`
	Score       int               `json:"score"`
	Generated   string            `json:"generated"`
	Preview     PreviewModel      `json:"preview"`
	CheckedAt   time.Time         `json:"checked_at"`
	Raw         map[string]string `json:"raw,omitempty"`
}

type CacheInfo struct {
	Hit        bool `json:"hit"`
	TTLSeconds int  `json:"ttl_seconds"`
}

type FetchInfo struct {
	Status      int           `json:"status"`
	ContentType string        `json:"content_type"`
	DurationMS  int64         `json:"duration_ms"`
	Redirects   int           `json:"redirects"`
	PageBytes   int64         `json:"page_bytes"`
	Error       string        `json:"error,omitempty"`
	Elapsed     time.Duration `json:"-"`
}

type Metadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Canonical   string   `json:"canonical"`
	Favicon     string   `json:"favicon"`
	Robots      string   `json:"robots"`
	ThemeColor  string   `json:"theme_color"`
	OG          OGMeta   `json:"og"`
	Twitter     TwMeta   `json:"twitter"`
	Images      []string `json:"images"`
}

type OGMeta struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageWidth  string `json:"image_width"`
	ImageHeight string `json:"image_height"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	SiteName    string `json:"site_name"`
	Locale      string `json:"locale"`
}

type TwMeta struct {
	Card        string `json:"card"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Site        string `json:"site"`
}

type ImageInfo struct {
	URL         string `json:"url"`
	Status      int    `json:"status"`
	ContentType string `json:"content_type"`
	Bytes       int64  `json:"bytes"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Reachable   bool   `json:"reachable"`
	Error       string `json:"error,omitempty"`
}

type Diagnostic struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Field    string `json:"field"`
	Message  string `json:"message"`
	Fix      string `json:"fix"`
	Snippet  string `json:"snippet,omitempty"`
}

type PreviewModel struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	URL         string `json:"url"`
	Domain      string `json:"domain"`
	SiteName    string `json:"site_name"`
}

type Work struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	FinalURL    string    `json:"final_url"`
	Domain      string    `json:"domain"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Score       int       `json:"score"`
	Rating      int       `json:"rating"`
	Wins        int       `json:"wins"`
	Losses      int       `json:"losses"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

func (m Metadata) Preview(finalURL string) PreviewModel {
	title := firstNonEmpty(m.OG.Title, m.Twitter.Title, m.Title)
	desc := firstNonEmpty(m.OG.Description, m.Twitter.Description, m.Description)
	image := firstNonEmpty(m.OG.Image, m.Twitter.Image)
	url := firstNonEmpty(m.OG.URL, m.Twitter.URL, m.Canonical, finalURL)
	return PreviewModel{
		Title:       title,
		Description: desc,
		Image:       image,
		URL:         url,
		Domain:      domainOf(url),
		SiteName:    firstNonEmpty(m.OG.SiteName, domainOf(url)),
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func domainOf(raw string) string {
	for _, prefix := range []string{"https://", "http://"} {
		if len(raw) > len(prefix) && raw[:len(prefix)] == prefix {
			raw = raw[len(prefix):]
			break
		}
	}
	for i, r := range raw {
		if r == '/' || r == '?' || r == '#' {
			return raw[:i]
		}
	}
	return raw
}
