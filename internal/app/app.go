package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dizao9/Fitness-Journal/internal/config"
	"github.com/Dizao9/Fitness-Journal/internal/infrastructure/auth"
	"github.com/Dizao9/Fitness-Journal/internal/service"
	"github.com/Dizao9/Fitness-Journal/internal/storage"
	"github.com/Dizao9/Fitness-Journal/internal/transport"
)

const cost = 12 //следует вынести в .env

func Run(ctx context.Context) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("config load: %w", err)
	}
	db, err := storage.ConnectToDB(cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed connect to db: %w", err)
	}
	defer db.Close()

	jwtManager := auth.NewJWTManager(cfg.AccessJWTSecret, cfg.RefreshJWTSecret)
	ph := auth.NewBcryptPasswordHasher(cost)
	txManager := storage.NewTransactionManager(db)

	athlStorage := storage.NewAthleteStorage(db)
	exercStorage := storage.NewExerciseStorage(db)
	workoutStorage := storage.NewWorkoutStorage(db)

	authService := service.NewAuthService(athlStorage, jwtManager, ph)
	athlService := service.NewAthleteService(athlStorage)
	exercService := service.NewExerciseService(exercStorage)
	workoutService := service.NewWorkoutService(workoutStorage, exercStorage, txManager)

	router := transport.NewRouter(transport.NewAuthHandler(authService), transport.NewAthleteHandler(athlService), transport.NewExerciseHandler(exercService), transport.NewWorkoutHandler(workoutService))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Server is running on: %s \n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverError <- fmt.Errorf("listen: %s\n", err)
		}
	}()
	select {
	case err := <-serverError:
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
