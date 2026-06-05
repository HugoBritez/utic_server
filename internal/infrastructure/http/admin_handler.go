package http

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/middleware"
)

//go:embed templates/admin/*.html
var adminTemplates embed.FS

type AdminHandler struct {
	projectRepo repositories.ProjectRepository
	messageRepo repositories.MessageRepository
	sessions    *middleware.SessionStore
	tmpl        *template.Template
}

func NewAdminHandler(projectRepo repositories.ProjectRepository, messageRepo repositories.MessageRepository, sessions *middleware.SessionStore) *AdminHandler {
	return &AdminHandler{
		projectRepo: projectRepo,
		messageRepo: messageRepo,
		sessions:    sessions,
		tmpl:        template.Must(template.ParseFS(adminTemplates, "templates/admin/*.html")),
	}
}

type statsData struct {
	Projects     int
	Messages     int
	PhoneNumbers int
}

func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := r.FormValue("username")
		pass := r.FormValue("password")

		expectedUser := os.Getenv("ADMIN_USER")
		expectedPass := os.Getenv("ADMIN_PASSWORD")

		if user == expectedUser && pass == expectedPass {
			token := h.sessions.Create(user)
			http.SetCookie(w, &http.Cookie{
				Name:     "admin_session",
				Value:    token,
				Path:     "/admin",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			w.Header().Set("Content-Type", "text/html")
			h.tmpl.ExecuteTemplate(w, "login.html", "Credenciales inválidas")
		}
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "login.html", "")
}

func (h *AdminHandler) Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("admin_session")
	if err == nil {
		h.sessions.Delete(c.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/admin/login", http.StatusFound)
}

func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "index.html", nil)
}

func (h *AdminHandler) Stats(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectRepo.List(r.Context())
	if err != nil {
		log.Printf("[ERROR] Stats projects: %v", err)
	}

	messages, err := h.messageRepo.GetMessages(r.Context())
	if err != nil {
		log.Printf("[ERROR] Stats messages: %v", err)
	}

	seen := map[string]bool{}
	for _, m := range messages {
		seen[m.SenderPhoneNumber] = true
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "stats.html", statsData{
		Projects:     len(projects),
		Messages:     len(messages),
		PhoneNumbers: len(seen),
	})
}

func (h *AdminHandler) ProjectsTable(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectRepo.List(r.Context())
	if err != nil {
		http.Error(w, "Error loading projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "projects.html", projects)
}

func (h *AdminHandler) MessagesTable(w http.ResponseWriter, r *http.Request) {
	messages, err := h.messageRepo.GetMessages(r.Context())
	if err != nil {
		http.Error(w, "Error loading messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	h.tmpl.ExecuteTemplate(w, "messages.html", messages)
}
