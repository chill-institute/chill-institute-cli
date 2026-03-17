package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chill-institute/chill-institute-cli/internal/config"
	"github.com/chill-institute/chill-institute-cli/internal/rpc"
	"github.com/spf13/cobra"
)

func newAuthCommand(app *appContext) *cobra.Command {
	command := &cobra.Command{
		Use:   "auth",
		Short: "Sign in and manage local auth tokens",
		Example: strings.TrimSpace(`
chilly auth login
chilly auth login --token "token-from-setup"
chilly auth logout --dry-run --output json
`),
	}

	command.AddCommand(newAuthLoginCommand(app))
	command.AddCommand(newAuthLogoutCommand(app))
	return command
}

func newAuthLoginCommand(app *appContext) *cobra.Command {
	var token string
	var rawRequest string
	var skipOpen bool
	var skipVerify bool
	var dryRun bool

	command := &cobra.Command{
		Use:   "login",
		Short: "Sign in through a browser or store a setup token",
		Long: strings.TrimSpace(`
Sign in with the hosted browser flow, or store a setup token directly
with --token for non-interactive or remote environments. If the browser
is running on another machine, open https://chill.institute/auth/cli-token
in a signed-in browser and copy the token into --token.
`),
		Example: strings.TrimSpace(`
chilly auth login
chilly auth login --no-browser
chilly auth login --token "token-from-setup"
printf '{"token":"token-from-setup","skip_verify":true}' | chilly auth login --json @- --dry-run --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := app.loadConfig()
			if err != nil {
				return err
			}

			resolvedToken, resolvedSkipOpen, resolvedSkipVerify, err := resolveAuthLoginInput(app, token, rawRequest, skipOpen, skipVerify)
			if err != nil {
				return err
			}
			if dryRun {
				request := map[string]any{
					"api_base_url": cfg.APIBaseURL,
					"skip_verify":  resolvedSkipVerify,
				}
				if strings.TrimSpace(resolvedToken) != "" {
					request["mode"] = "token"
					request["token_provided"] = true
				} else {
					request["mode"] = "browser"
					request["no_browser"] = resolvedSkipOpen
				}
				return app.writeLocalDryRunPreview("auth login", request)
			}
			if resolvedToken == "" {
				resolvedToken, err = app.loginWithBrowser(context.Background(), cfg, resolvedSkipOpen)
				if err != nil {
					return err
				}
			}

			if resolvedToken == "" {
				return usageError("empty_token", "token cannot be empty")
			}

			var verifyResponse rpc.CallResponse
			if !resolvedSkipVerify {
				verifyResponse, err = app.callRPC(
					context.Background(),
					cfg,
					procedureUserGetUserProfile,
					map[string]any{},
					rpc.AuthUser,
					resolvedToken,
				)
				if err != nil {
					return fmt.Errorf("verify auth token: %w", err)
				}
			}

			cfg.AuthToken = resolvedToken
			if err := app.saveConfig(cfg); err != nil {
				return fmt.Errorf("persist auth token: %w", err)
			}

			if resolvedSkipVerify {
				return app.writeJSONPayload(map[string]any{"status": "ok", "saved": true})
			}
			return app.writeJSONPayload(map[string]any{
				"status":     "ok",
				"saved":      true,
				"request_id": verifyResponse.RequestID,
				"user":       jsonMessage(verifyResponse.Body),
			})
		},
	}

	command.Flags().StringVar(&token, "token", "", "setup token to store (non-interactive)")
	command.Flags().StringVar(&rawRequest, "json", "", "raw JSON request body, or @- to read it from stdin")
	command.Flags().BoolVar(&skipOpen, "no-browser", false, "print the login URL instead of opening a browser automatically")
	command.Flags().BoolVar(&skipVerify, "skip-verify", false, "skip token verification call")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "preview the auth action without verifying or saving credentials")
	return command
}

func resolveAuthLoginInput(app *appContext, token string, rawRequest string, skipOpen bool, skipVerify bool) (string, bool, bool, error) {
	trimmedToken := strings.TrimSpace(token)
	trimmedRequest := strings.TrimSpace(rawRequest)

	if trimmedRequest == "" {
		return trimmedToken, skipOpen, skipVerify, nil
	}
	if trimmedToken != "" || skipOpen || skipVerify {
		return "", false, false, usageError("ambiguous_auth_login_input", "use either flags or --json for auth login input, not both")
	}

	payload, err := app.decodeJSONObjectFlag(rawRequest, "--json")
	if err != nil {
		return "", false, false, err
	}

	resolvedSkipOpen := false
	if value, ok := payload["no_browser"]; ok {
		typed, ok := value.(bool)
		if !ok {
			return "", false, false, usageError("invalid_json_payload", "--json payload field no_browser must be a boolean")
		}
		resolvedSkipOpen = typed
	}

	resolvedSkipVerify := false
	if value, ok := payload["skip_verify"]; ok {
		typed, ok := value.(bool)
		if !ok {
			return "", false, false, usageError("invalid_json_payload", "--json payload field skip_verify must be a boolean")
		}
		resolvedSkipVerify = typed
	}

	resolvedToken := ""
	if value, ok := payload["token"]; ok {
		typed, ok := value.(string)
		if !ok {
			return "", false, false, usageError("invalid_json_payload", "--json payload field token must be a string")
		}
		resolvedToken = strings.TrimSpace(typed)
	}

	return resolvedToken, resolvedSkipOpen, resolvedSkipVerify, nil
}

func (app *appContext) loginWithBrowser(ctx context.Context, cfg config.Config, skipOpen bool) (string, error) {
	flow, err := newLoopbackAuthFlow(cfg.APIBaseURL)
	if err != nil {
		return "", err
	}
	defer func() {
		if shutdownErr := flow.shutdown(); shutdownErr != nil {
			_, _ = fmt.Fprintf(app.stderr, "Unable to stop local auth server cleanly: %v\n", shutdownErr)
		}
	}()

	noticeWriter := app.stdout
	if app.opts.output == outputJSON {
		noticeWriter = app.stderr
	}
	if _, err := fmt.Fprintf(noticeWriter, "Open this URL to authenticate:\n%s\n\n", flow.loginURL); err != nil {
		return "", err
	}

	errCh := make(chan error, 1)
	flow.start(errCh)

	if !skipOpen {
		if err := app.openURL(flow.loginURL); err != nil {
			_, _ = fmt.Fprintf(app.stderr, "Unable to open browser automatically: %v\n", err)
		}
	}

	timeout := app.authFlowTimeout
	if timeout <= 0 {
		timeout = 2 * time.Minute
	}

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	token, err := flow.waitForToken(waitCtx, errCh)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(token), nil
}

func newAuthLogoutCommand(app *appContext) *cobra.Command {
	var dryRun bool
	var rawRequest string

	command := &cobra.Command{
		Use:   "logout",
		Short: "Clear the stored auth token from local config",
		Example: strings.TrimSpace(`
chilly auth logout
chilly auth logout --dry-run --output json
chilly auth logout --json '{"clear_auth_token":true}' --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := app.loadConfig()
			if err != nil {
				return err
			}
			request, err := resolveAuthLogoutInput(app, rawRequest, strings.TrimSpace(cfg.AuthToken) != "")
			if err != nil {
				return err
			}
			if dryRun {
				return app.writeLocalDryRunPreview("auth logout", request)
			}
			cfg.AuthToken = ""
			if err := app.saveConfig(cfg); err != nil {
				return fmt.Errorf("persist config: %w", err)
			}
			return app.writeJSONPayload(map[string]any{"status": "ok", "logged_out": true})
		},
	}

	command.Flags().StringVar(&rawRequest, "json", "", "raw JSON request body, or @- to read it from stdin")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "preview the local auth token removal without saving it")
	return command
}

func resolveAuthLogoutInput(app *appContext, rawRequest string, hadAuthToken bool) (map[string]any, error) {
	request := map[string]any{
		"clear_auth_token": true,
		"had_auth_token":   hadAuthToken,
	}

	if strings.TrimSpace(rawRequest) == "" {
		return request, nil
	}

	payload, err := app.decodeJSONObjectFlag(rawRequest, "--json")
	if err != nil {
		return nil, err
	}
	if value, ok := payload["clear_auth_token"]; ok {
		typed, ok := value.(bool)
		if !ok || !typed {
			return nil, usageError("invalid_json_payload", "--json payload field clear_auth_token must be true")
		}
	}
	return request, nil
}

func jsonMessage(raw []byte) any {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return map[string]any{"raw": string(raw)}
	}
	return value
}
