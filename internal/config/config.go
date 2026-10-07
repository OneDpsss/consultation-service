// Package config загружает настройки сервиса (адрес, на котором слушать,
// строку подключения к БД) из переменных окружения, подставляя значения
// по умолчанию для локальной разработки, если переменная не задана.
package config

import "os"

// Config — настройки, с которыми запускается сервис.
type Config struct {
	HTTPAddr string // адрес и порт HTTP-сервера (переменная HTTP_ADDR)
	DSN      string // строка подключения к PostgreSQL (переменная DSN)
}

// Load читает HTTP_ADDR и DSN из окружения и собирает Config, подставляя
// вместо незаданных переменных значения по умолчанию (":8080" и DSN для
// локального PostgreSQL).
func Load() Config {
	return Config{
		HTTPAddr: getEnv("HTTP_ADDR", ":8080"),
		DSN:      getEnv("DSN", "host=localhost user=postgres password=postgres dbname=consultations port=5432 sslmode=disable"),
	}
}

// getEnv возвращает значение переменной окружения key или fallback,
// если переменная не задана или пуста.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
