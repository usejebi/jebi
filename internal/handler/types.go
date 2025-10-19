package handler

import "github.com/jawahars16/jebi/internal/core"

type appService interface {
	CreateAppDir() error
}

type projectService interface {
	SaveProjectConfig(name, description string) error
}

type cryptService interface {
	GenerateKey() (string, error)
	SaveKey(encodedKey string) error
	GetKey() ([]byte, error)
	Encrypt(key []byte, plaintext string) (ciphertextB64, nonceB64 string, err error)
	Decrypt(key []byte, ciphertextB64, nonceB64 string) (string, error)
	LoadSecrets(env string) (map[string]string, error)
}

type envService interface {
	CreateEnv(env string) error
	ListEnvs() ([]string, error)
	CurrentEnv() (string, error)
	SetCurrentEnv(env string) error
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
	GetPendingChanges() ([]core.Change, error)
	ClearPendingChanges() error
}

type slate interface {
	PromptWithDefault(message, defaultValue string) string
	ShowHeader(title string)
	ShowList(title string, items []string, highlight string)
	WriteStatus(key, action string)
	RenderMarkdown(md string) (string, error)
}
