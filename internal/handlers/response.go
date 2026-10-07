package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON записывает data в тело ответа как JSON с указанным HTTP-статусом.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// errorResponse — форма JSON-ответа, которую пишет writeError.
type errorResponse struct {
	Error string `json:"error"` // понятное пользователю сообщение об ошибке
}

// writeError записывает тело {"error": message} с указанным статусом.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
