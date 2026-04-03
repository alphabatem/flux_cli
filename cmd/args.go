package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type usageError struct {
	message string
}

func (e *usageError) Error() string {
	return e.message
}

func exactArgsFromUse() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		required, _ := parseUseArgs(cmd.Use)
		expected := len(required)

		if len(args) < expected {
			return newMissingArgsError(cmd, required[len(args):])
		}
		if len(args) > expected {
			return newUsageError(cmd, fmt.Sprintf("too many arguments: got %d, want %d", len(args), expected))
		}
		return nil
	}
}

func rangeArgsFromUse(min, max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		required, _ := parseUseArgs(cmd.Use)

		if len(args) < min {
			missing := required[len(args):]
			return newMissingArgsError(cmd, missing)
		}
		if len(args) > max {
			return newUsageError(cmd, fmt.Sprintf("too many arguments: got %d, want %d-%d", len(args), min, max))
		}
		return nil
	}
}

func maximumArgsFromUse(max int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > max {
			return newUsageError(cmd, fmt.Sprintf("too many arguments: got %d, want at most %d", len(args), max))
		}
		return nil
	}
}

func newMissingArgsError(cmd *cobra.Command, missing []string) error {
	if len(missing) == 0 {
		return newUsageError(cmd, "missing required argument")
	}

	label := "argument"
	if len(missing) > 1 {
		label = "arguments"
	}

	return newUsageError(cmd, fmt.Sprintf("missing required %s: %s", label, strings.Join(missing, ", ")))
}

func newUsageError(cmd *cobra.Command, message string) error {
	return &usageError{
		message: fmt.Sprintf("%s\nUsage: %s", message, usageLine(cmd)),
	}
}

func usageLine(cmd *cobra.Command) string {
	parts := strings.Fields(cmd.Use)
	if len(parts) <= 1 {
		return cmd.CommandPath()
	}
	return strings.TrimSpace(cmd.CommandPath() + " " + strings.Join(parts[1:], " "))
}

func parseUseArgs(use string) (required []string, optional []string) {
	parts := strings.Fields(use)
	if len(parts) <= 1 {
		return nil, nil
	}

	for _, part := range parts[1:] {
		switch {
		case strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">"):
			required = append(required, strings.Trim(part, "<>"))
		case strings.HasPrefix(part, "[") && strings.HasSuffix(part, "]"):
			optional = append(optional, strings.Trim(part, "[]"))
		}
	}

	return required, optional
}
