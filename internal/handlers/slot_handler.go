// Package handlers содержит HTTP-обработчики для консультационных слотов
// преподавателя и записей студентов.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"consultation-service/internal/models"
	"consultation-service/internal/repository"

	"github.com/gorilla/mux"
)

// SlotHandler implements the "Создать консультационный слот" and
// "Просмотреть список записавшихся студентов" use cases (teacher side).
type SlotHandler struct {
	repo *repository.SlotRepository
}

// NewSlotHandler creates a SlotHandler backed by the given repository.
func NewSlotHandler(repo *repository.SlotRepository) *SlotHandler {
	return &SlotHandler{repo: repo}
}

// createSlotRequest is the JSON body expected by Create.
type createSlotRequest struct {
	TeacherID    uint      `json:"teacher_id"`    // идентификатор преподавателя-владельца слота
	DisciplineID uint      `json:"discipline_id"` // идентификатор дисциплины
	StartsAt     time.Time `json:"starts_at"`     // время начала консультации
	EndsAt       time.Time `json:"ends_at"`       // время окончания консультации
	Capacity     int       `json:"capacity"`      // число студентов, вмещаемых слотом; <=0 трактуется как 1
}

// Create handles POST /slots — создание преподавателем нового
// консультационного слота. Request body is JSON (see createSlotRequest).
func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный формат запроса")
		return
	}

	if req.Capacity <= 0 {
		req.Capacity = 1
	}

	slot := &models.ConsultationSlot{
		TeacherID:    req.TeacherID,
		DisciplineID: req.DisciplineID,
		StartsAt:     req.StartsAt,
		EndsAt:       req.EndsAt,
		Capacity:     req.Capacity,
	}

	if err := h.repo.CreateSlot(slot); err != nil {
		if errors.Is(err, repository.ErrSlotOverlap) {
			writeError(w, http.StatusConflict, "слот пересекается с существующим у этого преподавателя")
			return
		}
		writeError(w, http.StatusInternalServerError, "не удалось создать слот")
		return
	}

	writeJSON(w, http.StatusCreated, slot)
}

// ListByTeacher handles GET /teachers/{id}/slots — список слотов
// преподавателя (path-параметр "id" — идентификатор преподавателя).
func (h *SlotHandler) ListByTeacher(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "некорректный id преподавателя")
		return
	}

	slots, err := h.repo.ListByTeacher(uint(id))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "не удалось получить слоты")
		return
	}

	writeJSON(w, http.StatusOK, slots)
}
