package app

import (
	"context"
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
	//auth block
	mux.HandleFunc("POST /auth/register", handlers.Auth.RegisterUser)
	mux.HandleFunc("POST /auth/login", handlers.Auth.Login)
	//user block
	mux.Handle("GET /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.GetProfile)))
	mux.Handle("PUT /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.UpdateUserProfile)))
	mux.Handle("DELETE /athlete/profile", protected(http.HandlerFunc(handlers.Athlete.DeleteUser)))

	return mux
}

func Run(ctx context.Context) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config data: %v", err)
	}
	db, err := storage.ConnectToDB(cfg.DSN)
	if err != nil {
		log.Fatalf("failed connect to db: %v", err)
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

	go func() {
		log.Println("Server is running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down..")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("Error closing DB: %v", err)
	}

	log.Println("Server exited properly")
}
