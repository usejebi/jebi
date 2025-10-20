package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jawahars16/jebi/cmd"
	"github.com/jawahars16/jebi/internal/ui"
)

func main() {
	slate := ui.NewSlate(lipgloss.Color("82"))
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slate.ShowError(fmt.Sprintf("Error: %v", err))
		os.Exit(1)
	}
}
