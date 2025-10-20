package core

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/jawahars16/jebi/internal/io"
)

type commitService struct {
	workingDir string
}

func NewCommitService(workingDir string) *commitService {
	return &commitService{
		workingDir: workingDir,
	}
}

// CommitList represents all commits for an environment.
type CommitList struct {
	Commits []Commit `json:"commits"`
}

func (s *commitService) commitFilePath(env string) string {
	return filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), EnvDirPath, env, CommitFileName)
}

// AddCommit appends a new commit to the environment.
func (s *commitService) AddCommit(env, message string, changes []Change) error {
	path := s.commitFilePath(env)

	cl, err := io.ReadJSONFile[CommitList](path)
	if err != nil {
		cl = CommitList{
			Commits: []Commit{},
		}
	}

	commit := Commit{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Message:   message,
		Timestamp: time.Now(),
		Changes:   changes,
	}
	cl.Commits = append(cl.Commits, commit)

	err = io.WriteJSONToFile(path, cl)
	if err != nil {
		return fmt.Errorf("failed to write commits: %w", err)
	}

	return nil
}

// ListCommits lists commits for an environment.
func (s *commitService) ListCommits(env string) ([]Commit, error) {
	path := s.commitFilePath(env)
	cl, err := io.ReadJSONFile[CommitList](path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse commits: %w", err)
	}
	return cl.Commits, nil
}
