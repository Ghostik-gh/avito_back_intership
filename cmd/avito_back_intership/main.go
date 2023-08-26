package main

import (
	"avito_back_intership/internal/config"
	"avito_back_intership/internal/http-server/handlers/segment/create_segment"
	"avito_back_intership/internal/http-server/handlers/segment/delete_segment"
	"avito_back_intership/internal/http-server/handlers/segment/segment_list"
	"avito_back_intership/internal/http-server/handlers/segment/segment_users"
	"avito_back_intership/internal/http-server/handlers/user/create_user"
	"avito_back_intership/internal/http-server/handlers/user/delete_user"
	"avito_back_intership/internal/http-server/handlers/user/user_list"
	"avito_back_intership/internal/http-server/handlers/user/user_segments"
	"avito_back_intership/internal/http-server/handlers/user_log"

	"avito_back_intership/internal/lib/logger/sl"
	"avito_back_intership/internal/storage/postgres"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "avito_back_intership/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title			Avito Intership (Backend)
// @version			1.0
// @description		Segmentation for A_B tests
// @contact.name	GhostikGH
// @contact.url	 	https://t.me/GhostikGH
// @contact.email	feodor200@mail.ru
// @host			localhost:8002
func main() {
	cfg := config.MustLoad()

	fmt.Println("CONFIG: ", cfg)
	log := setupLogger(cfg.Env)

	log.Info("starting slog")
	log.Debug("debug mod")

	storage, err := postgres.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("storage started success, tables created")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Создает сегмент
	router.Post("/segment/{segment}/{percentage}", create_segment.New(log, storage))
	// Удаляет сегмент
	router.Delete("/segment/{segment}", delete_segment.New(log, storage))
	// Получает список всех сегментов
	router.Get("/segment", segment_list.New(log, storage))
	// Получает список пользователей в данном сегменте
	router.Get("/segment/{segment}", segment_users.New(log, storage))

	// Создает юзера с 0 или более сегментами тут же происходит удаление
	router.Post("/user", create_user.New(log, storage))
	// Удаляет юзера
	router.Delete("/user/{user_id}", delete_user.New(log, storage))
	// Получает всех пользователей
	router.Get("/user", user_list.New(log, storage))
	// Получает сегменты юзера
	router.Get("/user/{user_id}", user_segments.New(log, storage))

	// Выводит csv файл с логом для юзера
	router.Get("/log/{user_id}", user_log.New(log, storage))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// Graceful shutdown

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()
	log.Info("server started", slog.String("address", cfg.Addres))

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage.Close()
	log.Info("storage closed success")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
