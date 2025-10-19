package handler

import (
	"context"

	"github.com/urfave/cli/v3"
)

type Status struct {
	changeRecordService changeRecordService
	slate               slate
}

func NewStatusHandler(changeRecordService changeRecordService, slate slate) *Status {
	return &Status{
		changeRecordService: changeRecordService,
		slate:               slate,
	}
}

func (h *Status) Handle(ctx context.Context, cmd *cli.Command) error {
	changes, err := h.changeRecordService.GetPendingChanges()
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		println("No pending changes")
		return nil
	}

	println("🔍 Pending Changes:")
	for _, change := range changes {
		h.slate.WriteStatus(change.Key, change.Action)
	}
	return nil
}
