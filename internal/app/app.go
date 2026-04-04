package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/Dizao9/Fitness-Journal/internal/storage"
	"github.com/Dizao9/Fitness-Journal/internal/transport"
)

func NewRouter(handlers *transport.Handlers) *http.ServeMux {
	mux := http.NewServeMux()
	//exercise block
	protected := handlers.Auth.AuthMiddlware
	mux.Handle("POST /exercise", protected(http.HandlerFunc(handlers.Exercise.PostExercise)))
	mux.Handle("GET /ListExercises", protected(http.HandlerFunc(handlers.Exercise.GetPageOfExercise)))
	mux.Handle("GET /exercise/{id}", protected(http.HandlerFunc(handlers.Exercise.GetExerciseByID)))
	mux.Handle("DELETE /exercise/{id}", protected(http.HandlerFunc(handlers.Exercise.DeleteExercise)))
	mux.Handle("PUT /exercise/{id}", protected(http.HandlerFunc(handlers.Exercise.UpdateExercise)))
	//auth block
	mux.HandleFunc("POST /auth/register", handlers.Auth.RegisterUser)
	mux.HandleFunc("POST /auth/login", handlers.Auth.Login)
	mux.HandleFunc("POST /auth/refresh", handlers.Auth.Refresh)
	mux.HandleFunc("POST /auth/logout", handlers.Auth.LogOut)
	//user block
	mux.Handle("GET /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.GetProfile)))
	mux.Handle("PUT /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.UpdateUserProfile)))
	mux.Handle("DELETE /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.DeleteUser)))
	//workout block
	mux.Handle("POST /workout", protected(http.HandlerFunc(handlers.Workout.CreateTraining)))
	return mux
}

func Run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config load: %w", err)
	}
	db, err := storage.ConnectToDB(cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed connect to db: %w", err)
	}

	storage := storage.NewStorage(db)
	services := service.NewServices(storage, cfg)
	handlers := transport.NewHandlers(services)

	mux := NewRouter(handlers)
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- fmt.Errorf("listen: %s\n", err)
		}
	}()
	select {
	case <-serverError:
		return err
	case <-ctx.Done():
		log.Println("Shutting down..")
	}

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		return fmt.Errorf("Shutdown server :%w", err)
	}

	if err := db.Close(); err != nil {
		return fmt.Errorf("closing DB: %w", err)
	}

	log.Println("Server exited properly")
	return nil
}
