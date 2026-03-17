package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alphabatem/flux_cli/dto"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

func parseCommitment(raw string) (pb.CommitmentLevel, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "processed":
		return pb.CommitmentLevel_PROCESSED, nil
	case "confirmed":
		return pb.CommitmentLevel_CONFIRMED, nil
	case "finalized":
		return pb.CommitmentLevel_FINALIZED, nil
	default:
		return pb.CommitmentLevel_PROCESSED, fmt.Errorf("invalid commitment %q (expected processed, confirmed, or finalized)", raw)
	}
}

func parseCSV(arg string) []string {
	parts := strings.Split(arg, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func printWatchUpdate(cmd *cobra.Command, endpoint string, update *pb.SubscribeUpdate) error {
	data, err := protoUpdateToData(update)
	if err != nil {
		return err
	}

	output.PrintSuccess(cmd, data, &dto.CLIMeta{
		Service:  "yellowstone",
		Endpoint: endpoint,
	})
	return nil
}

func protoUpdateToData(update *pb.SubscribeUpdate) (interface{}, error) {
	raw, err := protojson.Marshal(update)
	if err != nil {
		return nil, fmt.Errorf("marshal update: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("decode update json: %w", err)
	}
	return data, nil
}

func failUsage(cmd *cobra.Command, message string) {
	output.PrintError(cmd, "INVALID_ARGS", message, nil)
	os.Exit(dto.ExitUsageError)
}

func streamContextFromTimeoutFlag(cmd *cobra.Command) (context.Context, context.CancelFunc, error) {
	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return nil, nil, err
	}
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
		return ctx, cancel, nil
	}
	return cmd.Context(), func() {}, nil
}
