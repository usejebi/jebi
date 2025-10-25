package remote

import "github.com/jawahars16/jebi/internal/core"

type PushRequest struct {
	Project     core.Project     `json:"project"`
	Environment core.Environment `json:"environment"`
	Secrets     []core.Secret    `json:"secrets"`
}

type PushResponse struct {
}
