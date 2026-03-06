package transport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type ExerciseService interface {
	CreateCustomExercise(exerReq dto.CreateExerciseRequestDTO) (int, error)
	GetPageOfExercise(ctx context.Context, userID uuid.UUID, filter string, limit, offset int) ([]dto.ExerciseForPageDTO, error)
}

type ExerciseHandler struct {
	Service ExerciseService
}

func NewExerciseHandler(s ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{Service: s}
}

func (h *ExerciseHandler) PostExercise(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userIDStr, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Print("[CREATE_EXERCISE] false from context")
		http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Printf("[CREATE_EXERCISE] invalid format userID from token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
	userIDStr, ok := UserIDFromContext(ctx)
	if !ok {
		if !ok {
			log.Print("[GET_EXERCISES] false from context")
			http.Error(w, "missing userID from middleware", http.StatusInternalServerError)
			return
		}
	}

	context, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		log.Printf("[GET_EXERCISES] invalid format userID from token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

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
