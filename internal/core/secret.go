package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jawahars16/jebi/internal/io"
)

var (
	ErrSecretAlreadyExists = fmt.Errorf("secret already exists")
	ErrSecretNotFound      = fmt.Errorf("secret not found")
)

type secretService struct {
	workingDir string
}

func NewSecretService(workingDir string) *secretService {
	return &secretService{
		workingDir: workingDir,
	}
}

func (s *secretService) envDir(env string) string {
	return filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), EnvDirPath, env)
}

func (s *secretService) AddSecret(key, env string, secret Secret) error {
	envDir := s.envDir(env)
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
			return ErrSecretAlreadyExists
		}
	}

	data[key] = secret

	err := io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return fmt.Errorf("failed to write secrets: %w", err)
	}
	return nil
}

func (s *secretService) SetSecret(key, env string, secret Secret) (ChangeType, error) {
	envDir := s.envDir(env)
	secretPath := filepath.Join(envDir, SecretFileName)

	data, err := io.ReadJSONFile[map[string]Secret](secretPath)
	if err != nil {
		return "", fmt.Errorf("failed to read secrets: %w", err)
	}

	var action ChangeType
	_, exists := data[key]
	if !exists {
		action = ChangeTypeAdd
	} else {
		action = ChangeTypeModify
	}

	data[key] = secret
	err = io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return "", fmt.Errorf("failed to write secrets: %w", err)
	}
	return action, nil
}

func (s *secretService) RemoveSecret(key, env string) error {
	envDir := s.envDir(env)
	secretPath := filepath.Join(envDir, SecretFileName)

	data, err := io.ReadJSONFile[map[string]Secret](secretPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets: %w", err)
	}

	if _, exists := data[key]; !exists {
		return ErrSecretNotFound
	}

	delete(data, key)

	err = io.WriteJSONToFile(secretPath, data)
	if err != nil {
		return fmt.Errorf("failed to write secrets: %w", err)
	}
	return nil
}

func (s *secretService) ListSecrets(projectId, env string) ([]Secret, error) {
	envDir := s.envDir(env)
	secretPath := filepath.Join(envDir, SecretFileName)

	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		// If secret file does not exist, return empty slice
		return []Secret{}, nil
	}

	data, err := io.ReadJSONFile[map[string]Secret](secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets: %w", err)
	}

	var secrets []Secret
	for key, secret := range data {
		secrets = append(secrets, Secret{
			Key:             key,
			Value:           secret.Value,
			Nonce:           secret.Nonce,
			ProjectId:       projectId,
			EnvironmentName: env,
			NoSecret:        secret.NoSecret,
			CreatedAt:       secret.CreatedAt,
			UpdatedAt:       secret.UpdatedAt,
		})
	}
	return secrets, nil
}
