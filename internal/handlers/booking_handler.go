package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"consultation-service/internal/repository"

	"github.com/gorilla/mux"
)

// BookingHandler обрабатывает запись студента на консультацию, отмену
// записи и просмотр своих записей.
type BookingHandler struct {
	repo *repository.BookingRepository
}

// NewBookingHandler создаёт BookingHandler поверх переданного репозитория.
func NewBookingHandler(repo *repository.BookingRepository) *BookingHandler {
	return &BookingHandler{repo: repo}
}

// createBookingRequest — тело JSON-запроса, которое ожидает Create.
type createBookingRequest struct {
	StudentID uint `json:"student_id"` // идентификатор студента
	SlotID    uint `json:"slot_id"`    // идентификатор слота консультации
}

// Create обрабатывает POST /bookings — запись студента на консультацию.
//
// Возвращает 201 с созданной записью, 400 при некорректном теле запроса,
// 409 — если в слоте уже нет свободных мест.
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

// Cancel обрабатывает DELETE /bookings/{id}?student_id=... — отмену
// студентом собственной записи. Параметр пути "id" — идентификатор
// записи, параметр запроса "student_id" подтверждает, что отменяет
// именно владелец записи.
//
// Возвращает 204 при успехе, 400 — если "id" или "student_id" некорректны.
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

// ListByStudent обрабатывает GET /students/{id}/bookings — список записей
// студента (параметр пути "id" — идентификатор студента). Каждая запись
// возвращается вместе со связанным слотом.
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
