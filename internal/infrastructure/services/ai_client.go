package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/HugoBritez/utic.dev-server/internal/domain/services"
)

type GroqClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewGroqClient(apiKey string, model string) *GroqClient {
	if apiKey == "" {
		apiKey = os.Getenv("AI_API_KEY")
	}
	if model == "" {
		model = os.Getenv("AI_MODEL")
		if model == "" {
			model = "llama-3.3-70b-versatile"
		}
	}
	return &GroqClient{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{},
	}
}

func (c *GroqClient) ExtractProjectInfo(ctx context.Context, repoURL string) (*services.ProjectInfo, error) {
	repoMeta, err := c.fetchRepoMetadata(repoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repo metadata: %w", err)
	}

	ownerProfile, _ := c.fetchOwnerProfile(repoMeta.Owner)

	readmeContent, err := c.fetchRepoREADME(repoURL)
	if err != nil {
		readmeContent = "(README no disponible)"
	}

	prompt := fmt.Sprintf(`Analiza este repositorio de GitHub y extrae la información en formato JSON.

INSTRUCCIONES POR CAMPO:
- name: Usá el nombre del repo convertido a título. Ej: "golang-walking-skeleton" → "Golang Walking Skeleton"
- description: Descripción de 1-2 oraciones basada en la descripción del repo o README.
- student_email: Usá el email público del perfil del owner si existe. Si no, dejá "".
- student_name: Usá el nombre del perfil del owner. Si no tiene, usá el username de GitHub ("%s").
- tech_stack: ARRAY OBLIGATORIO. Incluí SIEMPRE el lenguaje principal ("%s") como primer elemento. Luego agregá frameworks, librerías y herramientas que veas en el README o README. Ejemplo: si el lenguaje es "Go" y el README menciona "chi", devolvé ["Go", "chi"].
- categories: 1-3 categorías en minúsculas. Opciones válidas: web, api, mobile, cli, desktop, machine-learning, devops, game, library, tool, educational, portfolio, backend, frontend, fullstack.
- stars: El número exacto de estrellas (%d).

JSON DE EJEMPLO (formato exacto que espero):
{
  "name": "Mi Proyecto",
  "description": "Descripción del proyecto",
  "student_email": "user@email.com",
  "student_name": "Nombre del Autor",
  "tech_stack": ["Go", "chi", "PostgreSQL"],
  "categories": ["api", "backend"],
  "stars": 42
}

---
DATOS DEL REPOSITORIO:
- Nombre: %s
- Descripción: %s
- Lenguaje: %s
- Estrellas: %d
- Owner username: %s
- Owner name: %s
- Owner email público: %s

---
README:
%s

---
Respondé SOLO con el JSON válido.`,
		repoMeta.Owner,
		repoMeta.Language,
		repoMeta.Stars,
		repoMeta.Name,
		repoMeta.Description,
		repoMeta.Language,
		repoMeta.Stars,
		repoMeta.Owner,
		ownerProfile.Name,
		ownerProfile.Email,
		readmeContent,
	)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `Extract GitHub repo metadata as JSON. ALWAYS include the primary language in tech_stack. Use owner profile data for student_name and student_email when available.`,
			},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.1,
		"response_format": map[string]string{"type": "json_object"},
	})

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Groq API error %d: %s", resp.StatusCode, string(body))
	}

	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &aiResp); err != nil {
		return nil, err
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Groq")
	}

	content := strings.TrimSpace(aiResp.Choices[0].Message.Content)
	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var info services.ProjectInfo
	if err := json.Unmarshal([]byte(content), &info); err != nil {
		return nil, fmt.Errorf("failed to parse Groq response: %w, raw: %s", err, content)
	}

	// Fallback: si la IA no devolvió tech_stack pero hay lenguaje, lo inyectamos
	if len(info.TechStack) == 0 && repoMeta.Language != "" {
		info.TechStack = []string{repoMeta.Language}
	}

	// Fallback: si no hay student_name, usamos el owner
	if info.StudentName == "" {
		info.StudentName = repoMeta.Owner
	}

	// Fallback: si no hay email, usamos el del profile
	if info.StudentEmail == "" && ownerProfile.Email != "" {
		info.StudentEmail = ownerProfile.Email
	}

	return &info, nil
}

type repoMetadata struct {
	Name        string
	Description string
	Language    string
	Stars       int
	Forks       int
	Owner       string
	CreatedAt   string
	UpdatedAt   string
}

type ownerProfile struct {
	Name  string
	Email string
}

func (c *GroqClient) fetchRepoMetadata(repoURL string) (*repoMetadata, error) {
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	var data struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Language    string `json:"language"`
		Stars       int    `json:"stargazers_count"`
		Forks       int    `json:"forks_count"`
		Owner       struct {
			Login string `json:"login"`
		} `json:"owner"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &repoMetadata{
		Name:        data.Name,
		Description: data.Description,
		Language:    data.Language,
		Stars:       data.Stars,
		Forks:       data.Forks,
		Owner:       data.Owner.Login,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}, nil
}

func (c *GroqClient) fetchOwnerProfile(username string) (*ownerProfile, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &ownerProfile{}, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &ownerProfile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &ownerProfile{}, nil
	}

	var data struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	json.NewDecoder(resp.Body).Decode(&data)

	return &ownerProfile{
		Name:  data.Name,
		Email: data.Email,
	}, nil
}

func (c *GroqClient) fetchRepoREADME(repoURL string) (string, error) {
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func parseGitHubURL(repoURL string) (owner string, repo string, err error) {
	parts := strings.Split(strings.TrimSuffix(strings.TrimPrefix(repoURL, "https://github.com/"), "/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL: %s", repoURL)
	}
	return parts[0], parts[1], nil
}
