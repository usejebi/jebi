package core

import (
	"fmt"
	"path/filepath"

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

func (p *projectService) SaveProjectConfig(name, description string) error {
	project := Project{
		Name:        name,
		Description: description,
	}
	path := filepath.Join(p.workingDir, fmt.Sprintf(".%s", AppName), ProjectConfigFile)
	err := io.WriteJSONToFile(path, project)
	if err != nil {
		return fmt.Errorf("failed to write project config: %w", err)
	}
	return nil
}
