package cmd

import (
	"testing"
	"time"
)

func TestStreamCommandsHaveTimeoutFlag(t *testing.T) {
	cases := []struct {
		name       string
		cmdDefault string
	}{
		{"account.watch", rpcAccountWatchCmd.Flag("timeout").DefValue},
		{"slot.watch", rpcSlotWatchCmd.Flag("timeout").DefValue},
		{"transaction.watch", rpcTxWatchCmd.Flag("timeout").DefValue},
		{"signature.watch", rpcSignatureWatchCmd.Flag("timeout").DefValue},
	}

	for _, tc := range cases {
		if tc.cmdDefault != "0s" {
			t.Fatalf("%s timeout default=%q want=0s", tc.name, tc.cmdDefault)
		}
	}
}

func TestSignatureConfirmTimeoutDefault(t *testing.T) {
	flag := rpcSignatureConfirmCmd.Flag("timeout")
	if flag == nil {
		t.Fatal("signature confirm timeout flag missing")
	}
	if flag.DefValue != (30 * time.Second).String() {
		t.Fatalf("timeout default=%q want=%q", flag.DefValue, (30 * time.Second).String())
	}
}

func TestDataTokensCandlesCountDefault(t *testing.T) {
	flag := dataTokensCandlesCmd.Flag("count")
	if flag == nil {
		t.Fatal("data tokens candles count flag missing")
	}
	if flag.DefValue != "10" {
		t.Fatalf("count default=%q want=%q", flag.DefValue, "10")
	}
}

func TestDataTokensCandlesIntervalDefault(t *testing.T) {
	flag := dataTokensCandlesCmd.Flag("interval")
	if flag == nil {
		t.Fatal("data tokens candles interval flag missing")
	}
	if flag.DefValue != "1" {
		t.Fatalf("interval default=%q want=%q", flag.DefValue, "1")
	}
}

func TestDataBulkCandlesIntervalDefault(t *testing.T) {
	flag := dataBulkCandlesCmd.Flag("interval")
	if flag == nil {
		t.Fatal("data bulk candles interval flag missing")
	}
	if flag.DefValue != "1" {
		t.Fatalf("interval default=%q want=%q", flag.DefValue, "1")
	}
}
