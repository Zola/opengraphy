package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port          string
	PublicBaseURL string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	SQLitePath    string
	BotUserAgent  string
	CacheTTL      time.Duration
	GalleryTTL    time.Duration
	FetchTimeout  time.Duration
	MaxHTMLBytes  int64
	MaxRedirects  int
	RateSalt      string
	AdminToken    string
}

func Load() Config {
	return Config{
		Port:          env("PORT", "8080"),
		PublicBaseURL: env("PUBLIC_BASE_URL", "https://opengraphy.example.com"),
		RedisAddr:     env("REDIS_ADDR", "redis:6379"),
		RedisPassword: env("REDIS_PASSWORD", ""),
		RedisDB:       envInt("REDIS_DB", 0),
		SQLitePath:    env("SQLITE_PATH", "/data/opengraphy.db"),
		BotUserAgent:  env("BOT_USER_AGENT", "OpenGraphyBot/1.0 (+https://opengraphy.example.com/bot)"),
		CacheTTL:      envDuration("CACHE_TTL", 4*time.Hour),
		GalleryTTL:    envDuration("GALLERY_TTL", 30*24*time.Hour),
		FetchTimeout:  envDuration("FETCH_TIMEOUT", 8*time.Second),
		MaxHTMLBytes:  int64(envInt("MAX_HTML_BYTES", 2*1024*1024)),
		MaxRedirects:  envInt("MAX_REDIRECTS", 5),
		RateSalt:      env("RATE_SALT", "change-me-in-production"),
		AdminToken:    env("ADMIN_TOKEN", ""),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
