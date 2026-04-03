package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func captureCmdStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func initTestConfigEnv(t *testing.T) {
	t.Helper()

	home := t.TempDir()
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", home); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("HOME", oldHome)
	})

	for _, envVar := range []string{"FLUX_DATASTREAM_API_KEY", "FLUX_RPC_API_KEY", "FLUX_RUGCHECK_API_KEY"} {
		envVar := envVar
		old := os.Getenv(envVar)
		if err := os.Unsetenv(envVar); err != nil {
			t.Fatalf("failed to clear %s: %v", envVar, err)
		}
		t.Cleanup(func() {
			_ = os.Setenv(envVar, old)
		})
	}
}

func initTestContext(t *testing.T) {
	t.Helper()
	oldCtx := ctx
	initContext()
	t.Cleanup(func() {
		ctx = oldCtx
	})
}

func TestFormatConfigSetConfirmation(t *testing.T) {
	if got := formatConfigSetConfirmation("fluxrpc.api_key", "secret"); got != "Set fluxrpc.api_key" {
		t.Fatalf("unexpected API key confirmation: %q", got)
	}
	if got := formatConfigSetConfirmation("fluxrpc.region", "eu"); got != "Set fluxrpc.region=eu" {
		t.Fatalf("unexpected normal confirmation: %q", got)
	}
}

func TestConfigSetCmd_DefaultOutputIsConcise(t *testing.T) {
	initTestConfigEnv(t)

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
	})

	out := captureCmdStdout(t, func() {
		rootCmd.SetArgs([]string{"config", "set", "fluxrpc.api_key", "secret-value"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute failed: %v", err)
		}
	})

	if strings.Contains(out, `"success"`) {
		t.Fatalf("expected concise output, got JSON: %s", out)
	}
	if strings.Contains(out, "secret-value") {
		t.Fatalf("expected secret to be redacted from output: %s", out)
	}
	if strings.TrimSpace(out) != "Set fluxrpc.api_key" {
		t.Fatalf("unexpected output: %q", out)
	}

	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".flux-cli", "config.json"))
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	if !strings.Contains(string(data), "secret-value") {
		t.Fatalf("expected config file to persist the value")
	}
}

func TestConfigSetCmd_ExplicitJSONFormatKeepsStructuredOutput(t *testing.T) {
	initTestConfigEnv(t)

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
	})

	out := captureCmdStdout(t, func() {
		rootCmd.SetArgs([]string{"--format", "json", "config", "set", "fluxrpc.region", "eu"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute failed: %v", err)
		}
	})

	if !strings.Contains(out, `"success": true`) {
		t.Fatalf("expected JSON output, got: %s", out)
	}
}

func TestHelpHidesConfiguredAPIKeyFlags(t *testing.T) {
	initTestConfigEnv(t)

	configDir := filepath.Join(os.Getenv("HOME"), ".flux-cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(`{"datastream":{"api_key":"ds-key"},"fluxrpc":{"api_key":"rpc-key","region":"us"},"rugcheck":{"api_key":"rc-key"}}`), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
		updateHiddenConfigFlags(rootCmd, nil)
	})

	out := captureCmdStdout(t, func() {
		rootCmd.SetArgs([]string{"data", "tokens", "price", "--help"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute failed: %v", err)
		}
	})

	for _, flag := range []string{"--datastream-api-key", "--fluxrpc-api-key", "--rugcheck-api-key"} {
		if strings.Contains(out, flag) {
			t.Fatalf("expected %s to be hidden, got help:\n%s", flag, out)
		}
	}
}

func TestHelpShowsAPIKeyFlagsWhenUnset(t *testing.T) {
	initTestConfigEnv(t)

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
		updateHiddenConfigFlags(rootCmd, nil)
	})

	out := captureCmdStdout(t, func() {
		rootCmd.SetArgs([]string{"data", "tokens", "price", "--help"})
		if err := rootCmd.Execute(); err != nil {
			t.Fatalf("execute failed: %v", err)
		}
	})

	if !strings.Contains(out, "--datastream-api-key") {
		t.Fatalf("expected datastream flag in help, got:\n%s", out)
	}
}

func TestRunConfigInit_SavesInteractiveConfig(t *testing.T) {
	initTestConfigEnv(t)
	initTestContext(t)

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
	})

	var out bytes.Buffer
	input := strings.NewReader("ds-key\nrpc-key\neu\nrc-key\ntable\n")

	if err := runConfigInit(configInitCmd, input, &out, false); err != nil {
		t.Fatalf("runConfigInit failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".flux-cli", "config.json"))
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	saved := string(data)
	for _, want := range []string{`"api_key": "ds-key"`, `"api_key": "rpc-key"`, `"region": "eu"`, `"api_key": "rc-key"`, `"format": "table"`} {
		if !strings.Contains(saved, want) {
			t.Fatalf("expected config to contain %s, got %s", want, saved)
		}
	}
}

func TestRunConfigInit_PreservesValuesOnSkip(t *testing.T) {
	initTestConfigEnv(t)

	oldPreRun := rootCmd.PersistentPreRun
	rootCmd.PersistentPreRun = nil
	t.Cleanup(func() {
		rootCmd.PersistentPreRun = oldPreRun
	})

	configDir := filepath.Join(os.Getenv("HOME"), ".flux-cli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	existing := `{"datastream":{"api_key":"ds-key"},"fluxrpc":{"api_key":"rpc-key","region":"eu"},"rugcheck":{"api_key":"rc-key","base_url":"https://api.rugcheck.xyz/v1"},"output":{"format":"table"}}`
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), []byte(existing), 0600); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	initTestContext(t)

	var out bytes.Buffer
	input := strings.NewReader("\n\n\n\n\n")

	if err := runConfigInit(configInitCmd, input, &out, true); err != nil {
		t.Fatalf("runConfigInit failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(configDir, "config.json"))
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	saved := string(data)
	for _, want := range []string{`"api_key": "ds-key"`, `"api_key": "rpc-key"`, `"region": "eu"`, `"api_key": "rc-key"`, `"format": "table"`} {
		if !strings.Contains(saved, want) {
			t.Fatalf("expected config to contain %s, got %s", want, saved)
		}
	}
}

func TestShouldSkipAutoSetup(t *testing.T) {
	cases := []struct {
		cmd  *cobra.Command
		want bool
	}{
		{rootCmd, true},
		{configSetCmd, true},
		{versionCmd, true},
		{dataTokensPriceCmd, false},
	}

	for _, tc := range cases {
		if got := shouldSkipAutoSetup(tc.cmd); got != tc.want {
			t.Fatalf("%s skip=%v want=%v", tc.cmd.CommandPath(), got, tc.want)
		}
	}
}
