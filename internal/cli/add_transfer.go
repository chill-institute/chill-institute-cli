package cli

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/chill-institute/chill-institute-cli/internal/rpc"
	"github.com/spf13/cobra"
)

func newAddTransferCommand(app *appContext) *cobra.Command {
	var transferURL string
	var rawRequest string
	var dryRun bool

	command := &cobra.Command{
		Use:   "add-transfer",
		Short: "Add a transfer through chill.institute",
		Example: strings.TrimSpace(`
chilly add-transfer --url "magnet:?xt=urn:btih:..."
printf '{"url":"magnet:?xt=urn:btih:..."}' | chilly add-transfer --json @- --output json
chilly add-transfer --url "magnet:?xt=urn:btih:..." --dry-run --output json
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddTransfer(app, "add-transfer", transferURL, rawRequest, dryRun)
		},
	}

	command.Flags().StringVar(&transferURL, "url", "", "magnet or URL to add as transfer")
	command.Flags().StringVar(&rawRequest, "json", "", "raw JSON request body, or @- to read it from stdin")
	command.Flags().BoolVar(&dryRun, "dry-run", false, "validate input and print the request without executing it")
	return command
}

func runAddTransfer(app *appContext, commandID string, transferURL string, rawRequest string, dryRun bool) error {
	request, err := buildAddTransferRequest(app, transferURL, rawRequest)
	if err != nil {
		return err
	}
	if dryRun {
		return app.writeDryRunPreview(commandID, procedureUserAddTransfer, rpc.AuthUser, request)
	}

	cfg, err := app.loadConfig()
	if err != nil {
		return err
	}
	token, err := app.userToken(cfg)
	if err != nil {
		return err
	}

	response, err := app.callRPC(
		context.Background(),
		cfg,
		procedureUserAddTransfer,
		request,
		rpc.AuthUser,
		token,
	)
	if err != nil {
		return fmt.Errorf("add transfer: %w", err)
	}
	return app.writeSelectedResponseBodyWithRenderer(response.Body, nil, renderTransferPretty)
}

func buildAddTransferRequest(app *appContext, transferURL string, rawRequest string) (map[string]any, error) {
	trimmedURL := strings.TrimSpace(transferURL)
	trimmedRequest := strings.TrimSpace(rawRequest)

	if trimmedURL != "" && trimmedRequest != "" {
		return nil, usageError("ambiguous_transfer_input", "use either --url or --json, not both")
	}

	if trimmedRequest != "" {
		request, err := app.decodeJSONObjectFlag(rawRequest, "--json")
		if err != nil {
			return nil, err
		}
		urlValue, ok := request["url"].(string)
		if !ok {
			return nil, usageError("invalid_json_payload", "--json payload must include a string url field")
		}
		normalizedURL, err := normalizeTransferURL(urlValue)
		if err != nil {
			return nil, err
		}
		request["url"] = normalizedURL
		return request, nil
	}

	normalizedURL, err := normalizeTransferURL(transferURL)
	if err != nil {
		return nil, err
	}
	return map[string]any{"url": normalizedURL}, nil
}

func normalizeTransferURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", usageError("missing_url", "--url is required")
	}
	if strings.IndexFunc(trimmed, unicode.IsControl) >= 0 {
		return "", usageError("invalid_url", "--url must not contain control characters")
	}
	return trimmed, nil
}
