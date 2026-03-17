package cli

import (
	"encoding/json"
	"io"
	"strings"
)

func (app *appContext) decodeJSONObjectFlag(raw string, flagName string) (map[string]any, error) {
	payload, err := app.readJSONFlag(raw, flagName)
	if err != nil {
		return nil, err
	}

	var value any
	if err := json.Unmarshal(payload, &value); err != nil {
		return nil, usageError("invalid_json_payload", "parse %s payload: %v", flagName, err)
	}

	object, ok := value.(map[string]any)
	if !ok {
		return nil, usageError("invalid_json_payload", "%s payload must be a JSON object", flagName)
	}
	return object, nil
}

func (app *appContext) readJSONFlag(raw string, flagName string) ([]byte, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, usageError("missing_json_payload", "%s payload cannot be empty", flagName)
	}

	if trimmed != "@-" {
		return []byte(trimmed), nil
	}

	if app == nil || app.stdin == nil {
		return nil, usageError("missing_json_payload", "%s payload stdin is not available", flagName)
	}

	payload, err := io.ReadAll(app.stdin)
	if err != nil {
		return nil, wrapInternalError("json_payload_read_failed", "read JSON payload from stdin", err)
	}
	if strings.TrimSpace(string(payload)) == "" {
		return nil, usageError("missing_json_payload", "%s payload from stdin cannot be empty", flagName)
	}
	return payload, nil
}
