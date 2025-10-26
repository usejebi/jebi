package cmd

import (
	"context"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/jawahars16/jebi/internal/core"
	"github.com/jawahars16/jebi/internal/crypt"
	"github.com/jawahars16/jebi/internal/handler"
	"github.com/jawahars16/jebi/internal/remote"
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
	userService := core.NewUserService(workingDir)

	slate := ui.NewSlate(lipgloss.Color("82"))

	setHandler := handler.NewSetHandler(projectService, cryptService, envService, secretService, changeRecordService, slate)
	addHandler := handler.NewAddHandler(projectService, cryptService, envService, secretService, changeRecordService, slate)
	removeHandler := handler.NewRemoveHandler(cryptService, envService, secretService, changeRecordService, slate)
	projectHandler := handler.NewInitHandler(appService, projectService, envService, cryptService, slate)
	envHandler := handler.NewEnvHandler(envService, slate)
	commitHandler := handler.NewCommitHandler(envService, commitService, changeRecordService, userService, secretService, projectService, slate)
	exportHandler := handler.NewExportHandler(envService, cryptService, projectService, slate)
	statusHandler := handler.NewStatusHandler(envService, slate)
	runHandler := handler.NewRunHandler(envService, cryptService, projectService, slate)
	logHandler := handler.NewLogHandler(envService, commitService, slate)
	loginHandler := handler.NewLoginHandler(userService, slate)
	apiClient := remote.NewAPIClient(core.DefaultServerURL)
	pushHandler := handler.NewPushHandler(projectService, envService, secretService, commitService, apiClient, slate)

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
		newLoginCommand(loginHandler),
		newPushCommand(pushHandler),
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
