package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/domain"
	"github.com/google/uuid"
)

type mockWorkoutStorage struct {
	insertWorkoutFunc  func(ctx context.Context, workout domain.Workout) (uuid.UUID, error)
	insertSetsBulkFunc func(ctx context.Context, workoutID uuid.UUID, sets []domain.Set) error
}

type mockTxManager struct {
}

func (m *mockTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type mockExerciseWorkStorage struct {
	getExerciseByIDsFunc func(ctx context.Context, userID uuid.UUID, IDs []int) ([]domain.Exercise, error)
}

func (m *mockExerciseWorkStorage) GetExerciseByIDs(ctx context.Context, userID uuid.UUID, IDs []int) ([]domain.Exercise, error) {
	if m.getExerciseByIDsFunc != nil {
		return m.getExerciseByIDsFunc(ctx, userID, IDs)
	}

	return nil, nil
}

func (m *mockWorkoutStorage) InsertWorkout(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
	if m.insertWorkoutFunc != nil {
		return m.insertWorkoutFunc(ctx, workout)
	}

	return uuid.Nil, nil
}

func (m *mockWorkoutStorage) InsertSetsBulk(ctx context.Context, workoutID uuid.UUID, sets []domain.Set) error {
	if m.insertSetsBulkFunc != nil {
		return m.insertSetsBulkFunc(ctx, workoutID, sets)
	}
	return nil
}

func TestWorkoutService_CreateWorkout_Success(t *testing.T) {
	expectedID := uuid.New()

	mockStorage := &mockWorkoutStorage{
		insertWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
			if workout.Status != domain.WorkoutStatusInProgress {
				t.Errorf("expected status %s, got %s", domain.WorkoutStatusInProgress, workout.Status)
			}
			return expectedID, nil
		},
	}

	service := NewWorkoutService(mockStorage, &mockExerciseWorkStorage{}, &mockTxManager{})
	workout := CreateWorkoutInput{
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
		{"grade 0 valid", 0, false},
		{"grade 11 invalid", 11, true},
		{"grade -1 invalid", -1, true},
	}
	for _, tc := range tests {
		expectedID := uuid.Nil
		mockStorage := &mockWorkoutStorage{
			insertWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
				return expectedID, nil
			},
		}

		service := NewWorkoutService(mockStorage, &mockExerciseWorkStorage{}, &mockTxManager{})

		grade := tc.grade
		workout := CreateWorkoutInput{
			AthleteID:       uuid.New(),
			Status:          domain.WorkoutStatusInProgress,
			GradeOfTraining: &grade,
		}

		_, err := service.CreateWorkout(context.Background(), workout)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != domain.ErrInvalidGradeOfTraining {
				t.Fatalf("expected error: %v, got :%v", domain.ErrInvalidGradeOfTraining, err)
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error :%v", err)
			}
		}

	}
}

func intPtr(v int) *int {
	return &v
}

func TestWorkoutService_CreateWorkout_FutureDate(t *testing.T) {
	workServ := NewWorkoutService(&mockWorkoutStorage{}, &mockExerciseWorkStorage{}, &mockTxManager{})
	input := CreateWorkoutInput{AthleteID: uuid.New(), TotalTime: nil, DateOfTraining: time.Now().Add(100 * time.Second), GradeOfTraining: nil, Status: domain.WorkoutStatusFinished,
		Sets: []CreateSetInput{{
			ExerciseID: 1,
			Weight:     110.0,
			SetOrder:   intPtr(1),
			Reps:       8,
			Rpe:        nil,
		}, {
			ExerciseID: 1,
			Weight:     120.0,
			SetOrder:   intPtr(2),
			Reps:       5,
			Rpe:        nil,
		}}}

	_, err := workServ.CreateWorkout(context.Background(), input)
	if err == nil {
		t.Fatalf("ожидалась ошибка о времени в будущем. Не получено ошибки")
	}
	if !errors.Is(err, domain.DateOfTrainingInFuture) {
		t.Fatalf("получена неожиданная ошибка: %v", err)
	}
}

func TestWorkoutService_CreateWorkout_NothingInSets(t *testing.T) {
	workServ := NewWorkoutService(&mockWorkoutStorage{}, &mockExerciseWorkStorage{}, &mockTxManager{})
	input := CreateWorkoutInput{AthleteID: uuid.New(), TotalTime: nil, DateOfTraining: time.Now(), GradeOfTraining: nil, Status: domain.WorkoutStatusFinished,
		Sets: nil,
	}

	_, err := workServ.CreateWorkout(context.Background(), input)
	if err == nil {
		t.Fatalf("ожидалась ошибка ErrNoSetsInFinishedWorkout")
	}
	if !errors.Is(err, domain.ErrNoSetsInFinishedWorkout) {
		t.Fatalf("получена неожиданная ошибка: %v", err)
	}
}
