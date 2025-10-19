package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/io"
)

var (
	ErrCurrentEnvNotExist = fmt.Errorf("current environment does not exist")
)

type envService struct {
	workingDir string
}

func NewEnvService(workingDir string) *envService {
	return &envService{
		workingDir: workingDir,
	}
}

// currentEnvPath returns the path to the "current" file
func (e *envService) currentEnvPath() string {
	return filepath.Join(e.workingDir, fmt.Sprintf(".%s", AppName), CurrentFileName)
}

// envDir returns the directory for a given environment
func (e *envService) envDir(env string) string {
	return filepath.Join(e.workingDir, fmt.Sprintf(".%s", AppName), env)
}

// CurrentEnv reads the active environment from ".<AppName>/current"
func (e *envService) CurrentEnv() (string, error) {
	path := e.currentEnvPath()
	currentEnv, err := io.ReadJSONFile[CurrentEnv](path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrCurrentEnvNotExist
		}
		return "", err
	}
	return currentEnv.Env, nil
}

func (e *envService) GetCurrentEnv() (*CurrentEnv, error) {
	path := e.currentEnvPath()
	currentEnv, err := io.ReadJSONFile[CurrentEnv](path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCurrentEnvNotExist
		}
		return nil, err
	}
	return &currentEnv, nil
}

// SetCurrentEnv sets the current active environment in ".<AppName>/current"
func (e *envService) SetCurrentEnv(env string) error {
	path := e.currentEnvPath()
	currentEnv := CurrentEnv{
		Env:     env,
		Changes: []Change{},
	}
	return io.WriteJSONToFile(path, currentEnv)
}

// CreateEnv creates a new environment folder: ".<AppName>/<env>"
func (e *envService) CreateEnv(env string) error {
	envDir := e.envDir(env)
	if err := os.MkdirAll(envDir, 0700); err != nil {
		return fmt.Errorf("failed to create environment '%s': %w", env, err)
	}
	return nil
}

func (e *envService) EnvExists(env string) (bool, error) {
	envDir := e.envDir(env)
	info, err := os.Stat(envDir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func (e *envService) RemoveEnv(env string) error {
	// check if env exists
	exists, err := e.EnvExists(env)
	if err != nil {
		return fmt.Errorf("failed to check if environment exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("environment '%s' does not exist", env)
	}

	envDir := e.envDir(env)
	if err := os.RemoveAll(envDir); err != nil {
		return fmt.Errorf("failed to delete environment '%s': %w", env, err)
	}

	// if the deleted env was the current env, unset current env
	currentEnv, err := io.ReadJSONFile[CurrentEnv](e.currentEnvPath())
	if err == nil && currentEnv.Env == env {
		if err := os.Remove(e.currentEnvPath()); err != nil {
			return fmt.Errorf("failed to unset current environment: %w", err)
		}
	}

	return nil
}

func (e *envService) HasPendingChanges() (bool, error) {
	path := e.currentEnvPath()
	currentEnv, err := io.ReadJSONFile[CurrentEnv](path)
	if err != nil {
		return false, fmt.Errorf("failed to read current environment: %w", err)
	}
	return len(currentEnv.Changes) > 0, nil
}

// ListEnvs lists all environment folders inside ".<AppName>"
func (e *envService) ListEnvs() ([]string, error) {
	dir := filepath.Join(e.workingDir, fmt.Sprintf(".%s", AppName))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := []string{}
	for _, e := range entries {
		if e.IsDir() {
			envs = append(envs, e.Name())
		}
	}
	return envs, nil
}
