package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON writes data as a JSON response body with the given HTTP status.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// errorResponse is the JSON shape returned by writeError.
type errorResponse struct {
	Error string `json:"error"` // сообщение об ошибке, понятное пользователю (ТЗ 3.2.4)
}

// writeError writes a {"error": message} JSON body with the given status.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}
