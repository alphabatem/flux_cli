package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestParseCommitment(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{"", false},
		{"processed", false},
		{"confirmed", false},
		{"finalized", false},
		{"bad", true},
	}

	for _, tt := range tests {
		_, err := parseCommitment(tt.in)
		if (err != nil) != tt.wantErr {
			t.Fatalf("parseCommitment(%q) err=%v wantErr=%v", tt.in, err, tt.wantErr)
		}
	}
}

func TestParseCSV(t *testing.T) {
	got := parseCSV(" a, ,b,, c ")
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("len(parseCSV)=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("parseCSV[%d]=%q want=%q", i, got[i], want[i])
		}
	}
}

func TestStreamContextFromTimeoutFlag_NoTimeout(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	cmd.Flags().Duration("timeout", 0, "")

	ctx, cancel, err := streamContextFromTimeoutFlag(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()

	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline when timeout is 0")
	}
}

func TestStreamContextFromTimeoutFlag_WithTimeout(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	cmd.Flags().Duration("timeout", 0, "")
	if err := cmd.Flags().Set("timeout", "50ms"); err != nil {
		t.Fatalf("set timeout: %v", err)
	}

	ctx, cancel, err := streamContextFromTimeoutFlag(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cancel()

	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected deadline when timeout is set")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("expected deadline in the future")
	}
}
