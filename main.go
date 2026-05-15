package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HugoBritez/utic.dev-server/internal/application/projects"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/database"
	httphandler "github.com/HugoBritez/utic.dev-server/internal/infrastructure/http"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/middleware"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/repository"
	aiservice "github.com/HugoBritez/utic.dev-server/internal/infrastructure/services"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/db"
)

func main() {
	// --- Load .env ---
	godotenv.Load()

	// --- Database ---
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/app.db"
	}

	sqlDB, err := database.NewSQLite(dbPath)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	defer sqlDB.Close()

	if err := database.RunMigrations(sqlDB, "db/schema/001_init.sql"); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	// --- Dependencies ---
	queries := db.New(sqlDB)
	projectRepo := repository.NewProjectRepository(queries, sqlDB)
	aiService := aiservice.NewGroqClient("", "")
	createUseCase := projects.NewCreateProjectUseCase(projectRepo, aiService)
	projectHandler := httphandler.NewProjectHandler(createUseCase, projectRepo)

	// --- Router ---
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.APIKey())

		r.Post("/api/projects", projectHandler.CreateProject)
		r.Get("/api/projects", projectHandler.ListProjects)
		r.Get("/api/projects/{id}", projectHandler.GetProject)
		r.Post("/api/projects/{id}/star", projectHandler.StarProject)
		r.Get("/api/protected", func(w http.ResponseWriter, r *http.Request) {
			key, _ := middleware.GetAPIKeyFromContext(r.Context())
			w.Write([]byte("Access granted! API key: " + key))
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
 