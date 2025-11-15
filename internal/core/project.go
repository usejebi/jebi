package core

import (
	"fmt"
	"path/filepath"
	"time"

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

func (p *projectService) SaveProjectConfig(id, name, description, defaultEnvironment string) (string, error) {
	project := Project{
		ID:                 id,
		Name:               name,
		Description:        description,
		DefaultEnvironment: defaultEnvironment,
		CreatedAt:          time.Now().UTC(),
		UpdatedAt:          time.Now().UTC(),
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
