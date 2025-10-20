package handler

import "github.com/jawahars16/jebi/internal/core"

type appService interface {
	CreateAppDir() error
	Exists() (bool, error)
}

type projectService interface {
	SaveProjectConfig(name, description string) (string, error)
	LoadProjectConfig() (*core.Project, error)
}

type cryptService interface {
	GenerateKey() (string, error)
	Encrypt(key []byte, plaintext string) (ciphertextB64, nonceB64 string, err error)
	Decrypt(key []byte, ciphertextB64, nonceB64 string) (string, error)
	LoadSecrets(env string) (map[string]string, error)
	SaveKey(key string) error
	LoadKey() ([]byte, error)
}

type envService interface {
	CreateEnv(env string) error
	ListEnvs() ([]string, error)
	CurrentEnv() (string, error)
	SetCurrentEnv(env string) error
	GetCurrentEnv() (*core.CurrentEnv, error)
	RemoveEnv(env string) error
}

type secretService interface {
	SetSecret(key, env string, secret core.Secret) (string, error)
	AddSecret(key, env string, secret core.Secret) error
	RemoveSecret(key, env string) error
}

type commitService interface {
	AddCommit(env, message string, changes []core.Change) error
	ListCommits(env string) ([]core.Commit, error)
}

type changeRecordService interface {
	AddChangeRecord(env, action, key string) error
	ClearPendingChanges() error
}

type slate interface {
	PromptWithDefault(message, defaultValue string) string
	ShowHeader(title string)
	ShowList(title string, items []string, highlight string)
	WriteStatus(changes []core.Change)
	RenderMarkdown(md string)
	ShowWarning(msg string)
	ShowError(msg string)
}
