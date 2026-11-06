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
	getWorkoutsFunc    func(ctx context.Context, athleteID uuid.UUID, limit, offset int) ([]domain.Workout, error)
	getWorkoutByIDFunc func(ctx context.Context, workoutID uuid.UUID) (domain.Workout, error)
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

func (m *mockWorkoutStorage) GetWorkouts(ctx context.Context, athleteID uuid.UUID, limit, offset int) ([]domain.Workout, error) {
	if m.getWorkoutsFunc != nil {
		return m.getWorkoutsFunc(ctx, athleteID, limit, offset)
	}

	return nil, nil
}

func (m *mockWorkoutStorage) GetWorkoutByID(ctx context.Context, workoutID uuid.UUID) (domain.Workout, error) {
	if m.getWorkoutByIDFunc != nil {
		return m.getWorkoutByIDFunc(ctx, workoutID)
	}
	return domain.Workout{}, nil
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
	if !errors.Is(err, domain.ErrDateOfTrainingInFuture) {
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

func TestWorkoutService_CreateWorkout_ExercicesAndSets(t *testing.T) {
	mockExercise := &mockExerciseWorkStorage{
		getExerciseByIDsFunc: func(ctx context.Context, userID uuid.UUID, IDs []int) ([]domain.Exercise, error) {
			return []domain.Exercise{
				{
					ID:          IDs[0],
					Name:        "Bench_press",
					MuscleGroup: "Chest",
					AthleteID:   &userID,
				}, {
					ID:          IDs[1],
					Name:        "Squat",
					MuscleGroup: "Legs",
					AthleteID:   &userID,
				},
			}, nil
		},
	}

	expectedWorkoutID := uuid.New()

	bulkInserted := false
	mockWorkoutStorage := &mockWorkoutStorage{
		insertWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
			return expectedWorkoutID, nil
		},
		insertSetsBulkFunc: func(ctx context.Context, workoutID uuid.UUID, sets []domain.Set) error {
			if workoutID != expectedWorkoutID {
				t.Errorf("ожидался WorkoutID %d, получен %d", expectedWorkoutID, workoutID)
			}

			if len(sets) != 3 {
				t.Errorf("ожидалось 3 подхода, получено %d,", len(sets))
			}
			bulkInserted = true
			return nil
		},
	}

	workServ := NewWorkoutService(mockWorkoutStorage, mockExercise, &mockTxManager{})

	input := CreateWorkoutInput{
		AthleteID:      uuid.New(),
		DateOfTraining: time.Now(),
		Status:         domain.WorkoutStatusFinished,
		Sets: []CreateSetInput{
			{
				ExerciseID: 1,
				Weight:     110,
				Reps:       5,
				Rpe:        intPtr(5),
			},
			{
				ExerciseID: 1,
				Weight:     110,
				Reps:       4,
				Rpe:        intPtr(6),
			},
			{
				ExerciseID: 2,
				Weight:     130,
				Reps:       6,
				Rpe:        intPtr(9),
			},
		},
	}

	id, err := workServ.CreateWorkout(context.Background(), input)
	if err != nil {
		t.Errorf("ошибка не ожидалась, получена %v", err)
	}

	if id != expectedWorkoutID {
		t.Errorf("айди не совпали. Ожидался %s, получен: %s", expectedWorkoutID, id)
	}

	if !bulkInserted {
		t.Errorf("bulk insert не был выполнен")
	}
}

func TestWorkoutService_CreateWorkout_SomeonesElseExercise(t *testing.T) {
	exercStore := &mockExerciseWorkStorage{
		getExerciseByIDsFunc: func(ctx context.Context, userID uuid.UUID, IDs []int) ([]domain.Exercise, error) {
			return []domain.Exercise{}, nil
		},
	}

	isCalled := false

	workStore := &mockWorkoutStorage{
		insertWorkoutFunc: func(ctx context.Context, workout domain.Workout) (uuid.UUID, error) {
			isCalled = true
			return uuid.Nil, nil
		},
	}

	mockWorkoutService := NewWorkoutService(workStore, exercStore, &mockTxManager{})
	input := CreateWorkoutInput{
		AthleteID:      uuid.New(),
		DateOfTraining: time.Now(),
		Status:         domain.WorkoutStatusFinished,
		Sets: []CreateSetInput{
			{
				ExerciseID: 1,
				Weight:     110,
				Reps:       5,
				Rpe:        intPtr(5),
			},
		},
	}

	_, err := mockWorkoutService.CreateWorkout(context.Background(), input)
	if !errors.Is(err, domain.ErrNotEnoughPermission) {
		t.Errorf("Ожидалась ошибка: %v, получена :%v", domain.ErrNotEnoughPermission, err)
	}

	if isCalled {
		t.Errorf("вызов inserWorkout не должен был произойти, но произошёл")
	}
}

func TestWorkoutService_GetWorkouts_EdgeCases(t *testing.T) {
	tests := []struct {
		Name           string
		Limit          int
		Page           int
		ExpectedOffSet int
	}{
		{
			Name:           "Limit is over of rage",
			Limit:          100,
			Page:           2,
			ExpectedOffSet: 50,
		}, {
			Name:           "Limit is lower of rage",
			Limit:          -1,
			Page:           2,
			ExpectedOffSet: 10,
		}, {
			Name:           "Normal",
			Limit:          35,
			Page:           2,
			ExpectedOffSet: 35,
		}, {
			Name:           "Incorrect page",
			Limit:          20,
			Page:           0,
			ExpectedOffSet: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {

			var capturedOffset int

			mockStorage := &mockWorkoutStorage{
				getWorkoutsFunc: func(ctx context.Context, athleteID uuid.UUID, limit int, offset int) ([]domain.Workout, error) {
					capturedOffset = offset
					return nil, nil
				},
			}

			service := NewWorkoutService(mockStorage, &mockExerciseWorkStorage{}, &mockTxManager{})

			_, err := service.GetWorkouts(context.Background(), GetWorkoutsInput{
				AthleteID: uuid.New(),
				Limit:     tc.Limit,
				Page:      tc.Page,
			})

			if err != nil {
				t.Errorf("вызов не должен был привлечь ошибку: %v", err)
			}

			if capturedOffset != tc.ExpectedOffSet {
				t.Errorf("значения ожидаемого смещения и полученного не совпали. Ожидалось: %d. Получено: %d", tc.ExpectedOffSet, capturedOffset)
			}
		})

	}
}

func TestWorkout_GetWorkoutByID_NotMyWorkout(t *testing.T) {
	workoutID := uuid.New()
	workStore := &mockWorkoutStorage{
		getWorkoutByIDFunc: func(ctx context.Context, workoutID uuid.UUID) (domain.Workout, error) {
			return domain.Workout{
				ID:        workoutID,
				AthleteID: uuid.New(),
			}, nil
		},
	}

	workoutService := NewWorkoutService(workStore, &mockExerciseWorkStorage{}, &mockTxManager{})

	_, err := workoutService.GetWorkout(context.Background(), GetWorkoutInput{
		AthleteID: uuid.New(),
		WorkoutID: workoutID,
	})
	if err == nil || !errors.Is(err, domain.ErrNotEnoughPermission) {
		t.Errorf("ожидалась ошибка: %v, получена: %v", domain.ErrNotEnoughPermission, err)
	}
}

func TestWorkout_GetWorkoutByID_Success(t *testing.T) {
	workoutID := uuid.New()
	athleteID := uuid.New()
	dataFromStorage := domain.Workout{
		ID:              workoutID,
		AthleteID:       athleteID,
		TotalTime:       intPtr(50),
		GradeOfTraining: intPtr(9),
		Status:          domain.WorkoutStatusFinished,
		Sets: []domain.Set{
			{ID: uuid.New(),
				ExerciseName: "bench_press",
				ExerciseID:   1,
				WorkoutID:    workoutID,
				Weight:       100,
				SetOrder:     intPtr(1),
				Reps:         10},
			{
				ID:           uuid.New(),
				ExerciseName: "bench_press",
				ExerciseID:   1,
				WorkoutID:    workoutID,
				Weight:       100,
				SetOrder:     intPtr(2),
				Reps:         12},
		},
	}
	workStore := &mockWorkoutStorage{
		getWorkoutByIDFunc: func(ctx context.Context, workoutID uuid.UUID) (domain.Workout, error) {
			return dataFromStorage, nil
		},
	}

	service := NewWorkoutService(workStore, &mockExerciseWorkStorage{}, &mockTxManager{})
	input := GetWorkoutInput{
		AthleteID: athleteID,
		WorkoutID: workoutID,
	}
	workout, err := service.GetWorkout(context.Background(), input)
	if err != nil {
		t.Fatalf("не ожидалось ошибка, получена: %v", err)
	}

	if workout.ID != workoutID {
		t.Errorf("ожидался workoutID: %v, получен: %v", workoutID, workout.ID)
	}

	if len(workout.Sets) != len(dataFromStorage.Sets) {
		t.Fatalf("ожидалось %v подходов, получено %v", len(dataFromStorage.Sets), len(workout.Sets))
	}

	for i, expectedSet := range workout.Sets {
		actualSet := dataFromStorage.Sets[i]
		if actualSet.ExerciseName != expectedSet.ExerciseName {
			t.Errorf("set [%d]: ожидалось exercise_name: %v, получено: %v", i, expectedSet.ExerciseName, actualSet.ExerciseName)
		}
	}

}
