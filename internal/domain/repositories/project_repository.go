package repositories

import (
	"context"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) (*entities.Project, error)
	GetByID(ctx context.Context, id string) (*entities.Project, error)
	GetByRepoURL(ctx context.Context, repoURL string) (*entities.Project, error)
	List(ctx context.Context) ([]entities.Project, error)
	GetByStudentEmail(ctx context.Context, email string) ([]entities.Project, error)
	Update(ctx context.Context, project *entities.Project) (*entities.Project, error)
	Delete(ctx context.Context, id string) error
	Star(ctx context.Context, id string) (*entities.Project, error)

	GetByCategory(ctx context.Context, category string) ([]entities.Project, error)
	GetByTechStack(ctx context.Context, techStack string) ([]entities.Project, error)
	GetByAnyCategory(ctx context.Context, categories []string) ([]entities.Project, error)
	GetByAnyTechStack(ctx context.Context, techStacks []string) ([]entities.Project, error)
}
