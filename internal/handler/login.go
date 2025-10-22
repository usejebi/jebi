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
	authResult, err := h.userService.AuthenticateWithBrowser(core.LoginURL)
	if err != nil {
		h.slate.ShowError(fmt.Sprintf("Authentication failed: %v", err))
		return nil
	}

	// Save the authentication tokens
	if err := h.userService.SaveAuthToken(authResult.AccessToken); err != nil {
		return fmt.Errorf("failed to save authentication token: %w", err)
	}

	// Store the current user information
	user := core.User{
		Username: authResult.Username,
		Server:   core.LoginURL,
	}
	if err := h.userService.SaveCurrentUser(user); err != nil {
		return fmt.Errorf("failed to save user information: %w", err)
	}

	fmt.Printf("✅ Successfully authenticated as '%s'\n", authResult.Username)
	fmt.Printf("🔒 Authentication token saved securely\n")
	return nil
}
