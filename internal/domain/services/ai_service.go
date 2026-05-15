package services

import "context"

type ProjectInfo struct {
	Name        string
	Description string
	StudentEmail string
	StudentName string
	TechStack   []string
	Categories  []string
	Stars       int
}

type AIService interface {
	ExtractProjectInfo(ctx context.Context, repoURL string) (*ProjectInfo, error)
}
