package main

import (
	"log"
	"net/http"

	"consultation-service/internal/config"
	"consultation-service/internal/handlers"
	"consultation-service/internal/repository"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDB(cfg.DSN)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	router := handlers.NewRouter(db)

	log.Printf("listening on %s", cfg.HTTPAddr)
	if err := http.ListenAndServe(cfg.HTTPAddr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
