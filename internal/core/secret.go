package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/io"
)

type secretService struct {
	workingDir string
}

func NewSecretService(workingDir string) *secretService {
	return &secretService{
		workingDir: workingDir,
	}
}

func (s *secretService) AddSecret(key, env string, secret Secret) error {
	envDir := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), env)
	secretPath := filepath.Join(envDir, SecretFileName)

	var data map[string]Secret
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		// If secret file does not exist, create an empty one
		data = make(map[string]Secret)
	} else {
		data, err = io.ReadJSONFile[map[string]Secret](secretPath)
		if err != nil {
			return fmt.Errorf("failed to read secrets: %w", err)
		}
		if _, exists := data[key]; exists {
			return fmt.Errorf("secret with key %q already exists", key)
		}
	}

	data[key] = secret

	err := io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return fmt.Errorf("failed to write secrets: %w", err)
	}
	return nil
}

func (s *secretService) SetSecret(key, env string, secret Secret) (string, error) {
	envDir := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), env)
	secretPath := filepath.Join(envDir, SecretFileName)

	data, err := io.ReadJSONFile[map[string]Secret](secretPath)
	if err != nil {
		return "", fmt.Errorf("failed to read secrets: %w", err)
	}

	var action string
	_, exists := data[key]
	if !exists {
		action = ActionAdd
	} else {
		action = ActionUpdate
	}

	data[key] = secret
	err = io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return "", fmt.Errorf("failed to write secrets: %w", err)
	}
	return action, nil
}

func (s *secretService) RemoveSecret(key, env string) error {
	envDir := filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), env)
	secretPath := filepath.Join(envDir, SecretFileName)

	data, err := io.ReadJSONFile[map[string]Secret](secretPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets: %w", err)
	}

	if _, exists := data[key]; !exists {
		return fmt.Errorf("secret with key %q does not exist", key)
	}

	delete(data, key)

	err = io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return fmt.Errorf("failed to write secrets: %w", err)
	}
	return nil
}
