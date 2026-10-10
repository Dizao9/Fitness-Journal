package transport

import "net/http"

func NewRouter(auth *AuthHandler, athlete *AthleteHandler, exercise *ExerciseHandler, workout *WorkoutHandler) *http.ServeMux {
	mux := http.NewServeMux()
	//exercise block
	protected := auth.AuthMiddlware
	mux.Handle("POST /exercise", protected(http.HandlerFunc(exercise.PostExercise)))
	mux.Handle("GET /ListExercises", protected(http.HandlerFunc(exercise.GetPageOfExercise)))
	mux.Handle("GET /exercise/{id}", protected(http.HandlerFunc(exercise.GetExerciseByID)))
	mux.Handle("DELETE /exercise/{id}", protected(http.HandlerFunc(exercise.DeleteExercise)))
	mux.Handle("PUT /exercise/{id}", protected(http.HandlerFunc(exercise.UpdateExercise)))
	//auth block
	mux.HandleFunc("POST /auth/register", auth.RegisterUser)
	mux.HandleFunc("POST /auth/login", auth.Login)
	mux.HandleFunc("POST /auth/refresh", auth.Refresh)
	mux.HandleFunc("POST /auth/logout", auth.LogOut)
	//user block
	mux.Handle("GET /athlete/profile", protected(http.HandlerFunc(athlete.GetProfile)))
	mux.Handle("PUT /athlete/profile", protected(http.HandlerFunc(athlete.UpdateUserProfile)))
	mux.Handle("DELETE /athlete/profile", protected(http.HandlerFunc(athlete.DeleteUser)))
	//workout block
	mux.Handle("POST /workout", protected(http.HandlerFunc(workout.CreateTraining)))
	return mux
}
