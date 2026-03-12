package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type ExerciseService interface {
	CreateCustomExercise(exerReq dto.CreateExerciseRequestDTO) (int, error)
	GetPageOfExercise(ctx context.Context, userID uuid.UUID, filter string, limit, offset int) ([]dto.ExerciseForPageDTO, error)
	GetExerciseByID(userID uuid.UUID, exerciseID int) (dto.ExerciseDTO, error)
	DeleteExerciseByID(userID uuid.UUID, exerciseID int) error
	UpdateExercise(userID uuid.UUID, exerciseID int, exerciseUPD dto.ExerciseUpdateReqDTO) error
}

type ExerciseHandler struct {
	Service ExerciseService
}

func NewExerciseHandler(s ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{Service: s}
}

func (h *ExerciseHandler) PostExercise(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Print("[CREATE_EXERCISE] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	var newExercise dto.CreateExerciseRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&newExercise); err != nil {
		http.Error(w, "something wrong in request", http.StatusBadRequest)
		return
	}
	if newExercise.Name == "" {
		http.Error(w, "fields required", http.StatusBadRequest)
		return
	}

	newExercise.AthleteID = &userID
	id, err := h.Service.CreateCustomExercise(newExercise)
	if err != nil {
		log.Printf("[CREATE_EXERCISE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": strconv.Itoa(id)}); err != nil {
		log.Printf("[CREATE_EXERCISE] encoder was failed: %v", err)
	}
}

func (h *ExerciseHandler) GetPageOfExercise(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		log.Print("[GET_EXERCISES] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit < 0 || limit > 300 {
		limit = 10
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 0 {
		page = 1
	}

	filter := query.Get("filter")

	ListOfExercises, err := h.Service.GetPageOfExercise(context, userID, filter, limit, page)
	if err != nil {
		log.Printf("[GET_EXERCISES internal server error: %v]", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(ListOfExercises); err != nil {
		log.Printf("[GET_EXERCISES] encode problem :%v", err)
	}
}

func (h *ExerciseHandler) GetExerciseByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Printf("[GET_EXERCISE] false from context ")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}

	exerciseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exercise_id format", http.StatusBadRequest)
		return
	}

	exercise, err := h.Service.GetExerciseByID(userID, exerciseID)
	if err != nil {
		if errors.Is(err, domain.ErrExerciseNotFound) {
			http.Error(w, "invalid exercise_id", http.StatusBadRequest)
			return
		}
		log.Printf("[GET_EXERCISE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(exercise); err != nil {
		log.Printf("[GET_EXERCISE] encode problem :%v", err)
	}
}

func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Printf("[DELETE_EXERCISE] false from context")
		http.Error(w, "missing user_id from middleware", http.StatusInternalServerError)
		return
	}
	exerciseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid exerciseID format", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteExerciseByID(userID, exerciseID); err != nil {
		if errors.Is(err, domain.ErrNotEnoughPermission) {
			http.Error(w, "not enough permission to delete this", http.StatusForbidden)
			return
		}
		if errors.Is(err, domain.ErrExerciseNotFound) {
			http.Error(w, "resource not found", http.StatusNotFound)
			return
		}
		log.Printf("[DELETE_EXERCISE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Printf("[PUT_EXERCISE] false from context")
		http.Error(w, "missing user_id from middleware", http.StatusInternalServerError)
		return
	}

	var exerciseUPD dto.ExerciseUpdateReqDTO
	if err := json.NewDecoder(r.Body).Decode(&exerciseUPD); err != nil {
		http.Error(w, "something wrong in request", http.StatusBadRequest)
		return
	}

	exerciseID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id in path value", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateExercise(userID, exerciseID, exerciseUPD)
	if err != nil {
		if err == domain.ErrExerciseNotFound {
			http.Error(w, "nothing is changed", http.StatusNoContent) //need to know !
			return
		}
		log.Printf("[PUT_EXERCISE] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "ok"}); err != nil {
		log.Printf("[PUT_EXERCISE] encoder is failed: %v", err)
	}
}
