package service

import (
	"context"
	"testing"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type mockWorkoutStorage struct {
	createWorkoutFunc func(ctx context.Context, workout domain.Workout) (uuid.UUID, error)
}

func (m *mockWorkoutStorage) CreateWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	return m.createWorkoutFunc(ctx, workout)
}

func TestWorkoutService_CreateWorkout_Success(t *testing.T) {
	expectedID := uuid.New()

	mockStorage := &mockWorkoutStorage{
		createWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
			if workout.Status != domain.WorkoutStatusInProgress {
				t.Errorf("expected status %s, got %s", domain.WorkoutStatusInProgress, workout.Status)
			}
			return expectedID, nil
		},
	}

	service := NewWorkoutService(mockStorage)
	workout := domain.Workout{
		AthleteID: uuid.New(),
		Status:    domain.WorkoutStatusInProgress,
	}

	id, err := service.CreateWorkout(context.Background(), workout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != expectedID {
		t.Errorf("expected id %s got %s", expectedID, id)
	}
}

func TestWorkoutService_CreateWorkout_GradeBoundary(t *testing.T) {
	tests := []struct {
		name    string
		grade   int
		wantErr bool
	}{
		{"grade 1 valid", 1, false},
		{"grade 10 valid", 10, false},
		{"grade 0 invalid", 0, true},
		{"grade 11 invalid", 11, true},
	}
	for _, tc := range tests {
		expectedID := uuid.Nil
		mockStorage := &mockWorkoutStorage{
			createWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
				return expectedID, nil
			},
		}

		service := NewWorkoutService(mockStorage)

		grade := tc.grade
		workout := domain.Workout{
			AthleteID:       uuid.New(),
			Status:          domain.WorkoutStatusInProgress,
			GradeOfTraining: &grade,
		}

		_, err := service.CreateWorkout(context.Background(), workout)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidGrade {
				t.Fatalf("expected error: %v, got :%v", domain.ErrInvalidGrade, err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error :%v", err)
			}
		}

	}
}
