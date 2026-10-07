// Command server — точка входа HTTP-сервиса записи на консультации: он
// загружает конфигурацию, открывает соединение с БД, собирает таблицу
// маршрутов и запускает прослушивание.
package main

import (
	"log"
	"net/http"

	"github.com/OneDpsss/consultation-service/internal/config"
	"github.com/OneDpsss/consultation-service/internal/handlers"
	"github.com/OneDpsss/consultation-service/internal/repository"
)

// main загружает Config, подключается к PostgreSQL (попутно прогоняя
// автомиграцию через repository.NewDB), собирает роутер и блокируется на
// http.ListenAndServe. Любая ошибка на старте фатальна — без БД или
// открытого сокета сервису нечего делать.
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
