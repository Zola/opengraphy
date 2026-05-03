package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"opengraphy/internal/cache"
	"opengraphy/internal/config"
	"opengraphy/internal/fetcher"
	"opengraphy/internal/gallery"
	"opengraphy/internal/i18n"
	"opengraphy/internal/model"
	"opengraphy/internal/parser"
	"opengraphy/internal/preview"
	"opengraphy/internal/security"
	"opengraphy/internal/stats"
)

type App struct {
	Config  config.Config
	Cache   *cache.Redis
	Stats   *stats.Stats
	I18n    *i18n.Bundle
	Fetcher *fetcher.Fetcher
	Parser  *parser.Parser
	Gallery *gallery.Store

	templates *template.Template
}

type PageData struct {
	Template      string
	Title         string
	Description   string
	CanonicalURL  string
	DefaultURL    string
	AlternateURLs map[string]string
	NoIndex       bool
	Locale        string
	Locales       []string
	LocaleLabels  map[string]string
	T             func(string) string
	DiagMessage   func(model.Diagnostic) string
	DiagFix       func(model.Diagnostic) string
	FormatBytes   func(int64) string
	ImageSize     func(model.ImageInfo) string
	URL           string
	Result        model.CheckResult
	Stats         stats.Snapshot
	Recent        []model.Work
	Random        []model.Work
	Leaderboard   []model.Work
	Error         string
}

func (a *App) Routes() http.Handler {
	funcs := template.FuncMap{
		"safe": func(s string) template.HTML { return template.HTML(s) },
	}
	a.templates = template.Must(template.New("").Funcs(funcs).ParseGlob(filepath.Join("internal", "templates", "*.html")))

	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.RequestID, middleware.Recoverer, middleware.Compress(5))
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))).ServeHTTP(w, r)
	})

	r.Get("/", a.home)
	r.Get("/gallery", a.gallery)
	r.Post("/check", a.submitCheck)
	r.Get("/preview", a.previewPage)
	r.Get("/panel", a.previewPage)
	r.Get("/api/check", a.apiCheck)
	r.Post("/api/check", a.apiCheck)
	r.Get("/api/stats", a.apiStats)
	r.Post("/api/presence", a.apiPresence)
	r.Post("/api/lang", a.apiLang)
	r.Get("/api/gallery/recent", a.apiRecent)
	r.Get("/api/gallery/random", a.apiRandom)
	r.Get("/api/gallery/leaderboard", a.apiLeaderboard)
	r.Get("/api/pk/pair", a.apiPKPair)
	r.Post("/api/pk/vote", a.apiPKVote)
	r.Get("/privacy", a.staticPage("privacy.html", "Privacy"))
	r.Get("/terms", a.staticPage("terms.html", "Terms"))
	r.Get("/robots.txt", a.robots)
	r.Get("/sitemap.xml", a.sitemap)
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	return r
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	locale := a.locale(r)
	a.Stats.Visit(r.Context(), locale)
	a.render(w, r, "home.html", PageData{
		Title:       "OpenGraphy",
		Stats:       a.Stats.Snapshot(r.Context()),
		Recent:      a.Gallery.Recent(r.Context(), 8),
		Leaderboard: a.Gallery.Leaderboard(r.Context(), 6),
	})
}

func (a *App) gallery(w http.ResponseWriter, r *http.Request) {
	locale := a.locale(r)
	a.render(w, r, "gallery.html", PageData{
		Title:       a.I18n.T(locale, "gallery") + " - OpenGraphy",
		Recent:      a.Gallery.Recent(r.Context(), 8),
		Random:      a.Gallery.Random(r.Context(), 20),
		Leaderboard: a.Gallery.Leaderboard(r.Context(), 20),
	})
}

func (a *App) submitCheck(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	http.Redirect(w, r, "/preview?url="+urlQueryEscape(r.FormValue("url")), http.StatusSeeOther)
}

func (a *App) previewPage(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	result, err := a.check(r.Context(), r, rawURL, r.URL.Query().Get("refresh") == "1")
	data := PageData{Title: "Preview - OpenGraphy", URL: rawURL, NoIndex: true}
	if err != nil {
		data.Error = err.Error()
	} else {
		data.Result = result
	}
	a.render(w, r, "preview.html", data)
}

func (a *App) apiCheck(w http.ResponseWriter, r *http.Request) {
	rawURL := r.URL.Query().Get("url")
	if rawURL == "" && r.Method == http.MethodPost {
		var body struct {
			URL string `json:"url"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		rawURL = body.URL
	}
	result, err := a.check(r.Context(), r, rawURL, r.URL.Query().Get("refresh") == "1")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (a *App) check(ctx context.Context, r *http.Request, rawURL string, refresh bool) (model.CheckResult, error) {
	sessionID := a.ensureSession(nil, r)
	ipHash := security.HashWithSalt(clientIP(r), a.Config.RateSalt)
	if !stats.AllowFixedWindow(ctx, a.Cache.Client(), "rl:s:"+sessionID, 10, time.Minute) ||
		!stats.AllowFixedWindow(ctx, a.Cache.Client(), "rl:ip:"+ipHash, 30, time.Minute) {
		return model.CheckResult{}, http.ErrHandlerTimeout
	}

	normalized, _, err := security.Normalize(rawURL)
	if err != nil {
		return model.CheckResult{}, err
	}
	key := "og:metadata:" + security.Hash(normalized)
	var cached model.CheckResult
	if !refresh {
		if ok, err := a.Cache.GetJSON(ctx, key, &cached); err == nil && ok {
			cached.Cache.Hit = true
			cached.Cache.TTLSeconds = a.Cache.TTLSeconds()
			a.Stats.CacheHit(ctx, true)
			return cached, nil
		}
	}
	a.Stats.CacheHit(ctx, false)
	a.Stats.Check(ctx)

	htmlResult, err := a.Fetcher.FetchHTML(ctx, normalized)
	if err != nil {
		return model.CheckResult{}, err
	}
	meta, raw, err := a.Parser.Parse(htmlResult.Body, htmlResult.FinalURL)
	if err != nil {
		return model.CheckResult{}, err
	}
	previewModel := meta.Preview(htmlResult.FinalURL)
	image := a.Fetcher.InspectImage(ctx, previewModel.Image, htmlResult.FinalURL)
	score, diagnostics, generated := preview.Build(meta, image, htmlResult.Info, htmlResult.FinalURL)

	result := model.CheckResult{
		URL:         htmlResult.NormalizedURL,
		FinalURL:    htmlResult.FinalURL,
		Cache:       model.CacheInfo{Hit: false, TTLSeconds: a.Cache.TTLSeconds()},
		Fetch:       htmlResult.Info,
		Meta:        meta,
		Image:       image,
		Diagnostics: diagnostics,
		Score:       score,
		Generated:   generated,
		Preview:     previewModel,
		CheckedAt:   time.Now().UTC(),
		Raw:         raw,
	}
	_ = a.Cache.SetJSON(ctx, key, result)
	_ = a.Gallery.Upsert(ctx, result)
	return result, nil
}

func (a *App) apiStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.Stats.Snapshot(r.Context()))
}

func (a *App) apiPresence(w http.ResponseWriter, r *http.Request) {
	id := a.ensureSession(w, r)
	a.Stats.Presence(r.Context(), id)
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func (a *App) apiLang(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	locale := r.FormValue("lang")
	if !a.I18n.Has(locale) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unsupported locale"})
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "lang", Value: locale, Path: "/", MaxAge: 365 * 24 * 3600, SameSite: http.SameSiteLaxMode})
	if !strings.Contains(r.Header.Get("Accept"), "application/json") {
		target := r.Referer()
		if target == "" {
			target = "/"
		}
		http.Redirect(w, r, target, http.StatusSeeOther)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func (a *App) apiRecent(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.Gallery.Recent(r.Context(), 8))
}

func (a *App) apiRandom(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.Gallery.Random(r.Context(), 20))
}

func (a *App) apiLeaderboard(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.Gallery.Leaderboard(r.Context(), 20))
}

func (a *App) apiPKPair(w http.ResponseWriter, r *http.Request) {
	left, right, err := a.Gallery.PKPair(r.Context())
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]model.Work{"left": left, "right": right})
}

func (a *App) apiPKVote(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Winner string `json:"winner"`
		Loser  string `json:"loser"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
		return
	}
	if err := a.Gallery.Vote(r.Context(), body.Winner, body.Loser); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	left, right, err := a.Gallery.PKPair(r.Context())
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]model.Work{"left": left, "right": right})
}

func (a *App) staticPage(name, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := a.locale(r)
		if name == "privacy.html" {
			title = a.I18n.T(locale, "privacy_title")
		}
		if name == "terms.html" {
			title = a.I18n.T(locale, "terms_title")
		}
		a.render(w, r, name, PageData{Title: title + " - OpenGraphy"})
	}
}

func (a *App) robots(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("User-agent: *\nAllow: /\nSitemap: " + strings.TrimRight(a.Config.PublicBaseURL, "/") + "/sitemap.xml\n"))
}

func (a *App) sitemap(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	base := strings.TrimRight(a.Config.PublicBaseURL, "/")
	paths := []string{"/", "/gallery", "/privacy", "/terms"}
	locales := a.I18n.Locales()
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">`)
	for _, path := range paths {
		for _, locale := range locales {
			loc := a.localizedURL(base, path, locale)
			b.WriteString(`<url><loc>`)
			b.WriteString(xmlEscape(loc))
			b.WriteString(`</loc>`)
			for _, alternate := range locales {
				b.WriteString(`<xhtml:link rel="alternate" hreflang="`)
				b.WriteString(xmlEscape(alternate))
				b.WriteString(`" href="`)
				b.WriteString(xmlEscape(a.localizedURL(base, path, alternate)))
				b.WriteString(`"/>`)
			}
			b.WriteString(`<xhtml:link rel="alternate" hreflang="x-default" href="`)
			b.WriteString(xmlEscape(a.absoluteURL(base, path)))
			b.WriteString(`"/></url>`)
		}
	}
	b.WriteString(`</urlset>`)
	_, _ = w.Write([]byte(b.String()))
}

func (a *App) render(w http.ResponseWriter, r *http.Request, tmpl string, data PageData) {
	locale := a.locale(r)
	base := strings.TrimRight(a.Config.PublicBaseURL, "/")
	data.Template = tmpl
	data.Locale = locale
	data.Locales = a.I18n.Locales()
	data.LocaleLabels = a.I18n.LocaleLabels()
	data.T = func(key string) string { return a.I18n.T(locale, key) }
	data.DiagMessage = func(d model.Diagnostic) string { return a.diagnosticText(locale, d, "message") }
	data.DiagFix = func(d model.Diagnostic) string { return a.diagnosticText(locale, d, "fix") }
	data.FormatBytes = formatBytes
	data.ImageSize = func(image model.ImageInfo) string {
		if image.Width <= 0 || image.Height <= 0 {
			return a.I18n.T(locale, "unknown_label")
		}
		return fmt.Sprintf("%d x %d", image.Width, image.Height)
	}
	if data.Description == "" {
		data.Description = a.I18n.T(locale, "meta_description")
	}
	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	if data.CanonicalURL == "" {
		if r.URL.Query().Get("lang") != "" {
			data.CanonicalURL = a.localizedURL(base, path, locale)
		} else {
			data.CanonicalURL = a.absoluteURL(base, path)
		}
	}
	if data.DefaultURL == "" {
		data.DefaultURL = a.absoluteURL(base, path)
	}
	if !data.NoIndex && data.AlternateURLs == nil {
		data.AlternateURLs = make(map[string]string, len(data.Locales))
		for _, alternate := range data.Locales {
			data.AlternateURLs[alternate] = a.localizedURL(base, path, alternate)
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.templates.ExecuteTemplate(w, tmpl, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) diagnosticText(locale string, d model.Diagnostic, suffix string) string {
	if d.Code == "" {
		if suffix == "fix" {
			return d.Fix
		}
		return d.Message
	}
	key := "diag_" + d.Code + "_" + suffix
	value := a.I18n.T(locale, key)
	if value == key {
		if suffix == "fix" {
			return d.Fix
		}
		return d.Message
	}
	return value
}

func (a *App) absoluteURL(base, path string) string {
	if path == "/" {
		return base + "/"
	}
	return base + path
}

func (a *App) localizedURL(base, path, locale string) string {
	return a.absoluteURL(base, path) + "?lang=" + url.QueryEscape(locale)
}

func xmlEscape(value string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&apos;")
	return replacer.Replace(value)
}

func formatBytes(n int64) string {
	if n <= 0 {
		return "0 B"
	}
	units := []string{"B", "KB", "MB"}
	value := float64(n)
	for _, unit := range units {
		if value < 1024 || unit == units[len(units)-1] {
			if unit == "B" {
				return fmt.Sprintf("%d %s", n, unit)
			}
			return fmt.Sprintf("%.1f %s", value, unit)
		}
		value /= 1024
	}
	return fmt.Sprintf("%d B", n)
}

func (a *App) locale(r *http.Request) string {
	return a.I18n.Locale(r)
}

func (a *App) ensureSession(w http.ResponseWriter, r *http.Request) string {
	if c, err := r.Cookie("og_session"); err == nil && c.Value != "" {
		return c.Value
	}
	if w == nil {
		return "anon-" + security.HashWithSalt(clientIP(r), a.Config.RateSalt)
	}
	id := randomID()
	http.SetCookie(w, &http.Cookie{Name: "og_session", Value: id, Path: "/", MaxAge: 30 * 60, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	return id
}

func randomID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func urlQueryEscape(s string) string {
	return url.QueryEscape(strings.TrimSpace(s))
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
