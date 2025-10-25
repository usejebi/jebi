package core

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jawahars16/jebi/internal/io"
)

type projectService struct {
	workingDir string
}

func NewProjectService(workingDir string) *projectService {
	return &projectService{
		workingDir: workingDir,
	}
}

func (p *projectService) SaveProjectConfig(name, description string) (string, error) {
	project := Project{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	path := filepath.Join(p.workingDir, fmt.Sprintf(".%s", AppName), ProjectConfigFile)
	err := io.WriteJSONToFile(path, project)
	if err != nil {
		return "", fmt.Errorf("failed to write project config: %w", err)
	}
	return project.ID, nil
}

func (p *projectService) LoadProjectConfig() (*Project, error) {
	path := filepath.Join(p.workingDir, fmt.Sprintf(".%s", AppName), ProjectConfigFile)
	project, err := io.ReadJSONFile[Project](path)
	if err != nil {
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}
	return &project, nil
}
