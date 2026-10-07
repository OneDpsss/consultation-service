// Command server is the entrypoint of the consultation-booking HTTP
// service: it loads configuration, opens the database, builds the route
// table and starts listening.
package main

import (
	"log"
	"net/http"

	"consultation-service/internal/config"
	"consultation-service/internal/handlers"
	"consultation-service/internal/repository"
)

// main loads Config, connects to PostgreSQL (running GORM auto-migration
// as a side effect of repository.NewDB), builds the router and blocks
// on http.ListenAndServe. Any setup failure is fatal — the service has
// nothing useful to do without a database or a listening socket.
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
