package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"consultation-service/internal/repository"

	"github.com/gorilla/mux"
)

// BookingHandler implements the "Записаться на консультацию", "Отменить
// запись" and student-side listing use cases from п. 3.1.1 ТЗ.
type BookingHandler struct {
	repo *repository.BookingRepository
}

// NewBookingHandler creates a BookingHandler backed by the given repository.
func NewBookingHandler(repo *repository.BookingRepository) *BookingHandler {
	return &BookingHandler{repo: repo}
}

// createBookingRequest is the JSON body expected by Create.
type createBookingRequest struct {
	StudentID uint `json:"student_id"` // идентификатор студента
	SlotID    uint `json:"slot_id"`    // идентификатор консультационного слота
}

// Create handles POST /bookings — запись студента на консультацию.
//
// Возвращает 201 с созданной записью, 400 при некорректном теле запроса,
// 409 если у слота уже нет свободных мест (ErrSlotFull, см. п. 3.2.1 ТЗ —
// запрет двойной записи на один слот).
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

// Cancel handles DELETE /bookings/{id}?student_id=... — отмена студентом
// собственной записи (путевой параметр "id" — идентификатор записи,
// query-параметр "student_id" подтверждает, что отменяет владелец записи).
//
// Возвращает 204 при успехе, 400 если "id" или "student_id" некорректны.
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

// ListByStudent handles GET /students/{id}/bookings — "Просмотр студентом
// списка своих активных и завершенных записей" (путевой параметр "id" —
// идентификатор студента). Каждая запись возвращается вместе со связанным
// слотом (см. BookingRepository.ListByStudent).
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
