package handler

import (
	"context"
	"fmt"

	"github.com/jawahars16/jebi/internal/core"
	"github.com/urfave/cli/v3"
)

type Login struct {
	userService userService
	slate       slate
}

func NewLoginHandler(userService userService, slate slate) *Login {
	return &Login{
		userService: userService,
		slate:       slate,
	}
}

func (h *Login) Handle(ctx context.Context, cmd *cli.Command) error {
	h.slate.ShowHeader("Opening browser window for authentication...")
	h.slate.RenderMarkdown(fmt.Sprintf(`A browser window will open for you to authenticate with Jebi.
	Please complete the login process in your browser.
	The CLI will wait for up to 30 seconds for authentication to complete.
	(Click this link if not redirected automatically)
	<%s>`, core.LoginURL))

	// Attempt browser-based authentication
	// The userService handles saving all authentication details internally
	authResult, err := h.userService.AuthenticateWithBrowser(core.LoginURL)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Authentication failed: %v", err))
		return nil
	}

	fmt.Printf("Successfully authenticated as %s\n", authResult.User.DisplayName)

	return nil
}
