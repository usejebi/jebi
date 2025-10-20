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

	var output string
	for i := len(commits) - 1; i >= 0; i-- { // reverse order (latest first)
		c := commits[i]
		output += fmt.Sprintf("commit %s \nDate: %s\n `%s` \n\n", c.ID, c.Timestamp.Format("Mon Jan 2 15:04:05 2006 -0700"), c.Message)
	}
	h.slate.RenderMarkdown(output)
	return nil
}
