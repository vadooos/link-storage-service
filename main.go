package main

import (
	"context"
	"link-storage-service/cache"
	"link-storage-service/config"
	"link-storage-service/handler"
	"link-storage-service/repository"
	"link-storage-service/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	linkRepo, err := repository.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open repository: %v", err)
	}
	defer linkRepo.Close()

	var h = handler.New(service.New(linkRepo, cache.New(cfg.CacheTTL), cfg.ShortCodeLength))

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: h.Routes()}

	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown: %v", err)
	}

	log.Println("server stopped")
}
