package i18n

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestLocaleFallbackFromAcceptLanguage(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Language", "zh-HK,zh;q=0.9,en;q=0.7")
	if got := New().Locale(req); got != "zh-TW" {
		t.Fatalf("locale = %q", got)
	}
}

func TestLocaleCookieWinsOverQueryAndAcceptLanguage(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/?lang=ja", nil)
	req.Header.Set("Accept-Language", "de;q=1.0")
	req.AddCookie(&http.Cookie{Name: "lang", Value: "zh-TW"})
	if got := New().Locale(req); got != "zh-TW" {
		t.Fatalf("locale = %q", got)
	}
}

func TestLocaleAcceptLanguageWinsBeforeQueryFallback(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/?lang=ja", nil)
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
	if got := New().Locale(req); got != "fr" {
		t.Fatalf("locale = %q", got)
	}
}

func TestTemplateKeysExistForAllLocales(t *testing.T) {
	required := templateKeys(t)
	bundle := New()
	for _, locale := range bundle.Locales() {
		for key := range required {
			if bundle.data[locale][key] == "" {
				t.Fatalf("%s missing translation key %q", locale, key)
			}
		}
	}
}

func templateKeys(t *testing.T) map[string]struct{} {
	t.Helper()
	pattern := regexp.MustCompile(`call \.T "([^"]+)"`)
	keys := map[string]struct{}{}
	files, err := filepath.Glob(filepath.Join("..", "templates", "*.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range pattern.FindAllStringSubmatch(string(raw), -1) {
			keys[match[1]] = struct{}{}
		}
	}
	return keys
}
