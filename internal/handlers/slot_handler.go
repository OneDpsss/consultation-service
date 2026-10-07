// @file slot_handler.go
// @brief HTTP-обработчики для консультационных слотов преподавателя.

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

// / @brief Обработчик прецедентов "Создать консультационный слот" и
// / "Просмотреть список записавшихся студентов" (со стороны преподавателя).
type SlotHandler struct {
	repo *repository.SlotRepository
}

// / @brief Создаёт SlotHandler.
// / @param repo репозиторий слотов
// / @return указатель на новый SlotHandler
func NewSlotHandler(repo *repository.SlotRepository) *SlotHandler {
	return &SlotHandler{repo: repo}
}

type createSlotRequest struct {
	TeacherID    uint      `json:"teacher_id"`
	DisciplineID uint      `json:"discipline_id"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	Capacity     int       `json:"capacity"`
}

// / @brief Обрабатывает POST /slots — создание преподавателем нового
// / консультационного слота.
// / @param w ответ HTTP
// / @param r запрос HTTP, тело — JSON createSlotRequest
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

// / @brief Обрабатывает GET /teachers/{id}/slots — список слотов
// / преподавателя (основа для просмотра записавшихся студентов).
// / @param w ответ HTTP
// / @param r запрос HTTP, параметр пути "id" — идентификатор преподавателя
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
