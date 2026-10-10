package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type WorkoutService interface {
	CreateWorkout(ctx context.Context, input service.CreateWorkoutInput) (uuid.UUID, error)
}

type WorkoutHandler struct {
	Service WorkoutService
}

func NewWorkoutHandler(s WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{Service: s}
}

func (h *WorkoutHandler) CreateTraining(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		log.Printf("[CREATE_TRAINING] false from context")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	var workoutCreateReq dto.CreateWorkoutReq
	if err := json.NewDecoder(r.Body).Decode(&workoutCreateReq); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}

	serviceSets := make([]service.CreateSetInput, len(workoutCreateReq.Sets))
	for i, s := range workoutCreateReq.Sets {
		serviceSets[i] = service.CreateSetInput{
			ExerciseID: s.ExerciseID,
			Weight:     s.Weight,
			SetOrder:   s.SetOrder,
			Reps:       s.Reps,
			Rpe:        s.Rpe,
		}
	}

	workoutInput := service.CreateWorkoutInput{
		AthleteID:       userID,
		TotalTime:       workoutCreateReq.TotalTime,
		GradeOfTraining: workoutCreateReq.GradeOfTraining,
		DateOfTraining:  workoutCreateReq.DateOfTraining,
		Status:          workoutCreateReq.Status,
		Sets:            serviceSets,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id, err := h.Service.CreateWorkout(ctx, workoutInput)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidGradeOfTraining) {
			http.Error(w, "invalid grade of training, should be from 1 to 10", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.DateOfTrainingInFuture) {
			http.Error(w, "date of training shouldn't be in the future", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrNoSetsInFinishedWorkout) {
			http.Error(w, "finished workout must be with sets", http.StatusBadRequest)
			return
		}

		if errors.Is(err, domain.ErrNotEnoughPermission) {
			http.Error(w, "not enough permission for this", http.StatusForbidden)
			return
		}
		http.Error(w, "invalid server error", http.StatusInternalServerError)
		log.Printf("invalid server error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id.String()}); err != nil {
		log.Printf("[CREATE_WORKOUT] Encoder failed :%v", err)
	}
}
