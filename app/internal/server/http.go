package server

import (
	"context"
	"fmt"
	"image-app/internal/config"
	"image-app/internal/core"
	"image-app/internal/image"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func StartServer(ctx context.Context) {
	// init config
	cfg := config.InitConfig()

	// setup dependencies
	deps := core.NewDependency(ctx, cfg)

	// setup server
	var srv http.Server

	// register routes
	routes := InitRoutes(ctx, deps)
	srv.Handler = routes

	// setup graceful shutdown
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		deps.Logger.InfoContext(ctx, "server is shutting down")

		ctx := context.Background()
		ctxCncl, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
		cancel()
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctxCncl); err != nil {
			// Error from closing listeners, or context timeout:
			deps.Logger.InfoContext(ctx, "fail to shut down", slog.Any("error", err))
		}
		close(idleConnectionClosed)
	}()

	srv.Addr = fmt.Sprintf(":%d", cfg.Port)
	deps.Logger.Info("server starts...")

	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func InitRoutes(ctx context.Context, deps *core.Dependency) http.Handler {
	// wiring shared packages

	// wiring repository layer
	imgRepo := image.NewImageRepository(deps)

	// wiring logic layer
	imgLogic := image.NewImageLogic(deps, imgRepo)

	// wiring handler layer
	imgHandler := image.NewImageHandler(deps, imgLogic)

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Post("/upload", imgHandler.UploadImage)
	r.Get("/images", imgHandler.GetImages)

	return r
}
