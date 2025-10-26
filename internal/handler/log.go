package handler

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

type Log struct {
	envService    envService
	commitService commitService
	slate         slate
}

func NewLogHandler(envService envService, commitService commitService, slate slate) *Log {
	return &Log{
		envService:    envService,
		commitService: commitService,
		slate:         slate,
	}
}

func (h *Log) Handle(ctx context.Context, cmd *cli.Command) error {
	env, err := h.envService.CurrentEnv()
	if err != nil {
		fmt.Printf("⚠️ Current environment is not set. Use `jebi use <env name>` to set the current environment\n")
		return nil
	}

	commits, err := h.commitService.ListCommits(env)
	if err != nil {
		return err
	}

	if len(commits) == 0 {
		fmt.Println("(no commits yet)")
		return nil
	}

	head, err := h.commitService.GetHead(env)
	if err != nil {
		return err
	}

	var output string
	for _, c := range commits {
		if c.ID == head.LocalHead {
			output += fmt.Sprintf("commit %s  (HEAD -> %s) \nAuthor: %s\nDate: %s\n `%s` \n\n", c.ID, env, c.Author, c.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"), c.Message)
			continue
		}
		if c.ID == head.RemoteHead {
			output += fmt.Sprintf("commit %s  (origin/%s) \nAuthor: %s\nDate: %s\n `%s` \n\n", c.ID, env, c.Author, c.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"), c.Message)
			continue
		}
		output += fmt.Sprintf("commit %s \nAuthor: %s\nDate: %s\n `%s` \n\n", c.ID, c.Author, c.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"), c.Message)
	}
	h.slate.RenderMarkdown(output)
	return nil
}
