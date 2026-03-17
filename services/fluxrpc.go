package services

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/alphabatem/common/context"
	"github.com/alphabatem/flux_cli/dto"
	"github.com/alphabatem/flux_cli/internal/httpclient"
)

const FLUXRPC_SVC = "fluxrpc_svc"

type FluxRPCService struct {
	context.DefaultService

	Client *httpclient.Client
}

func (s *FluxRPCService) Id() string {
	return FLUXRPC_SVC
}

// ResolveRPCBaseURL returns the base URL for the configured region.
// Region takes precedence over base_url if set.
func ResolveRPCBaseURL(cfg *dto.FluxRPCConfig) string {
	switch cfg.Region {
	case "eu":
		return "https://eu.fluxrpc.com"
	case "us":
		return "https://us.fluxrpc.com"
	default:
		if cfg.BaseURL != "" {
			return cfg.BaseURL
		}
		return "https://us.fluxrpc.com"
	}
}

func (s *FluxRPCService) Configure(ctx *context.Context) error {
	cfg := ctx.Service(CONFIG_SVC).(*ConfigService).Config()

	s.Client = httpclient.New(
		ResolveRPCBaseURL(&cfg.FluxRPC),
		cfg.FluxRPC.APIKey,
		"X-API-KEY",
	)

	return s.DefaultService.Configure(ctx)
}

// jsonRPCRequest is the standard Solana JSON-RPC request format.
type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

// jsonRPCResponse is the standard Solana JSON-RPC response format.
type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Call executes an arbitrary JSON-RPC method.
func (s *FluxRPCService) Call(method string, params []interface{}) (interface{}, error) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %w", err)
	}

	var resp jsonRPCResponse
	if err := s.Client.Post("", bytes.NewReader(body), &resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	var result interface{}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("decoding result: %w", err)
	}

	return result, nil
}

// Convenience methods for common RPC calls

func (s *FluxRPCService) GetBalance(pubkey string, commitment string) (interface{}, error) {
	params := []interface{}{pubkey}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getBalance", params)
}

func (s *FluxRPCService) GetAccountInfo(pubkey string, encoding string, commitment string) (interface{}, error) {
	params := []interface{}{pubkey}
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getAccountInfo", params)
}

func (s *FluxRPCService) GetMultipleAccounts(pubkeys []string, encoding string, commitment string) (interface{}, error) {
	params := []interface{}{pubkeys}
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getMultipleAccounts", params)
}

func (s *FluxRPCService) GetProgramAccounts(programId string, encoding string, filters string, commitment string) (interface{}, error) {
	params := []interface{}{programId}
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if filters != "" {
		var f []interface{}
		if err := json.Unmarshal([]byte(filters), &f); err == nil {
			opts["filters"] = f
		}
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getProgramAccounts", params)
}

func (s *FluxRPCService) GetMinimumBalanceForRentExemption(dataLength uint64, commitment string) (interface{}, error) {
	params := []interface{}{dataLength}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getMinimumBalanceForRentExemption", params)
}

func (s *FluxRPCService) GetLargestAccounts(filter string, commitment string) (interface{}, error) {
	opts := map[string]interface{}{}
	if filter != "" {
		opts["filter"] = filter
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	var params []interface{}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getLargestAccounts", params)
}

func (s *FluxRPCService) GetBlock(slot uint64, encoding string, txDetail string, commitment string) (interface{}, error) {
	params := []interface{}{slot}
	opts := map[string]interface{}{
		"maxSupportedTransactionVersion": 0,
	}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if txDetail != "" {
		opts["transactionDetails"] = txDetail
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	params = append(params, opts)
	return s.Call("getBlock", params)
}

func (s *FluxRPCService) GetBlocks(startSlot, endSlot uint64, commitment string) (interface{}, error) {
	params := []interface{}{startSlot}
	if endSlot > 0 {
		params = append(params, endSlot)
	}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getBlocks", params)
}

func (s *FluxRPCService) GetBlocksWithLimit(startSlot, limit uint64, commitment string) (interface{}, error) {
	params := []interface{}{startSlot, limit}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getBlocksWithLimit", params)
}

func (s *FluxRPCService) GetBlockHeight(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getBlockHeight", params)
}

func (s *FluxRPCService) GetBlockTime(slot uint64) (interface{}, error) {
	return s.Call("getBlockTime", []interface{}{slot})
}

func (s *FluxRPCService) GetBlockCommitment(slot uint64) (interface{}, error) {
	return s.Call("getBlockCommitment", []interface{}{slot})
}

func (s *FluxRPCService) GetBlockProduction(identity string, commitment string) (interface{}, error) {
	opts := map[string]interface{}{}
	if identity != "" {
		opts["identity"] = identity
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	var params []interface{}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getBlockProduction", params)
}

func (s *FluxRPCService) GetLatestBlockhash(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getLatestBlockhash", params)
}

func (s *FluxRPCService) IsBlockhashValid(blockhash string, commitment string) (interface{}, error) {
	params := []interface{}{blockhash}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("isBlockhashValid", params)
}

func (s *FluxRPCService) GetTransaction(signature string, encoding string, commitment string) (interface{}, error) {
	params := []interface{}{signature}
	opts := map[string]interface{}{
		"maxSupportedTransactionVersion": 0,
	}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	params = append(params, opts)
	return s.Call("getTransaction", params)
}

func (s *FluxRPCService) GetTransactionCount(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getTransactionCount", params)
}

func (s *FluxRPCService) GetSignaturesForAddress(address string, limit int, before, until, commitment string) (interface{}, error) {
	params := []interface{}{address}
	opts := map[string]interface{}{}
	if limit > 0 {
		opts["limit"] = limit
	}
	if before != "" {
		opts["before"] = before
	}
	if until != "" {
		opts["until"] = until
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getSignaturesForAddress", params)
}

func (s *FluxRPCService) GetSignatureStatuses(signatures []string, searchHistory bool) (interface{}, error) {
	params := []interface{}{signatures}
	if searchHistory {
		params = append(params, map[string]bool{"searchTransactionHistory": true})
	}
	return s.Call("getSignatureStatuses", params)
}

func (s *FluxRPCService) GetFeeForMessage(message string, commitment string) (interface{}, error) {
	params := []interface{}{message}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getFeeForMessage", params)
}

func (s *FluxRPCService) GetTokenAccountBalance(pubkey string, commitment string) (interface{}, error) {
	params := []interface{}{pubkey}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getTokenAccountBalance", params)
}

func (s *FluxRPCService) GetTokenAccountsByOwner(owner string, filterKey, filterValue string, encoding string, commitment string) (interface{}, error) {
	params := []interface{}{owner}
	filter := map[string]string{filterKey: filterValue}
	params = append(params, filter)
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getTokenAccountsByOwner", params)
}

func (s *FluxRPCService) GetTokenAccountsByDelegate(delegate string, filterKey, filterValue string, encoding string, commitment string) (interface{}, error) {
	params := []interface{}{delegate}
	filter := map[string]string{filterKey: filterValue}
	params = append(params, filter)
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getTokenAccountsByDelegate", params)
}

func (s *FluxRPCService) GetTokenLargestAccounts(mint string, commitment string) (interface{}, error) {
	params := []interface{}{mint}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getTokenLargestAccounts", params)
}

func (s *FluxRPCService) GetTokenSupply(mint string, commitment string) (interface{}, error) {
	params := []interface{}{mint}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getTokenSupply", params)
}

func (s *FluxRPCService) GetSlot(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getSlot", params)
}

func (s *FluxRPCService) GetSlotLeader(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getSlotLeader", params)
}

func (s *FluxRPCService) GetSlotLeaders(startSlot, limit uint64) (interface{}, error) {
	return s.Call("getSlotLeaders", []interface{}{startSlot, limit})
}

func (s *FluxRPCService) MinimumLedgerSlot() (interface{}, error) {
	return s.Call("minimumLedgerSlot", nil)
}

func (s *FluxRPCService) GetMaxRetransmitSlot() (interface{}, error) {
	return s.Call("getMaxRetransmitSlot", nil)
}

func (s *FluxRPCService) GetMaxShredInsertSlot() (interface{}, error) {
	return s.Call("getMaxShredInsertSlot", nil)
}

func (s *FluxRPCService) GetHighestSnapshotSlot() (interface{}, error) {
	return s.Call("getHighestSnapshotSlot", nil)
}

func (s *FluxRPCService) GetEpochInfo(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getEpochInfo", params)
}

func (s *FluxRPCService) GetEpochSchedule() (interface{}, error) {
	return s.Call("getEpochSchedule", nil)
}

func (s *FluxRPCService) GetLeaderSchedule(slot uint64, identity string, commitment string) (interface{}, error) {
	var params []interface{}
	if slot > 0 {
		params = append(params, slot)
	} else {
		params = append(params, nil)
	}
	opts := map[string]interface{}{}
	if identity != "" {
		opts["identity"] = identity
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getLeaderSchedule", params)
}

func (s *FluxRPCService) GetHealth() (interface{}, error) {
	return s.Call("getHealth", nil)
}

func (s *FluxRPCService) GetVersion() (interface{}, error) {
	return s.Call("getVersion", nil)
}

func (s *FluxRPCService) GetIdentity() (interface{}, error) {
	return s.Call("getIdentity", nil)
}

func (s *FluxRPCService) GetClusterNodes() (interface{}, error) {
	return s.Call("getClusterNodes", nil)
}

func (s *FluxRPCService) GetGenesisHash() (interface{}, error) {
	return s.Call("getGenesisHash", nil)
}

func (s *FluxRPCService) GetFirstAvailableBlock() (interface{}, error) {
	return s.Call("getFirstAvailableBlock", nil)
}

func (s *FluxRPCService) GetRecentPerformanceSamples(limit int) (interface{}, error) {
	var params []interface{}
	if limit > 0 {
		params = append(params, limit)
	}
	return s.Call("getRecentPerformanceSamples", params)
}

func (s *FluxRPCService) GetRecentPrioritizationFees(addresses []string) (interface{}, error) {
	var params []interface{}
	if len(addresses) > 0 {
		params = append(params, addresses)
	}
	return s.Call("getRecentPrioritizationFees", params)
}

func (s *FluxRPCService) GetSupply(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getSupply", params)
}

func (s *FluxRPCService) GetVoteAccounts(votePubkey string, commitment string) (interface{}, error) {
	opts := map[string]interface{}{}
	if votePubkey != "" {
		opts["votePubkey"] = votePubkey
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	var params []interface{}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getVoteAccounts", params)
}

func (s *FluxRPCService) GetStakeMinimumDelegation(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getStakeMinimumDelegation", params)
}

func (s *FluxRPCService) GetInflationGovernor(commitment string) (interface{}, error) {
	var params []interface{}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("getInflationGovernor", params)
}

func (s *FluxRPCService) GetInflationRate() (interface{}, error) {
	return s.Call("getInflationRate", nil)
}

func (s *FluxRPCService) GetInflationReward(addresses []string, epoch uint64, commitment string) (interface{}, error) {
	params := []interface{}{addresses}
	opts := map[string]interface{}{}
	if epoch > 0 {
		opts["epoch"] = epoch
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("getInflationReward", params)
}

func (s *FluxRPCService) SendTransaction(tx string, encoding string, skipPreflight bool, commitment string) (interface{}, error) {
	params := []interface{}{tx}
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if skipPreflight {
		opts["skipPreflight"] = true
	}
	if commitment != "" {
		opts["preflightCommitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("sendTransaction", params)
}

func (s *FluxRPCService) SimulateTransaction(tx string, encoding string, sigVerify bool, commitment string) (interface{}, error) {
	params := []interface{}{tx}
	opts := map[string]interface{}{}
	if encoding != "" {
		opts["encoding"] = encoding
	}
	if sigVerify {
		opts["sigVerify"] = true
	}
	if commitment != "" {
		opts["commitment"] = commitment
	}
	if len(opts) > 0 {
		params = append(params, opts)
	}
	return s.Call("simulateTransaction", params)
}

func (s *FluxRPCService) RequestAirdrop(pubkey string, lamports uint64, commitment string) (interface{}, error) {
	params := []interface{}{pubkey, lamports}
	if commitment != "" {
		params = append(params, map[string]string{"commitment": commitment})
	}
	return s.Call("requestAirdrop", params)
}
