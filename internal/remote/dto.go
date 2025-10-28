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
	Message string           `json:"message"`
	Code    string           `json:"code"`
	Data    PushResponseData `json:"data,omitempty"`
}

type PushResponseData struct {
	CommitHead string `json:"commitHead"`
}
