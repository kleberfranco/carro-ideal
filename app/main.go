package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"carro-ideal/app/clients"
	"carro-ideal/app/db"
	"carro-ideal/app/internal/admin"
	"carro-ideal/app/internal/api"
	"carro-ideal/app/internal/health"
	platformmw "carro-ideal/app/internal/platform"
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
	sessionRepo := repository.NewSessionRepository(db.GetDB())
	questionRepo := repository.NewQuestionRepository(db.GetDB())
	vehicleRepo := repository.NewVehicleRepository(db.GetDB())
	recommendationRepo := repository.NewRecommendationRepository(db.GetDB())
	adminRepo := repository.NewAdminRepository(db.GetDB())
	catalogCache := service.NewCatalogCache(time.Duration(cfg.CacheTTL) * time.Second)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(sessionRepo)
	questionnaireService := service.NewQuestionnaireService(questionRepo, catalogCache)
	vehicleService := service.NewVehicleService(vehicleRepo, catalogCache)
	adminService := service.NewAdminService(adminRepo, catalogCache)

	var aiService *service.AIService
	if cfg.OpenAIAPIKey != "" {
		openAIClient := clients.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.OpenAITimeoutSecs)
		aiService = service.NewAIService(openAIClient)
	}
	recommendationService := service.NewRecommendationService(questionnaireService, vehicleService, recommendationRepo, aiService)
	secureCookie := strings.EqualFold(cfg.Environment, "production")
	logger := newLogger(cfg.LogLevel)

	webHandler := web.NewHandler(userService, authService, secureCookie)
	apiHandler := api.NewHandler(
		userService,
		authService,
		questionnaireService,
		recommendationService,
		vehicleService,
		secureCookie,
	)
	adminHandler := admin.NewHandler(userService, authService, adminService)
	healthHandler := health.NewHandler(db.GetDB())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(platformmw.RequestIDHeader)
	r.Use(platformmw.Recovery(logger))
	r.Use(platformmw.RequestLogger(logger))
	r.Use(platformmw.SecurityHeaders)
	r.Use(platformmw.CORS(cfg.AllowedOrigins))
	r.Use(platformmw.RateLimit(cfg.RateLimit, time.Duration(cfg.RateWindow)*time.Second))
	r.Use(platformmw.CSRF(secureCookie))

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
		r.Handle("/", api.RequireAuth(authService, http.HandlerFunc(apiHandler.Placeholder)))
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(api.JSONMiddleware)

		r.With(func(next http.Handler) http.Handler {
			return api.RequireAuth(authService, next)
		}).Get("/questions", apiHandler.Questions)

		r.Route("/questions", func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return api.RequireAuth(authService, next)
			})
			r.Get("/", apiHandler.Questions)
		})

		r.With(func(next http.Handler) http.Handler {
			return api.RequireAuth(authService, next)
		}).Get("/recommendations", apiHandler.RecommendationHistory)

		r.Route("/recommendations", func(r chi.Router) {
			r.Use(func(next http.Handler) http.Handler {
				return api.RequireAuth(authService, next)
			})
			r.Post("/generate", apiHandler.GenerateRecommendations)
			r.Get("/", apiHandler.RecommendationHistory)
			r.Get("/{id}", apiHandler.RecommendationDetails)
		})

		r.Get("/vehicles", apiHandler.Vehicles)

		r.Route("/vehicles", func(r chi.Router) {
			r.Get("/", apiHandler.Vehicles)
			r.Get("/{id}", apiHandler.VehicleDetails)
		})
	})

	r.Route("/api/admin", func(r chi.Router) {
		r.Use(api.JSONMiddleware)
		r.Use(func(next http.Handler) http.Handler {
			return api.RequireAuth(authService, next)
		})
		r.Use(func(next http.Handler) http.Handler {
			return api.RequireAdmin(userService, next)
		})
		admin.RegisterRoutes(r, adminHandler)
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
	r.Get("/admin", adminHandler.Page)
	r.Get("/logout", webHandler.LogoutHandler)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	r.Get("/", webHandler.HomeHandler)
	r.NotFound(platformmw.NotFound)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		var serveErr error
		if cfg.TLSCertFile != "" {
			logger.Info("starting HTTPS server", "port", cfg.Port)
			serveErr = srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			logger.Info("starting HTTP server", "port", cfg.Port)
			serveErr = srv.ListenAndServe()
		}
		if serveErr != nil && serveErr != http.ErrServerClosed {
			log.Fatalf("server error: %v", serveErr)
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

func newLogger(level string) *slog.Logger {
	var parsed slog.Level
	switch strings.ToLower(level) {
	case "debug":
		parsed = slog.LevelDebug
	case "warn":
		parsed = slog.LevelWarn
	case "error":
		parsed = slog.LevelError
	default:
		parsed = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parsed}))
}
