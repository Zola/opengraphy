package i18n

import (
	"net/http"
	"testing"
)

func TestLocaleFallbackFromAcceptLanguage(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "zh-HK,zh;q=0.9,en;q=0.7")
	if got := New().Locale(req); got != "zh-TW" {
		t.Fatalf("locale = %q", got)
	}
}
