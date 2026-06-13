package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"carro-ideal/app/db"
	"carro-ideal/app/internal/admin"
	"carro-ideal/app/internal/api"
	"carro-ideal/app/internal/health"
	"carro-ideal/app/internal/web"
	"carro-ideal/app/repository"
	"carro-ideal/app/service"
	"carro-ideal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.RunMigrations(cfg.DatabaseURL()); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db.GetDB())
	userService := service.NewUserService(userRepo)

	webHandler := web.NewHandler(userService)
	apiHandler := api.NewHandler(userService)
	adminHandler := admin.NewHandler(userService)
	healthHandler := health.NewHandler(db.GetDB())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler.Health)

	r.Route("/api/auth", func(r chi.Router) {
		r.Use(api.JSONMiddleware)
		r.Post("/register", apiHandler.Register)
		r.Post("/login", apiHandler.Login)
		r.Post("/logout", apiHandler.Logout)
		r.Get("/me", apiHandler.Me)
	})

	r.Route("/api/user", func(r chi.Router) {
		r.Use(api.JSONMiddleware)
		r.Get("/", api.RequireAuth(http.HandlerFunc(apiHandler.Placeholder)))
	})

	r.Route("/api/admin", func(r chi.Router) {
		r.Get("/", adminHandler.AdminHandler)
	})

	r.Route("/web", func(r chi.Router) {
		r.Get("/", webHandler.HomeHandler)
		r.Get("/login", webHandler.LoginHandler)
		r.Get("/register", webHandler.RegisterHandler)
		r.Get("/recommend", webHandler.RecommendHandler)
	})

	r.Get("/login", webHandler.LoginHandler)
	r.Get("/register", webHandler.RegisterHandler)
	r.Get("/recommend", webHandler.RecommendHandler)
	r.Get("/logout", webHandler.LogoutHandler)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.Get("/", webHandler.HomeHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("starting server on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("server stopped cleanly")
}
