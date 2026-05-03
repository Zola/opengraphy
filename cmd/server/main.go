package main

import (
	"context"
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
