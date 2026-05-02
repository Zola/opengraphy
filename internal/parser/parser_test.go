package parser

import "testing"

func TestParseMetadataAndResolveRelativeImage(t *testing.T) {
	html := []byte(`<!doctype html><html><head>
<title>Fallback title</title>
<meta name="description" content="Fallback description">
<meta property="og:title" content="OG title">
<meta property="og:description" content="OG description">
<meta property="og:image" content="/preview.png">
<meta name="twitter:card" content="summary_large_image">
<link rel="canonical" href="/article">
</head></html>`)
	meta, _, err := New().Parse(html, "https://example.com/path/page")
	if err != nil {
		t.Fatal(err)
	}
	if meta.OG.Title != "OG title" {
		t.Fatalf("og title = %q", meta.OG.Title)
	}
	if meta.OG.Image != "https://example.com/preview.png" {
		t.Fatalf("image = %q", meta.OG.Image)
	}
	if meta.Canonical != "https://example.com/article" {
		t.Fatalf("canonical = %q", meta.Canonical)
	}
}
