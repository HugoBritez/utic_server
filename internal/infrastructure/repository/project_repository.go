package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
	"github.com/HugoBritez/utic.dev-server/internal/infrastructure/db"
)

type ProjectRepository struct {
	queries *db.Queries
	db      *sql.DB
}

func NewProjectRepository(queries *db.Queries, db *sql.DB) repositories.ProjectRepository {
	return &ProjectRepository{queries: queries, db: db}
}

func toEntity(p db.Projects) entities.Project {
	return entities.Project{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		StudentEmail: p.StudentEmail,
		StudentName:  p.StudentName,
		TechStack:    p.TechStack,
		Categories:   p.Categories,
		Stars:        p.Stars,
		RepoUrl:      p.RepoUrl,
		CreatedAt:    p.CreatedAt,
	}
}

func toEntityList(projects []db.Projects) []entities.Project {
	result := make([]entities.Project, len(projects))
	for i, p := range projects {
		result[i] = toEntity(p)
	}
	return result
}

func (r *ProjectRepository) Create(ctx context.Context, project *entities.Project) (*entities.Project, error) {
	p, err := r.queries.CreateProject(ctx, db.CreateProjectParams{
		ID:           project.ID,
		Name:         project.Name,
		Description:  project.Description,
		StudentEmail: project.StudentEmail,
		StudentName:  project.StudentName,
		TechStack:    project.TechStack,
		Categories:   project.Categories,
		Stars:        project.Stars,
		RepoUrl:      project.RepoUrl,
	})
	if err != nil {
		return nil, err
	}
	e := toEntity(p)
	return &e, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*entities.Project, error) {
	p, err := r.queries.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	e := toEntity(p)
	return &e, nil
}

func (r *ProjectRepository) GetByRepoURL(ctx context.Context, repoURL string) (*entities.Project, error) {
	p, err := r.queries.GetProjectByRepoURL(ctx, repoURL)
	if err != nil {
		return nil, err
	}
	e := toEntity(p)
	return &e, nil
}

func (r *ProjectRepository) List(ctx context.Context) ([]entities.Project, error) {
	projects, err := r.queries.GetProjects(ctx)
	if err != nil {
		return nil, err
	}
	return toEntityList(projects), nil
}

func (r *ProjectRepository) GetByStudentEmail(ctx context.Context, email string) ([]entities.Project, error) {
	projects, err := r.queries.GetProjectsByStudent(ctx, email)
	if err != nil {
		return nil, err
	}
	return toEntityList(projects), nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *entities.Project) (*entities.Project, error) {
	p, err := r.queries.UpdateProject(ctx, db.UpdateProjectParams{
		ID:           project.ID,
		Name:         project.Name,
		Description:  project.Description,
		StudentEmail: project.StudentEmail,
		StudentName:  project.StudentName,
		TechStack:    project.TechStack,
		Categories:   project.Categories,
		Stars:        project.Stars,
		RepoUrl:      project.RepoUrl,
	})
	if err != nil {
		return nil, err
	}
	e := toEntity(p)
	return &e, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return r.queries.DeleteProject(ctx, id)
}

func (r *ProjectRepository) Star(ctx context.Context, id string) (*entities.Project, error) {
	p, err := r.queries.StarProject(ctx, id)
	if err != nil {
		return nil, err
	}
	e := toEntity(p)
	return &e, nil
}

func (r *ProjectRepository) GetByCategory(ctx context.Context, category string) ([]entities.Project, error) {
	query := `SELECT * FROM projects WHERE EXISTS (SELECT 1 FROM json_each(categories) WHERE value = ?) ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Projects
	for rows.Next() {
		var p db.Projects
		if err := scanProject(rows, &p); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return toEntityList(projects), rows.Err()
}

func (r *ProjectRepository) GetByTechStack(ctx context.Context, techStack string) ([]entities.Project, error) {
	query := `SELECT * FROM projects WHERE EXISTS (SELECT 1 FROM json_each(tech_stack) WHERE value = ?) ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, techStack)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Projects
	for rows.Next() {
		var p db.Projects
		if err := scanProject(rows, &p); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return toEntityList(projects), rows.Err()
}

func (r *ProjectRepository) GetByAnyCategory(ctx context.Context, categories []string) ([]entities.Project, error) {
	projects, err := r.queryByJSONList(ctx, "categories", categories)
	if err != nil {
		return nil, err
	}
	return toEntityList(projects), nil
}

func (r *ProjectRepository) GetByAnyTechStack(ctx context.Context, techStacks []string) ([]entities.Project, error) {
	projects, err := r.queryByJSONList(ctx, "tech_stack", techStacks)
	if err != nil {
		return nil, err
	}
	return toEntityList(projects), nil
}

func (r *ProjectRepository) queryByJSONList(ctx context.Context, column string, values []string) ([]db.Projects, error) {
	args := make([]interface{}, len(values))
	placeholders := make([]string, len(values))
	for i, v := range values {
		args[i] = v
		placeholders[i] = "?"
	}

	query := `SELECT * FROM projects WHERE EXISTS (SELECT 1 FROM json_each(` + column + `) WHERE value IN (` + strings.Join(placeholders, ",") + `)) ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []db.Projects
	for rows.Next() {
		var p db.Projects
		if err := scanProject(rows, &p); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func scanProject(rows *sql.Rows, p *db.Projects) error {
	return rows.Scan(&p.ID, &p.Name, &p.Description, &p.StudentEmail, &p.StudentName, &p.TechStack, &p.Categories, &p.Stars, &p.RepoUrl, &p.CreatedAt)
}
