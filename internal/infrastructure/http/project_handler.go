package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/HugoBritez/utic.dev-server/internal/application/projects"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	createUseCase *projects.CreateProjectUseCase
	repo          repositories.ProjectRepository
}

func NewProjectHandler(createUseCase *projects.CreateProjectUseCase, repo repositories.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{createUseCase: createUseCase, repo: repo}
}

type createProjectRequest struct {
	RepoURL string `json:"repo_url"`
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req createProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.RepoURL == "" {
		http.Error(w, `{"error":"repo_url is required"}`, http.StatusBadRequest)
		return
	}

	project, err := h.createUseCase.Execute(r.Context(), req.RepoURL)
	if err != nil {
		log.Printf("[ERROR] CreateProject: %v", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to list projects"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	project, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"project not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) StarProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	project, err := h.repo.Star(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"failed to star project"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
