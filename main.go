package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HugoBritez/utic.dev-server/internal/application/messages"
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

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC] %v", err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

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
	if err := database.RunMigrations(sqlDB, "db/schema/002_add_messages_features.sql"); err != nil {
		panic("failed to run messages migration: " + err.Error())
	}

	// --- Dependencies ---
	queries := db.New(sqlDB)
	projectRepo := repository.NewProjectRepository(queries, sqlDB)
	aiService := aiservice.NewGroqClient("", "")
	createUseCase := projects.NewCreateProjectUseCase(projectRepo, aiService)
	projectHandler := httphandler.NewProjectHandler(createUseCase, projectRepo)

	messageRepo := repository.NewMessageRepository(queries, sqlDB)
	createMessageUseCase := messages.NewCreateMessageUseCase(messageRepo)
	messageHandler := httphandler.NewMessageHandler(messageRepo, createMessageUseCase)

	sessionStore := middleware.NewSessionStore()
	adminHandler := httphandler.NewAdminHandler(projectRepo, messageRepo, sessionStore)

	// --- Router ---
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.RequestID)
	r.Use(recoveryMiddleware)

	// Public routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Admin UI (session-based auth)
	r.Get("/admin/login", adminHandler.Login)
	r.Post("/admin/login", adminHandler.Login)
	r.Get("/admin/logout", adminHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AdminSession(sessionStore))

		r.Get("/admin", adminHandler.Dashboard)
		r.Get("/admin/stats", adminHandler.Stats)
		r.Get("/admin/projects", adminHandler.ProjectsTable)
		r.Get("/admin/messages", adminHandler.MessagesTable)
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

		r.Post("/api/messages", messageHandler.CreateMessage)
		r.Get("/api/messages", messageHandler.GetMessages)
		r.Get("/api/messages/phone-numbers", messageHandler.GetNewPhoneNumbers)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server running on :%s\n", port)
	log.Printf("AI_API_KEY set: %v", os.Getenv("AI_API_KEY") != "")
	log.Printf("API_KEY set: %v", os.Getenv("API_KEY") != "")
	log.Printf("DB_PATH: %s", dbPath)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
 