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

type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Changes   []Change  `json:"changes"`
}

type CurrentEnv struct {
	Env     string   `json:"env"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Action string `json:"action"`
	Key    string `json:"key"`
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
