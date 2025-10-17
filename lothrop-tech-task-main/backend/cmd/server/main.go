package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/api"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()
	logger.Info("Starting server", zap.String("port", cfg.Port))

	// Initialize database connection
	db, err := database.NewPostgresConnection(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repository, service, and handlers
	companyRepo := repository.NewPostgresCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepo)
	companyHandlers := handlers.NewCompanyHandlers(companyService, logger)

	// Create router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	// CORS middleware for development
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/", handleApiStatus(logger))

		// Company routes
		r.Get("/companies", companyHandlers.GetCompanies)
		r.Post("/companies", companyHandlers.CreateCompany)
		r.Get("/companies/{id}", companyHandlers.GetCompanyByID)
		r.Delete("/companies/{id}", companyHandlers.DeleteCompany)
	})

	// Start server
	logger.Info("Server starting", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func handleApiStatus(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("API status endpoint called")

		response := api.ApiResponse{
			Error: false,
			Msg:   "hello world",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			logger.Error("Failed to encode response", zap.Error(err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
