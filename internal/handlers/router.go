package handlers

import (
	"consultation-service/internal/repository"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// NewRouter wires all handlers to their routes.
func NewRouter(db *gorm.DB) *mux.Router {
	bookingRepo := repository.NewBookingRepository(db)
	bookingHandler := NewBookingHandler(bookingRepo)

	slotRepo := repository.NewSlotRepository(db)
	slotHandler := NewSlotHandler(slotRepo)

	r := mux.NewRouter()

	r.HandleFunc("/bookings", bookingHandler.Create).Methods("POST")
	r.HandleFunc("/bookings/{id}", bookingHandler.Cancel).Methods("DELETE")
	r.HandleFunc("/students/{id}/bookings", bookingHandler.ListByStudent).Methods("GET")

	r.HandleFunc("/slots", slotHandler.Create).Methods("POST")
	r.HandleFunc("/teachers/{id}/slots", slotHandler.ListByTeacher).Methods("GET")

	return r
}
