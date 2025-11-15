package core

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"sort"
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

// getCommitsPath returns the path to commits file for an environment
func (s *commitService) getCommitsPath(env string) string {
	return filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), EnvDirPath, env, CommitFileName)
}

// getHeadPath returns the path to HEAD file for an environment
func (s *commitService) getHeadPath(env string) string {
	return filepath.Join(s.workingDir, fmt.Sprintf(".%s", AppName), EnvDirPath, env, "HEAD")
}

// generateCommitID generates a unique commit ID based on timestamp and content
func (s *commitService) generateCommitID(message string, author string, timestamp time.Time) string {
	content := fmt.Sprintf("%s-%s-%d", message, author, timestamp.Unix())
	hash := sha1.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)[:12] // Use first 12 characters like Git
}

// AddCommit creates a new commit with the given changes
func (s *commitService) AddCommit(id, env, message, author string, changes []Change, timestamp time.Time) (*Commit, error) {
	// Load existing commits
	commits, err := s.loadCommits(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load existing commits: %w", err)
	}

	// Get current HEAD to set as parent
	head, err := s.GetHead(env)
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	if id == "" {
		id = s.generateCommitID(message, author, timestamp)
	}

	// Create new commit
	commit := &Commit{
		ID:        id,
		Message:   message,
		Author:    author,
		Timestamp: timestamp,
		Changes:   changes,
		ParentID:  head.LocalHead,
	}

	// Append to commits list
	commits = append(commits, *commit)

	// Save commits
	if err := s.saveCommits(env, commits); err != nil {
		return nil, fmt.Errorf("failed to save commits: %w", err)
	}

	// Update local HEAD
	if err := s.UpdateLocalHead(env, commit.ID); err != nil {
		return nil, fmt.Errorf("failed to update HEAD: %w", err)
	}

	return commit, nil
}

// GetCommit retrieves a specific commit by ID
func (s *commitService) GetCommit(env, commitID string) (*Commit, error) {
	commits, err := s.loadCommits(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load commits: %w", err)
	}

	for _, commit := range commits {
		if commit.ID == commitID {
			return &commit, nil
		}
	}

	return nil, fmt.Errorf("commit %s not found in environment %s", commitID, env)
}

// ListCommits returns all commits for an environment, sorted by timestamp (newest first)
func (s *commitService) ListCommits(env string) ([]Commit, error) {
	commits, err := s.loadCommits(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load commits: %w", err)
	}

	// Sort by timestamp, newest first
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.After(commits[j].Timestamp)
	})

	return commits, nil
}

// GetHead retrieves the HEAD pointers for an environment
func (s *commitService) GetHead(env string) (*Head, error) {
	head, err := io.ReadJSONFile[Head](s.getHeadPath(env))
	if err != nil {
		return &Head{}, fmt.Errorf("failed to read HEAD: %w", err)
	}
	return &head, nil
}

// UpdateLocalHead updates the local HEAD pointer
func (s *commitService) UpdateLocalHead(env, commitID string) error {
	head, err := s.GetHead(env)
	if err != nil {
		return fmt.Errorf("failed to get current HEAD: %w", err)
	}

	head.LocalHead = commitID

	if err := io.WriteJSONToFile(s.getHeadPath(env), head); err != nil {
		return fmt.Errorf("failed to update local HEAD: %w", err)
	}

	return nil
}

// UpdateRemoteHead updates the remote HEAD pointer
func (s *commitService) UpdateRemoteHead(env, commitID string) error {
	head, err := s.GetHead(env)
	if err != nil {
		return fmt.Errorf("failed to get current HEAD: %w", err)
	}

	head.RemoteHead = commitID

	if err := io.WriteJSONToFile(s.getHeadPath(env), head); err != nil {
		return fmt.Errorf("failed to update remote HEAD: %w", err)
	}

	return nil
}

// ComputeState computes the final state of secrets up to a specific commit
func (s *commitService) ComputeState(env, upToCommitID string) (map[string]Secret, error) {
	commits, err := s.loadCommits(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load commits: %w", err)
	}

	// Build commit chain up to the specified commit
	commitChain, err := s.buildCommitChain(commits, upToCommitID)
	if err != nil {
		return nil, fmt.Errorf("failed to build commit chain: %w", err)
	}

	// Apply changes in chronological order
	stateMap := make(map[string]Secret)
	for _, commit := range commitChain {
		for _, change := range commit.Changes {
			switch change.Type {
			case ChangeTypeAdd, ChangeTypeModify:
				stateMap[change.Key] = Secret{
					Key:      change.Key,
					Value:    change.Value,
					Nonce:    change.Nonce,
					NoSecret: change.NoSecret,
				}
			case ChangeTypeRemove:
				delete(stateMap, change.Key)
			}
		}
	}

	return stateMap, nil
}

// GetCommitsSinceRemoteHead returns all commits since the remote HEAD
func (s *commitService) GetCommitsSinceRemoteHead(env string) ([]Commit, error) {
	head, err := s.GetHead(env)
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	// If no remote HEAD, return all commits
	if head.RemoteHead == "" {
		return s.ListCommits(env)
	}

	commits, err := s.loadCommits(env)
	if err != nil {
		return nil, fmt.Errorf("failed to load commits: %w", err)
	}

	// Find commits since remote HEAD
	var newCommits []Commit
	remoteHeadFound := false

	// Sort commits by timestamp (oldest first for this operation)
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Timestamp.Before(commits[j].Timestamp)
	})

	for _, commit := range commits {
		if commit.ID == head.RemoteHead {
			remoteHeadFound = true
			continue
		}
		if remoteHeadFound {
			newCommits = append(newCommits, commit)
		}
	}

	// If remote HEAD not found, return all commits
	if !remoteHeadFound {
		return commits, nil
	}

	return newCommits, nil
}

// loadCommits loads commits from disk
func (s *commitService) loadCommits(env string) ([]Commit, error) {
	path := s.getCommitsPath(env)
	commits, err := io.ReadJSONFile[[]Commit](path)
	if err != nil {
		return nil, fmt.Errorf("failed to read commits file: %w", err)
	}
	return commits, nil
}

// saveCommits saves commits to disk
func (s *commitService) saveCommits(env string, commits []Commit) error {
	if err := io.WriteJSONToFile(s.getCommitsPath(env), commits); err != nil {
		return fmt.Errorf("failed to write commits file: %w", err)
	}
	return nil
}

// buildCommitChain builds a chronological chain of commits up to a specific commit
func (s *commitService) buildCommitChain(allCommits []Commit, upToCommitID string) ([]Commit, error) {
	if upToCommitID == "" {
		return []Commit{}, nil
	}

	// Create a map for fast lookup
	commitMap := make(map[string]Commit)
	for _, commit := range allCommits {
		commitMap[commit.ID] = commit
	}

	// Build chain by following parent links backwards
	var chain []Commit
	currentID := upToCommitID

	for currentID != "" {
		commit, exists := commitMap[currentID]
		if !exists {
			return nil, fmt.Errorf("commit %s not found", currentID)
		}

		chain = append([]Commit{commit}, chain...) // Prepend to maintain chronological order
		currentID = commit.ParentID
	}

	return chain, nil
}
