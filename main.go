package main

import (
	"link-storage-service/cache"
	"link-storage-service/config"
	"link-storage-service/handler"
	"link-storage-service/repository"
	"link-storage-service/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()
	linkRepo, err := repository.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open repository: %v", err)
	}

	var h = handler.New(service.New(linkRepo, cache.New(cfg.CacheTTL), cfg.ShortCodeLength))

	srv := &http.Server{Addr: ":" + cfg.HTTPPort, Handler: h.Routes()}
	log.Printf("listening on %s", srv.Addr)
	err = srv.ListenAndServe()
	log.Fatalf("server: %v", err)
}
