package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
)

// rpcRequest captures an incoming JSON-RPC request for assertions.
type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// newTestRPCService creates a FluxRPCService backed by an httptest server.
// The handler captures requests and returns the given result.
func newTestRPCService(t *testing.T, result interface{}, captured *rpcRequest) (*FluxRPCService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if captured != nil {
			json.Unmarshal(body, captured)
		}

		resultBytes, _ := json.Marshal(result)
		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  json.RawMessage(resultBytes),
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(resp)
	}))

	svc := &FluxRPCService{
		Client: httpclient.New(srv.URL, "test-key", "Authorization"),
	}
	return svc, srv
}

func newTestRPCServiceError(t *testing.T) (*FluxRPCService, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"error": map[string]interface{}{
				"code":    -32600,
				"message": "Invalid Request",
			},
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(resp)
	}))

	svc := &FluxRPCService{
		Client: httpclient.New(srv.URL, "test-key", "Authorization"),
	}
	return svc, srv
}

func TestFluxRPC_Call(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, "ok", &captured)
	defer srv.Close()

	result, err := svc.Call("getHealth", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got %v", result)
	}
	if captured.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %s", captured.JSONRPC)
	}
	if captured.Method != "getHealth" {
		t.Errorf("expected method getHealth, got %s", captured.Method)
	}
}

func TestFluxRPC_Call_RPCError(t *testing.T) {
	svc, srv := newTestRPCServiceError(t)
	defer srv.Close()

	_, err := svc.Call("badMethod", nil)
	if err == nil {
		t.Fatal("expected error for RPC error response")
	}
	if err.Error() != "RPC error -32600: Invalid Request" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFluxRPC_GetBalance(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, map[string]interface{}{"value": 1000000}, &captured)
	defer srv.Close()

	_, err := svc.GetBalance("pubkey123", "confirmed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if captured.Method != "getBalance" {
		t.Errorf("expected getBalance, got %s", captured.Method)
	}
	if len(captured.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(captured.Params))
	}
	if captured.Params[0] != "pubkey123" {
		t.Errorf("expected pubkey123, got %v", captured.Params[0])
	}
	opts, ok := captured.Params[1].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map for opts, got %T", captured.Params[1])
	}
	if opts["commitment"] != "confirmed" {
		t.Errorf("expected confirmed, got %v", opts["commitment"])
	}
}

func TestFluxRPC_GetBalance_NoCommitment(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, map[string]interface{}{"value": 500}, &captured)
	defer srv.Close()

	_, err := svc.GetBalance("pubkey456", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(captured.Params) != 1 {
		t.Errorf("expected 1 param when no commitment, got %d", len(captured.Params))
	}
}

func TestFluxRPC_GetAccountInfo(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetAccountInfo("abc", "jsonParsed", "finalized")
	if captured.Method != "getAccountInfo" {
		t.Errorf("expected getAccountInfo, got %s", captured.Method)
	}
	opts, ok := captured.Params[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected options map")
	}
	if opts["encoding"] != "jsonParsed" {
		t.Errorf("expected jsonParsed, got %v", opts["encoding"])
	}
	if opts["commitment"] != "finalized" {
		t.Errorf("expected finalized, got %v", opts["commitment"])
	}
}

func TestFluxRPC_GetAccountInfo_NoOptions(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetAccountInfo("abc", "", "")
	if len(captured.Params) != 1 {
		t.Errorf("expected 1 param with no options, got %d", len(captured.Params))
	}
}

func TestFluxRPC_GetBlock(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetBlock(12345, "jsonParsed", "signatures", "confirmed")
	if captured.Method != "getBlock" {
		t.Errorf("expected getBlock, got %s", captured.Method)
	}
	// First param is slot
	slot, ok := captured.Params[0].(float64) // JSON numbers are float64
	if !ok {
		t.Fatalf("expected float64 for slot, got %T", captured.Params[0])
	}
	if slot != 12345 {
		t.Errorf("expected slot 12345, got %v", slot)
	}
	opts, ok := captured.Params[1].(map[string]interface{})
	if !ok {
		t.Fatal("expected options map")
	}
	if opts["encoding"] != "jsonParsed" {
		t.Errorf("expected jsonParsed, got %v", opts["encoding"])
	}
	if opts["transactionDetails"] != "signatures" {
		t.Errorf("expected signatures, got %v", opts["transactionDetails"])
	}
	if opts["maxSupportedTransactionVersion"] != float64(0) {
		t.Errorf("expected maxSupportedTransactionVersion 0, got %v", opts["maxSupportedTransactionVersion"])
	}
}

func TestFluxRPC_GetTransaction(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetTransaction("sig123", "base64", "")
	if captured.Method != "getTransaction" {
		t.Errorf("expected getTransaction, got %s", captured.Method)
	}
	if captured.Params[0] != "sig123" {
		t.Errorf("expected sig123, got %v", captured.Params[0])
	}
	opts := captured.Params[1].(map[string]interface{})
	if opts["encoding"] != "base64" {
		t.Errorf("expected base64, got %v", opts["encoding"])
	}
	if _, exists := opts["commitment"]; exists {
		t.Error("commitment should not be set when empty")
	}
}

func TestFluxRPC_GetSignaturesForAddress(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetSignaturesForAddress("addr1", 10, "before-sig", "", "confirmed")
	opts := captured.Params[1].(map[string]interface{})
	if opts["limit"] != float64(10) {
		t.Errorf("expected limit 10, got %v", opts["limit"])
	}
	if opts["before"] != "before-sig" {
		t.Errorf("expected before-sig, got %v", opts["before"])
	}
	if _, exists := opts["until"]; exists {
		t.Error("until should not be set when empty")
	}
}

func TestFluxRPC_GetTokenAccountsByOwner(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetTokenAccountsByOwner("owner1", "mint", "mintAddr", "jsonParsed", "")
	if captured.Method != "getTokenAccountsByOwner" {
		t.Errorf("expected getTokenAccountsByOwner, got %s", captured.Method)
	}
	filter := captured.Params[1].(map[string]interface{})
	if filter["mint"] != "mintAddr" {
		t.Errorf("expected mint filter, got %v", filter)
	}
}

func TestFluxRPC_SendTransaction(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, "txn-signature", &captured)
	defer srv.Close()

	result, err := svc.SendTransaction("base64txdata", "base64", true, "confirmed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "txn-signature" {
		t.Errorf("expected txn-signature, got %v", result)
	}
	if captured.Method != "sendTransaction" {
		t.Errorf("expected sendTransaction, got %s", captured.Method)
	}
	opts := captured.Params[1].(map[string]interface{})
	if opts["skipPreflight"] != true {
		t.Error("expected skipPreflight=true")
	}
	if opts["encoding"] != "base64" {
		t.Errorf("expected base64 encoding, got %v", opts["encoding"])
	}
}

func TestFluxRPC_NoParamMethods(t *testing.T) {
	methods := []struct {
		name   string
		call   func(svc *FluxRPCService) (interface{}, error)
		expect string
	}{
		{"getHealth", func(s *FluxRPCService) (interface{}, error) { return s.GetHealth() }, "getHealth"},
		{"getVersion", func(s *FluxRPCService) (interface{}, error) { return s.GetVersion() }, "getVersion"},
		{"getIdentity", func(s *FluxRPCService) (interface{}, error) { return s.GetIdentity() }, "getIdentity"},
		{"getClusterNodes", func(s *FluxRPCService) (interface{}, error) { return s.GetClusterNodes() }, "getClusterNodes"},
		{"getGenesisHash", func(s *FluxRPCService) (interface{}, error) { return s.GetGenesisHash() }, "getGenesisHash"},
		{"getFirstAvailableBlock", func(s *FluxRPCService) (interface{}, error) { return s.GetFirstAvailableBlock() }, "getFirstAvailableBlock"},
		{"getEpochSchedule", func(s *FluxRPCService) (interface{}, error) { return s.GetEpochSchedule() }, "getEpochSchedule"},
		{"getInflationRate", func(s *FluxRPCService) (interface{}, error) { return s.GetInflationRate() }, "getInflationRate"},
		{"minimumLedgerSlot", func(s *FluxRPCService) (interface{}, error) { return s.MinimumLedgerSlot() }, "minimumLedgerSlot"},
		{"getMaxRetransmitSlot", func(s *FluxRPCService) (interface{}, error) { return s.GetMaxRetransmitSlot() }, "getMaxRetransmitSlot"},
		{"getMaxShredInsertSlot", func(s *FluxRPCService) (interface{}, error) { return s.GetMaxShredInsertSlot() }, "getMaxShredInsertSlot"},
		{"getHighestSnapshotSlot", func(s *FluxRPCService) (interface{}, error) { return s.GetHighestSnapshotSlot() }, "getHighestSnapshotSlot"},
	}

	for _, tt := range methods {
		t.Run(tt.name, func(t *testing.T) {
			var captured rpcRequest
			svc, srv := newTestRPCService(t, "ok", &captured)
			defer srv.Close()

			_, err := tt.call(svc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if captured.Method != tt.expect {
				t.Errorf("expected method %s, got %s", tt.expect, captured.Method)
			}
		})
	}
}

func TestFluxRPC_GetProgramAccounts_WithFilters(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetProgramAccounts("prog1", "base64", `[{"dataSize": 165}]`, "")
	opts := captured.Params[1].(map[string]interface{})
	filters, ok := opts["filters"].([]interface{})
	if !ok {
		t.Fatal("expected filters array")
	}
	if len(filters) != 1 {
		t.Errorf("expected 1 filter, got %d", len(filters))
	}
}

func TestFluxRPC_GetLeaderSchedule_NilSlot(t *testing.T) {
	var captured rpcRequest
	svc, srv := newTestRPCService(t, nil, &captured)
	defer srv.Close()

	svc.GetLeaderSchedule(0, "identity1", "")
	// slot=0 means nil
	if captured.Params[0] != nil {
		t.Errorf("expected nil for slot 0, got %v", captured.Params[0])
	}
	opts := captured.Params[1].(map[string]interface{})
	if opts["identity"] != "identity1" {
		t.Errorf("expected identity1, got %v", opts["identity"])
	}
}

func TestResolveRPCBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		cfg      dto.FluxRPCConfig
		expected string
	}{
		{"eu region", dto.FluxRPCConfig{Region: "eu"}, "https://eu.fluxrpc.com"},
		{"us region", dto.FluxRPCConfig{Region: "us"}, "https://us.fluxrpc.com"},
		{"custom base_url", dto.FluxRPCConfig{BaseURL: "https://custom.rpc.com"}, "https://custom.rpc.com"},
		{"region overrides base_url", dto.FluxRPCConfig{Region: "eu", BaseURL: "https://custom.rpc.com"}, "https://eu.fluxrpc.com"},
		{"empty defaults to us", dto.FluxRPCConfig{}, "https://us.fluxrpc.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveRPCBaseURL(&tt.cfg)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}
