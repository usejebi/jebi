package remote

import (
	"github.com/jawahars16/jebi/internal/core"
)

type PushRequest struct {
	Project        core.Project     `json:"project"`
	Environment    core.Environment `json:"environment"`
	Commits        []core.Commit    `json:"commits"`                  // New commits to push
	FinalState     []core.Secret    `json:"finalState"`               // Final computed secrets with all metadata
	RemoteHeadHash string           `json:"remoteHeadHash,omitempty"` // For conflict detection
}

type PushResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	IsFirstPush   bool   `json:"is_first_push"`
	CommitsPushed int    `json:"commits_pushed"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
