package entities

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID           string
	Name         string
	Description  string
	StudentEmail string
	StudentName  string
	TechStack    string
	Categories   string
	Stars        int64
	RepoUrl      string
	CreatedAt    time.Time
}

func NewProject(name string, description string, studentEmail string, studentName string, techStack string, categories string, stars int, repoUrl string) *Project {
	now := time.Now()

	return &Project{
		ID:           uuid.New().String(),
		Name:         name,
		Description:  description,
		StudentEmail: studentEmail,
		StudentName:  studentName,
		TechStack:    techStack,
		Categories:   categories,
		Stars:        int64(stars),
		RepoUrl:      repoUrl,
		CreatedAt:    now,
	}
}
