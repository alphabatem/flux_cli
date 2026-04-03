package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/alphabatem/flux_cli/dto"
)

func TestDataTokensPriceMissingArgError(t *testing.T) {
	err := dataTokensPriceCmd.Args(dataTokensPriceCmd, []string{})
	if err == nil {
		t.Fatal("expected missing arg error")
	}

	var usageErr *usageError
	if !errors.As(err, &usageErr) {
		t.Fatalf("expected usageError, got %T", err)
	}

	if !strings.Contains(err.Error(), "missing required argument: mint") {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(err.Error(), "Usage: flux data tokens price <mint>") {
		t.Fatalf("missing usage line: %s", err)
	}
}

func TestConfigSetMissingSecondArgError(t *testing.T) {
	err := configSetCmd.Args(configSetCmd, []string{"fluxrpc.api_key"})
	if err == nil {
		t.Fatal("expected missing arg error")
	}

	if !strings.Contains(err.Error(), "missing required argument: value") {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(err.Error(), "Usage: flux config set <key> <value>") {
		t.Fatalf("missing usage line: %s", err)
	}
}

func TestRPCCallTooManyArgsError(t *testing.T) {
	err := rpcCallCmd.Args(rpcCallCmd, []string{"getBalance", "[]", "extra"})
	if err == nil {
		t.Fatal("expected too many args error")
	}

	if !strings.Contains(err.Error(), "too many arguments: got 3, want 1-2") {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(err.Error(), "Usage: flux rpc call <method> [params-json]") {
		t.Fatalf("missing usage line: %s", err)
	}
}

func TestRootErrorJSONDoesNotEscapeAngleBrackets(t *testing.T) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(dto.CLIResponse{
		Success: false,
		Error: &dto.CLIError{
			Code:    "USAGE_ERROR",
			Message: "missing required argument: mint\nUsage: flux data tokens candles <mint>",
		},
	}); err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, `\u003c`) || strings.Contains(out, `\u003e`) {
		t.Fatalf("expected angle brackets to remain unescaped, got: %s", out)
	}
}
