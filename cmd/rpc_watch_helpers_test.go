package cmd

import (
	"context"
	"testing"
	"time"

	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
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

func TestProtoUpdateToDataConvertsTransactionIdentifiersToBase58(t *testing.T) {
	update := &pb.SubscribeUpdate{
		Filters: []string{"transactions"},
		UpdateOneof: &pb.SubscribeUpdate_Transaction{
			Transaction: &pb.SubscribeUpdateTransaction{
				Slot: 12,
				Transaction: &pb.SubscribeUpdateTransactionInfo{
					Signature: make([]byte, 64),
					Transaction: &pb.Transaction{
						Signatures: [][]byte{make([]byte, 64)},
						Message: &pb.Message{
							AccountKeys:     [][]byte{make([]byte, 32)},
							RecentBlockhash: make([]byte, 32),
							AddressTableLookups: []*pb.MessageAddressTableLookup{
								{AccountKey: make([]byte, 32)},
							},
						},
					},
				},
			},
		},
	}

	data, err := protoUpdateToData(update)
	if err != nil {
		t.Fatalf("protoUpdateToData() error = %v", err)
	}

	root := data.(*watchUpdateOutput)
	tx := root.Transaction
	txInfo := tx.Transaction
	message := txInfo.Transaction.Message

	if got := txInfo.Signature; got != repeatOnes(64) {
		t.Fatalf("signature=%v want=%q", got, repeatOnes(64))
	}

	signatures := txInfo.Transaction.Signatures
	if got := signatures[0]; got != repeatOnes(64) {
		t.Fatalf("signatures[0]=%v want=%q", got, repeatOnes(64))
	}

	if got := message.RecentBlockhash; got != repeatOnes(32) {
		t.Fatalf("recentBlockhash=%v want=%q", got, repeatOnes(32))
	}

	accountKeys := message.AccountKeys
	if got := accountKeys[0]; got != repeatOnes(32) {
		t.Fatalf("accountKeys[0]=%v want=%q", got, repeatOnes(32))
	}

	lookups := message.AddressTableLookups
	if got := lookups[0].AccountKey; got != repeatOnes(32) {
		t.Fatalf("accountKey=%v want=%q", got, repeatOnes(32))
	}
}

func TestProtoUpdateToDataConvertsAccountIdentifiersToBase58(t *testing.T) {
	update := &pb.SubscribeUpdate{
		Filters: []string{"accounts"},
		UpdateOneof: &pb.SubscribeUpdate_Account{
			Account: &pb.SubscribeUpdateAccount{
				Slot: 7,
				Account: &pb.SubscribeUpdateAccountInfo{
					Pubkey:       make([]byte, 32),
					Owner:        make([]byte, 32),
					Data:         []byte{1, 2, 3},
					TxnSignature: make([]byte, 64),
				},
			},
		},
	}

	data, err := protoUpdateToData(update)
	if err != nil {
		t.Fatalf("protoUpdateToData() error = %v", err)
	}

	root := data.(*watchUpdateOutput)
	account := root.Account.Account

	if got := account.Pubkey; got != repeatOnes(32) {
		t.Fatalf("pubkey=%v want=%q", got, repeatOnes(32))
	}
	if got := account.Owner; got != repeatOnes(32) {
		t.Fatalf("owner=%v want=%q", got, repeatOnes(32))
	}
	if got := account.TxnSignature; got != repeatOnes(64) {
		t.Fatalf("txnSignature=%v want=%q", got, repeatOnes(64))
	}
	if got := account.Data; got != "AQID" {
		t.Fatalf("data=%v want=%q", got, "AQID")
	}
}

func repeatOnes(n int) string {
	out := make([]byte, n)
	for i := range out {
		out[i] = '1'
	}
	return string(out)
}
