package transport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/Dizao9/Fitness-Journal/internal/transport/dto"
	"github.com/google/uuid"
)

type WorkoutService interface {
	CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error)
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
	var workout dto.CreateWorkoutReq
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, "invalid body request", http.StatusBadRequest)
		return
	}

	workoutDomain := domain.Workout{
		AthleteID:       userID,
		GradeOfTraining: workout.GradeOfTraining,
	}

	switch workout.Status {
	case domain.WorkoutStatusFinished:
		workoutDomain.Status = domain.WorkoutStatusFinished
	case domain.WorkoutStatusInProgress:
		workoutDomain.Status = domain.WorkoutStatusInProgress
	default:
		http.Error(w, "invalid workout status", http.StatusBadRequest)
		return
	}
	workoutDomain.TotalTime = workout.TotalTime

	workoutDomain.DateOfTraining = workout.DateOfTraining

	workoutDomain.Sets = make([]domain.Set, 0, len(workout.Sets))
	for _, v := range workout.Sets {
		workoutDomain.Sets = append(workoutDomain.Sets, domain.Set{
			ExerciseID: v.ExerciseID,
			Weight:     v.Weight,
			SetOrder:   v.SetOrder,
			Reps:       v.Reps,
			Rpe:        v.Rpe,
		})
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id, err := h.Service.CreateWorkout(ctx, workoutDomain)
	if err != nil {
		log.Printf("[CREATE_WORKOUT] internal server error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": id.String()}); err != nil {
		log.Printf("[CREATE_WORKOUT] Encoder failed :%v", err)
	}
}
