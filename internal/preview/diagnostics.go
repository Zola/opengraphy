package preview

import (
	"fmt"
	"strings"

	"opengraphy/internal/model"
)

func Build(meta model.Metadata, image model.ImageInfo, fetch model.FetchInfo, finalURL string) (int, []model.Diagnostic, string) {
	var items []model.Diagnostic
	score := 100
	add := func(sev, field, msg, fix, snippet string, penalty int) {
		items = append(items, model.Diagnostic{Severity: sev, Field: field, Message: msg, Fix: fix, Snippet: snippet})
		score -= penalty
	}

	preview := meta.Preview(finalURL)
	if fetch.Status < 200 || fetch.Status >= 400 {
		add("error", "http", fmt.Sprintf("HTTP status is %d.", fetch.Status), "Make sure the page returns 200 OK for crawlers.", "", 20)
	}
	if fetch.Redirects > 3 {
		add("warning", "redirects", "Redirect chain is long.", "Keep social crawler redirects short and stable.", "", 8)
	}
	if !strings.Contains(strings.ToLower(fetch.ContentType), "html") {
		add("warning", "content-type", "Content type is not clearly HTML.", "Serve text/html for pages that need social previews.", "", 8)
	}
	if preview.Title == "" {
		add("error", "og:title", "Missing title for social previews.", "Add og:title and a regular title tag.", `<meta property="og:title" content="Your page title">`, 16)
	} else if l := len([]rune(preview.Title)); l > 70 {
		add("warning", "title", "Title may be too long.", "Aim for about 40-60 characters.", "", 5)
	} else if l < 12 {
		add("info", "title", "Title is very short.", "Use a descriptive title that still scans quickly.", "", 2)
	}
	if preview.Description == "" {
		add("error", "og:description", "Missing description for social previews.", "Add og:description and meta description.", `<meta property="og:description" content="A useful one or two sentence summary.">`, 16)
	} else if l := len([]rune(preview.Description)); l > 180 {
		add("warning", "description", "Description may be too long.", "Aim for about 120-160 characters.", "", 5)
	}
	if preview.Image == "" {
		add("error", "og:image", "Missing preview image.", "Add an absolute HTTPS og:image URL.", `<meta property="og:image" content="https://example.com/preview.png">`, 20)
	} else {
		if !strings.HasPrefix(preview.Image, "https://") {
			add("warning", "og:image", "Image URL is not HTTPS.", "Use HTTPS image URLs for better crawler compatibility.", "", 5)
		}
		if !image.Reachable {
			add("error", "og:image", "Preview image could not be fetched.", "Check image status, redirects, robots, and content type.", "", 16)
		}
		if image.Width > 0 && image.Height > 0 {
			if image.Width < 1200 || image.Height < 627 {
				add("warning", "og:image", "Preview image is smaller than recommended.", "Use 1200x630 or a 1.91:1 image for broad compatibility.", "", 8)
			}
		}
	}
	if meta.Twitter.Card == "" {
		add("info", "twitter:card", "Twitter/X card type is missing.", "Add twitter:card=summary_large_image for larger cards.", `<meta name="twitter:card" content="summary_large_image">`, 3)
	}
	if score < 0 {
		score = 0
	}
	return score, items, GenerateTags(meta, image, finalURL)
}

func GenerateTags(meta model.Metadata, image model.ImageInfo, finalURL string) string {
	p := meta.Preview(finalURL)
	card := meta.Twitter.Card
	if card == "" {
		card = "summary_large_image"
	}
	width, height := image.Width, image.Height
	if width == 0 && meta.OG.ImageWidth != "" {
		width = atoiLoose(meta.OG.ImageWidth)
	}
	if height == 0 && meta.OG.ImageHeight != "" {
		height = atoiLoose(meta.OG.ImageHeight)
	}
	out := fmt.Sprintf(`<!-- HTML Meta Tags -->
<title>%s</title>
<meta name="description" content="%s">

<!-- Open Graph Meta Tags -->
<meta property="og:url" content="%s">
<meta property="og:type" content="%s">
<meta property="og:title" content="%s">
<meta property="og:description" content="%s">
<meta property="og:image" content="%s">`, esc(p.Title), esc(p.Description), esc(p.URL), first(meta.OG.Type, "website"), esc(p.Title), esc(p.Description), esc(p.Image))
	if width > 0 && height > 0 {
		out += fmt.Sprintf(`
<meta property="og:image:width" content="%d">
<meta property="og:image:height" content="%d">`, width, height)
	}
	out += fmt.Sprintf(`

<!-- Twitter Meta Tags -->
<meta name="twitter:card" content="%s">
<meta name="twitter:title" content="%s">
<meta name="twitter:description" content="%s">
<meta name="twitter:image" content="%s">`, esc(card), esc(p.Title), esc(p.Description), esc(p.Image))
	return out
}

func first(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

func esc(s string) string {
	replacer := strings.NewReplacer(`&`, `&amp;`, `"`, `&quot;`, `<`, `&lt;`, `>`, `&gt;`)
	return replacer.Replace(s)
}

func atoiLoose(s string) int {
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}
