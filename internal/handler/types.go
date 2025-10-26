package handler

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/remote"
	"github.com/jawahars16/jebi/internal/ui"
)

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
	LoadSecrets(project, env string) (map[string]string, error)
	SaveKey(key, project string) error
	LoadKey(project string) ([]byte, error)
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
	ListSecrets(projectId, env string) ([]core.Secret, error)
	RemoveSecret(key, env string) error
}

type changeRecordService interface {
	AddChangeRecord(env, action, key, value, nonce string, nosecret bool) error
	ClearPendingChanges() error
}

type userService interface {
	AuthenticateWithBrowser(server string) (*core.AuthResponse, error)
	SaveAuthToken(token string) error
	LoadAuthToken() (string, error)
	SaveCurrentUser(user core.User) error
	LoadCurrentUser() (*core.User, error)
	Logout() error
}

type commitService interface {
	// Commit operations
	AddCommit(env, message, author string, changes []core.Change) (*core.Commit, error)
	GetCommit(env, commitID string) (*core.Commit, error)
	ListCommits(env string) ([]core.Commit, error)

	// HEAD operations
	GetHead(env string) (*core.Head, error)
	UpdateLocalHead(env, commitID string) error
	UpdateRemoteHead(env, commitID string) error

	// Status and state operations
	ComputeState(env, upToCommitID string) (map[string]core.Secret, error)
	GetCommitsSinceRemoteHead(env string) ([]core.Commit, error)
}

type apiClient interface {
	Push(req remote.PushRequest) (remote.PushResponse, error)
}

type slate interface {
	PromptWithDefault(message, defaultValue string) string
	ShowHeader(title string)
	ShowList(title string, items []string, highlight string)
	WriteStatus(changes []core.Change)
	RenderMarkdown(md string)
	ShowWarning(msg string)
	ShowError(msg string)
	WriteStyledText(text string, options ui.StyleOptions)
	WriteColoredText(text string, color lipgloss.Color)
	WriteIndentedText(text string, options ui.StyleOptions)
	ShowSuccess(message string)
	ShowEnvironmentContext(env string)
	ShowSecretOperation(operation, key, env string, isPlaintext bool)
}
