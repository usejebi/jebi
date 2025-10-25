package core

import "time"

const (
	ActionAdd    = "add"
	ActionUpdate = "update"
	ActionRemove = "remove"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Environment struct {
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Secret struct {
	Key             string    `json:"key"`
	Value           string    `json:"value"`
	Nonce           string    `json:"nonce"`
	ProjectId       string    `json:"project_id"`
	EnvironmentName string    `json:"environment_name"`
	NoSecret        bool      `json:"nosecret"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedAt       time.Time `json:"created_at"`
}

// Commit represents a single commit with its changes
type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	Changes   []Change  `json:"changes"`
	ParentID  string    `json:"parent_id,omitempty"` // Empty for first commit
}

// Head represents the HEAD pointers for an environment
type Head struct {
	LocalHead  string `json:"local_head"`  // Latest local commit ID
	RemoteHead string `json:"remote_head"` // Latest remote commit ID
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
	Type  ChangeType `json:"type"`
	Key   string     `json:"key"`
	Value string     `json:"value,omitempty"` // Empty for remove operations; new value for add/modify
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
