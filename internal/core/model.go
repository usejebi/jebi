package core

import "time"

const (
	ActionAdd    = "add"
	ActionUpdate = "update"
	ActionRemove = "remove"
)

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Secret struct {
	Value     string    `json:"value"`
	Nonce     string    `json:"nonce"`
	UpdatedAt time.Time `json:"updated_at"`
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
	Username string `json:"username"`
	Server   string `json:"server"`
}

type AuthResult struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
}
