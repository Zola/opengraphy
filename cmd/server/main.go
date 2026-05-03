package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"opengraphy/internal/cache"
	"opengraphy/internal/config"
	"opengraphy/internal/fetcher"
	"opengraphy/internal/gallery"
	"opengraphy/internal/handler"
	"opengraphy/internal/i18n"
	"opengraphy/internal/parser"
	"opengraphy/internal/stats"
)

func main() {
	cfg := config.Load()
	if len(os.Args) > 1 {
		runCLI(cfg, os.Args[1])
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisCache := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	defer redisCache.Close()

	galleryStore, err := gallery.NewStore(cfg.SQLitePath, redisCache.Client(), cfg.GalleryTTL)
	if err != nil {
		log.Fatalf("gallery store: %v", err)
	}
	defer galleryStore.Close()
	_ = galleryStore.Sync(ctx)
	_ = galleryStore.SeedDefaultsIfEmpty(ctx)
	go galleryStore.SyncEvery(ctx, time.Hour)

	app := &handler.App{
		Config:  cfg,
		Cache:   redisCache,
		Stats:   stats.New(redisCache.Client()),
		I18n:    i18n.New(),
		Fetcher: fetcher.New(cfg),
		Parser:  parser.New(),
		Gallery: galleryStore,
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           app.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("opengraphy listening on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}

func runCLI(cfg config.Config, command string) {
	ctx := context.Background()
	redisCache := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	defer redisCache.Close()

	galleryStore, err := gallery.NewStore(cfg.SQLitePath, redisCache.Client(), cfg.GalleryTTL)
	if err != nil {
		log.Fatalf("gallery store: %v", err)
	}
	defer galleryStore.Close()

	switch command {
	case "gallery-health":
		writeJSON(galleryStore.Health(ctx))
	case "gallery-sync":
		if err := galleryStore.Sync(ctx); err != nil {
			log.Fatalf("gallery sync: %v", err)
		}
		writeJSON(galleryStore.Health(ctx))
	case "gallery-seed":
		if err := galleryStore.SeedDefaultsIfEmpty(ctx); err != nil {
			log.Fatalf("gallery seed: %v", err)
		}
		writeJSON(galleryStore.Health(ctx))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", command)
		fmt.Fprintln(os.Stderr, "available commands: gallery-health, gallery-sync, gallery-seed")
		os.Exit(2)
	}
}

func writeJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Fatalf("json: %v", err)
	}
}
