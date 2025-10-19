package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/crypt"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/jawahars16/jebi/internal/ui"
	"github.com/urfave/cli/v3"
)

func initializeCommands() []*cli.Command {
	workingDir := getWorkingDir()
	appService := core.NewAppService(workingDir)
	projectService := core.NewProjectService(workingDir)
	envService := core.NewEnvService(workingDir)
	cryptService := crypt.NewService(workingDir)
	secretService := core.NewSecretService(workingDir)
	commitService := core.NewCommitService(workingDir)
	changeRecordService := core.NewChangeRecordService(workingDir)

	slate := ui.NewSlate(lipgloss.Color("82"))

	setHandler := handler.NewSetHandler(cryptService, envService, secretService, changeRecordService)
	addHandler := handler.NewAddHandler(cryptService, envService, secretService, changeRecordService)
	removeHandler := handler.NewRemoveHandler(cryptService, envService, secretService, changeRecordService)
	projectHandler := handler.NewInitHandler(appService, projectService, envService, cryptService, slate)
	envHandler := handler.NewEnvHandler(envService, slate)
	commitHandler := handler.NewCommitHandler(envService, commitService, changeRecordService)
	exportHandler := handler.NewExportHandler(envService, cryptService, slate)
	statusHandler := handler.NewStatusHandler(envService, slate)
	runHandler := handler.NewRunHandler(envService, cryptService, slate)
	logHandler := handler.NewLogHandler(envService, commitService, slate)

	return []*cli.Command{
		newInitCommand(projectHandler),
		newSetCommand(setHandler),
		newAddCommand(addHandler),
		newRemoveCommand(removeHandler),
		newEnvCommand(envHandler),
		newCommitCommand(commitHandler),
		newExportCommand(exportHandler),
		newLogCommand(logHandler),
		newStatusCommand(statusHandler),
		newRunCommand(runHandler),
		newVersionCommand(),
	}
}

func Run(ctx context.Context, args []string) error {
	cmd := &cli.Command{
		Name:        core.AppName,
		Usage:       "A demo CLI built with urfave/cli/v3",
		Description: "Dummy description for CLI",
		Version:     "0.1.0",
		Commands:    initializeCommands(),
	}

	return cmd.Run(ctx, args)
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("failed to get working directory")
	}
	return dir
}
