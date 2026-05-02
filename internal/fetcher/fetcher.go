package fetcher

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"opengraphy/internal/config"
	"opengraphy/internal/model"
	"opengraphy/internal/security"
)

type Fetcher struct {
	cfg config.Config
}

type HTMLResult struct {
	NormalizedURL string
	FinalURL      string
	Body          []byte
	Info          model.FetchInfo
}

func New(cfg config.Config) *Fetcher {
	return &Fetcher{cfg: cfg}
}

func (f *Fetcher) FetchHTML(ctx context.Context, rawURL string) (HTMLResult, error) {
	normalized, u, err := security.Normalize(rawURL)
	if err != nil {
		return HTMLResult{}, err
	}
	if err := security.ValidatePublicHost(ctx, u.Hostname()); err != nil {
		return HTMLResult{}, err
	}

	start := time.Now()
	redirects := 0
	client := &http.Client{
		Timeout: f.cfg.FetchTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirects = len(via)
			if len(via) >= f.cfg.MaxRedirects {
				return errors.New("redirect limit exceeded")
			}
			if err := security.ValidatePublicHost(req.Context(), req.URL.Hostname()); err != nil {
				return err
			}
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, normalized, nil)
	if err != nil {
		return HTMLResult{}, err
	}
	req.Header.Set("User-Agent", f.cfg.BotUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.3")
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return HTMLResult{}, err
	}
	defer resp.Body.Close()

	bodyReader := resp.Body
	if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(resp.Body)
		if err == nil {
			defer gz.Close()
			bodyReader = gz
		}
	}
	limited := io.LimitReader(bodyReader, f.cfg.MaxHTMLBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return HTMLResult{}, err
	}
	if int64(len(body)) > f.cfg.MaxHTMLBytes {
		body = body[:f.cfg.MaxHTMLBytes]
	}
	finalURL := resp.Request.URL.String()
	info := model.FetchInfo{
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		DurationMS:  time.Since(start).Milliseconds(),
		Redirects:   redirects,
		PageBytes:   int64(len(body)),
		Elapsed:     time.Since(start),
	}
	if !isHTML(info.ContentType) {
		info.Error = "response content-type is not HTML"
	}
	return HTMLResult{NormalizedURL: normalized, FinalURL: finalURL, Body: body, Info: info}, nil
}

func (f *Fetcher) InspectImage(ctx context.Context, rawURL string, base string) model.ImageInfo {
	imageURL := ResolveURL(base, rawURL)
	info := model.ImageInfo{URL: imageURL}
	if imageURL == "" {
		info.Error = "missing image URL"
		return info
	}
	_, u, err := security.Normalize(imageURL)
	if err != nil {
		info.Error = err.Error()
		return info
	}
	if err := security.ValidatePublicHost(ctx, u.Hostname()); err != nil {
		info.Error = err.Error()
		return info
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	if err != nil {
		info.Error = err.Error()
		return info
	}
	req.Header.Set("User-Agent", f.cfg.BotUserAgent)
	req.Header.Set("Accept", "image/*,*/*;q=0.2")
	client := &http.Client{Timeout: f.cfg.FetchTimeout}
	resp, err := client.Do(req)
	if err != nil {
		info.Error = err.Error()
		return info
	}
	defer resp.Body.Close()
	info.Status = resp.StatusCode
	info.ContentType = resp.Header.Get("Content-Type")
	info.Bytes = resp.ContentLength
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		info.Error = resp.Status
		return info
	}
	buf, _ := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if info.Bytes < 0 {
		info.Bytes = int64(len(buf))
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(buf))
	if err == nil {
		info.Width = cfg.Width
		info.Height = cfg.Height
	}
	info.Reachable = true
	return info
}

func ResolveURL(base, ref string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	u, err := url.Parse(ref)
	if err == nil && u.IsAbs() {
		return u.String()
	}
	b, err := url.Parse(base)
	if err != nil {
		return ref
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}

func isHTML(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return contentType == "" || strings.Contains(contentType, "text/html") || strings.Contains(contentType, "application/xhtml+xml")
}
