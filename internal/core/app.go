package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type appService struct {
	workingDir string
}

func NewAppService(workingDir string) *appService {
	return &appService{
		workingDir: workingDir,
	}
}

func (s *appService) CreateAppDir() error {
	dirName := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName))
	info, err := os.Stat(dirName)
	if err == nil && info.IsDir() {
		// Directory exists; check if it's empty
		entries, err := os.ReadDir(dirName)
		if err != nil {
			return fmt.Errorf("failed to read directory contents: %w", err)
		}

		if len(entries) > 0 {
			return fmt.Errorf(
				"project already initialized in %q. Use a different directory or remove %q to reinitialize",
				dirName, dirName,
			)
		}

		// Directory exists but is empty — fine to reuse
		return nil
	}

	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		if err := os.Mkdir(dirName, 0755); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", dirName, err)
		}
		return nil
	}

	// Some other error
	return fmt.Errorf("failed to check directory %q: %w", dirName, err)
}

func (s *appService) Exists() (bool, error) {
	dirName := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName))
	info, err := os.Stat(dirName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check directory %q: %w", dirName, err)
	}
	return info.IsDir(), nil
}
