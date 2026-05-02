package preview

import (
	"strings"
	"testing"

	"opengraphy/internal/model"
)

func TestBuildDiagnosticsMissingCoreTags(t *testing.T) {
	score, items, _ := Build(model.Metadata{}, model.ImageInfo{}, model.FetchInfo{Status: 200, ContentType: "text/html"}, "https://example.com")
	if score >= 80 {
		t.Fatalf("score should be low, got %d", score)
	}
	var fields []string
	for _, item := range items {
		fields = append(fields, item.Field)
	}
	joined := strings.Join(fields, ",")
	for _, want := range []string{"og:title", "og:description", "og:image"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing diagnostic %s in %s", want, joined)
		}
	}
}
