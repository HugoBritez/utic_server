package services

import "context"

type ProjectInfo struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	StudentEmail string   `json:"student_email"`
	StudentName  string   `json:"student_name"`
	TechStack    []string `json:"tech_stack"`
	Categories   []string `json:"categories"`
	Stars        int      `json:"stars"`
}

type AIService interface {
	ExtractProjectInfo(ctx context.Context, repoURL string) (*ProjectInfo, error)
}
