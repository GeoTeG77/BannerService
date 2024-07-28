package main

import (
	"BannerService/cmd/internal/config"
	banner "BannerService/cmd/internal/http-server/handlers/url/banner"
	feature "BannerService/cmd/internal/http-server/handlers/url/feature"
	"BannerService/cmd/internal/http-server/handlers/url/password"
	auth "BannerService/cmd/internal/http-server/middleware/auth"

	//"BannerService/cmd/internal/http-server/handlers/url/password"
	tag "BannerService/cmd/internal/http-server/handlers/url/tag"

	//"BannerService/cmd/internal/lib/logger/sl"
	"BannerService/cmd/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	os.Setenv("JWT KEY", "Secret228")

	log := setupLogger(cfg.Env)
	log.Info("debug messages are enabled", slog.String("env", cfg.Env))

	auth.LoadTemplates()
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	defer storage.Close()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Get("/", auth.LoginPageHandler)
	router.Post("/login", auth.LoginHandler)
	router.Get("/register", auth.RegisterPageHandler)
	router.Post("/register", auth.RegisterHandler)
	router.Post("/auth", password.CreateUser(log, storage))
	router.Get("/auth", password.CheckPassword(log, storage))
	router.Get("/logout", auth.Logout)

	router.Route("/", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Use(middleware.SetHeader("Content-Type", "application/json"))
		r.Use(middleware.Recoverer)
		r.Use(middleware.URLFormat)

		r.Post("/banner", banner.SetBanner(log, storage))
		r.Post("/tag", tag.CreateTag(log, storage))
		r.Post("/feature", feature.CreateFeature(log, storage))

		r.Patch("/banner/{id}", banner.UpdateBanner(log, storage))
		r.Patch("/feature/{id}/{name}", feature.UpdateFeatureName(log, storage))
		r.Patch("/tag/{id}/{name}", tag.UpdateTagName(log, storage))

		r.Get("/user_banner", banner.GetBanner(log, storage))
		r.Get("/banner", banner.GetBanners(log, storage))

		r.Delete("/banner/{id}", banner.DeleteBanner(log, storage))
		r.Delete("/tag/{id}", tag.DeleteTag(log, storage))
		r.Delete("/feature/{id}", feature.DeleteFeature(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.TimeOut,
		WriteTimeout: cfg.HTTPServer.TimeOut,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
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
