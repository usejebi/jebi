package core

import "time"

type Project struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	DefaultEnvironment string    `json:"defaultEnvironment,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`

	Key string `json:"key,omitempty"` // Base64-encoded encryption key for the project
}

type Environment struct {
	Name      string    `json:"name"`
	ProjectID string    `json:"projectId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Secret struct {
	Key             string    `json:"key"`
	Value           string    `json:"value"`
	Nonce           string    `json:"nonce"`
	ProjectId       string    `json:"projectId"`
	EnvironmentName string    `json:"environmentName"`
	NoSecret        bool      `json:"nosecret"`
	UpdatedAt       time.Time `json:"updatedAt"`
	CreatedAt       time.Time `json:"createdAt"`
}

// Commit represents a single commit with its changes
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Changes   []Change  `json:"changes"`
	ParentID  string    `json:"parentId,omitempty"` // Empty for first commit

	ProjectID       string `json:"projectId,omitempty"`
	EnvironmentName string `json:"environmentName,omitempty"`
}

// Head represents the HEAD pointers for an environment
type Head struct {
	LocalHead  string `json:"localHead"`  // Latest local commit ID
	RemoteHead string `json:"remoteHead"` // Latest remote commit ID
}

type CurrentEnv struct {
	Env     string   `json:"env"`
	Changes []Change `json:"changes"`
}

// ChangeType represents the type of change made to a secret
type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "add"
	ChangeTypeRemove ChangeType = "remove"
	ChangeTypeModify ChangeType = "modify"
)

// Change represents a single change to a secret
// Design Note: This struct only stores the new value, not the old value.
// This follows Git-like behavior where commits are forward-looking.
// Old values can be computed by reconstructing state at the parent commit.
// This keeps the data structure simple while still supporting:
// - State reconstruction at any commit
// - Diff computation (by comparing parent state vs new value)
// - Rollback operations (by applying reverse changes)
// Design Note: This struct only stores the new value, not the old value.
// This follows Git-like behavior where commits are forward-looking.
// Old values can be computed by reconstructing state at the parent commit.
// This keeps the data structure simple while still supporting:
// - State reconstruction at any commit
// - Diff computation (by comparing parent state vs new value)
// - Rollback operations (by applying reverse changes)
type Change struct {
	Type     ChangeType `json:"type"`
	Key      string     `json:"key"`
	Value    string     `json:"value,omitempty"`    // Empty for remove operations; new value for add/modify
	Nonce    string     `json:"nonce,omitempty"`    // Nonce for encrypted secrets; empty for no-secret entries
	NoSecret bool       `json:"nosecret,omitempty"` // Whether the secret is a no-secret entry
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
	Username    string `json:"username,omitempty"`
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
}

type AuthResponse struct {
	Tokens Tokens `json:"tokens"`
	User   User   `json:"user"`
}
