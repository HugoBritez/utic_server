package projects

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/HugoBritez/utic.dev-server/internal/domain/entities"
	"github.com/HugoBritez/utic.dev-server/internal/domain/repositories"
	"github.com/HugoBritez/utic.dev-server/internal/domain/services"
)

type CreateProjectUseCase struct {
	repo repositories.ProjectRepository
	ai   services.AIService
}

func NewCreateProjectUseCase(repo repositories.ProjectRepository, ai services.AIService) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo, ai: ai}
}

func marshalStringList(list []string) string {
	if list == nil {
		return "[]"
	}
	b, _ := json.Marshal(list)
	return string(b)
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, repoURL string) (*entities.Project, error) {
	existing, err := uc.repo.GetByRepoURL(ctx, repoURL)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	info, err := uc.ai.ExtractProjectInfo(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		existing.Name = info.Name
		existing.Description = info.Description
		existing.StudentEmail = info.StudentEmail
		existing.StudentName = info.StudentName
		existing.TechStack = marshalStringList(info.TechStack)
		existing.Categories = marshalStringList(info.Categories)
		existing.Stars = int64(info.Stars)

		return uc.repo.Update(ctx, existing)
	}

	project := entities.NewProject(
		info.Name,
		info.Description,
		info.StudentEmail,
		info.StudentName,
		marshalStringList(info.TechStack),
		marshalStringList(info.Categories),
		info.Stars,
		repoURL,
	)

	return uc.repo.Create(ctx, project)
}
