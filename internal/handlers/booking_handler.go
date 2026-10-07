package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"consultation-service/internal/repository"

	"github.com/gorilla/mux"
)

type BookingHandler struct {
	repo *repository.BookingRepository
}

func NewBookingHandler(repo *repository.BookingRepository) *BookingHandler {
	return &BookingHandler{repo: repo}
}

type createBookingRequest struct {
	StudentID uint `json:"student_id"`
	SlotID    uint `json:"slot_id"`
}

// Create handles POST /bookings — запись студента на консультацию.
func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный формат запроса")
		return
	}

	booking, err := h.repo.CreateBooking(req.StudentID, req.SlotID)
	if err != nil {
		if errors.Is(err, repository.ErrSlotFull) {
			writeError(w, http.StatusConflict, "слот уже занят")
			return
		}
		writeError(w, http.StatusInternalServerError, "не удалось создать запись")
		return
	}

	writeJSON(w, http.StatusCreated, booking)
}

// Cancel handles DELETE /bookings/{id} — отмена записи студентом.
func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id записи")
		return
	}

	studentIDStr := r.URL.Query().Get("student_id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "не указан student_id")
		return
	}

	if err := h.repo.CancelBooking(uint(id), uint(studentID)); err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось отменить запись")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListByStudent handles GET /students/{id}/bookings.
func (h *BookingHandler) ListByStudent(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id студента")
		return
	}

	bookings, err := h.repo.ListByStudent(uint(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить записи")
		return
	}

	writeJSON(w, http.StatusOK, bookings)
}
