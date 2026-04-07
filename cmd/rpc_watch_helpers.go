package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alphabatem/flux_cli/dto"
	pb "github.com/alphabatem/flux_cli/internal/yellowstonepb"
	"github.com/alphabatem/flux_cli/output"
	"github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"
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
	if update == nil {
		return nil, fmt.Errorf("nil update")
	}

	data := &watchUpdateOutput{
		Filters: update.GetFilters(),
	}
	if ts := update.GetCreatedAt(); ts != nil {
		data.CreatedAt = ts.AsTime().UTC().Format(time.RFC3339Nano)
	}

	switch typed := update.GetUpdateOneof().(type) {
	case *pb.SubscribeUpdate_Account:
		data.Account = accountUpdateToData(typed.Account)
	case *pb.SubscribeUpdate_Slot:
		data.Slot = slotUpdateToData(typed.Slot)
	case *pb.SubscribeUpdate_Transaction:
		data.Transaction = transactionUpdateToData(typed.Transaction)
	case *pb.SubscribeUpdate_TransactionStatus:
		data.TransactionStatus = transactionStatusUpdateToData(typed.TransactionStatus)
	case *pb.SubscribeUpdate_Block:
		data.Block = blockUpdateToData(typed.Block)
	case *pb.SubscribeUpdate_BlockMeta:
		data.BlockMeta = blockMetaUpdateToData(typed.BlockMeta)
	case *pb.SubscribeUpdate_Entry:
		data.Entry = entryUpdateToData(typed.Entry)
	case *pb.SubscribeUpdate_Ping:
		data.Ping = &watchPingOutput{}
	case *pb.SubscribeUpdate_Pong:
		data.Pong = &watchPongOutput{ID: typed.Pong.GetId()}
	}

	return data, nil
}

func accountUpdateToData(update *pb.SubscribeUpdateAccount) *watchAccountUpdateOutput {
	if update == nil {
		return nil
	}
	return &watchAccountUpdateOutput{
		Account:   accountInfoToData(update.GetAccount()),
		Slot:      update.GetSlot(),
		IsStartup: update.GetIsStartup(),
	}
}

func accountInfoToData(info *pb.SubscribeUpdateAccountInfo) *watchAccountInfoOutput {
	if info == nil {
		return nil
	}
	data := &watchAccountInfoOutput{
		Pubkey:       publicKeyString(info.GetPubkey()),
		Lamports:     info.GetLamports(),
		Owner:        publicKeyString(info.GetOwner()),
		Executable:   info.GetExecutable(),
		RentEpoch:    info.GetRentEpoch(),
		Data:         base64.StdEncoding.EncodeToString(info.GetData()),
		WriteVersion: info.GetWriteVersion(),
	}
	if sig := info.GetTxnSignature(); len(sig) > 0 {
		data.TxnSignature = signatureString(sig)
	}
	return data
}

func slotUpdateToData(update *pb.SubscribeUpdateSlot) *watchSlotUpdateOutput {
	if update == nil {
		return nil
	}
	data := &watchSlotUpdateOutput{
		Slot:   update.GetSlot(),
		Status: update.GetStatus().String(),
	}
	if update.Parent != nil {
		parent := update.GetParent()
		data.Parent = &parent
	}
	if update.DeadError != nil {
		data.DeadError = update.GetDeadError()
	}
	return data
}

func transactionUpdateToData(update *pb.SubscribeUpdateTransaction) *watchTransactionUpdateOutput {
	if update == nil {
		return nil
	}
	return &watchTransactionUpdateOutput{
		Transaction: transactionInfoToData(update.GetTransaction()),
		Slot:        update.GetSlot(),
	}
}

func transactionInfoToData(info *pb.SubscribeUpdateTransactionInfo) *watchTransactionInfoOutput {
	if info == nil {
		return nil
	}
	return &watchTransactionInfoOutput{
		Signature:   signatureString(info.GetSignature()),
		IsVote:      info.GetIsVote(),
		Transaction: transactionToData(info.GetTransaction()),
		Meta:        transactionMetaToData(info.GetMeta()),
		Index:       info.GetIndex(),
	}
}

func transactionToData(tx *pb.Transaction) *watchTransactionBodyOutput {
	if tx == nil {
		return nil
	}
	signatures := make([]string, 0, len(tx.GetSignatures()))
	for _, sig := range tx.GetSignatures() {
		signatures = append(signatures, signatureString(sig))
	}
	return &watchTransactionBodyOutput{
		Signatures: signatures,
		Message:    messageToData(tx.GetMessage()),
	}
}

func messageToData(msg *pb.Message) *watchMessageOutput {
	if msg == nil {
		return nil
	}
	accountKeys := make([]string, 0, len(msg.GetAccountKeys()))
	for _, key := range msg.GetAccountKeys() {
		accountKeys = append(accountKeys, publicKeyString(key))
	}
	instructions := make([]watchCompiledInstructionOutput, 0, len(msg.GetInstructions()))
	for _, ix := range msg.GetInstructions() {
		instructions = append(instructions, compiledInstructionToData(ix))
	}
	lookups := make([]watchAddressTableLookupOutput, 0, len(msg.GetAddressTableLookups()))
	for _, lookup := range msg.GetAddressTableLookups() {
		lookups = append(lookups, addressTableLookupToData(lookup))
	}
	return &watchMessageOutput{
		Header:              messageHeaderToData(msg.GetHeader()),
		AccountKeys:         accountKeys,
		RecentBlockhash:     publicKeyString(msg.GetRecentBlockhash()),
		Instructions:        instructions,
		Versioned:           msg.GetVersioned(),
		AddressTableLookups: lookups,
	}
}

func messageHeaderToData(header *pb.MessageHeader) *watchMessageHeaderOutput {
	if header == nil {
		return nil
	}
	return &watchMessageHeaderOutput{
		NumRequiredSignatures:       header.GetNumRequiredSignatures(),
		NumReadonlySignedAccounts:   header.GetNumReadonlySignedAccounts(),
		NumReadonlyUnsignedAccounts: header.GetNumReadonlyUnsignedAccounts(),
	}
}

func compiledInstructionToData(ix *pb.CompiledInstruction) watchCompiledInstructionOutput {
	if ix == nil {
		return watchCompiledInstructionOutput{}
	}
	return watchCompiledInstructionOutput{
		ProgramIDIndex: ix.GetProgramIdIndex(),
		Accounts:       base64.StdEncoding.EncodeToString(ix.GetAccounts()),
		Data:           base64.StdEncoding.EncodeToString(ix.GetData()),
	}
}

func addressTableLookupToData(lookup *pb.MessageAddressTableLookup) watchAddressTableLookupOutput {
	if lookup == nil {
		return watchAddressTableLookupOutput{}
	}
	return watchAddressTableLookupOutput{
		AccountKey:      publicKeyString(lookup.GetAccountKey()),
		WritableIndexes: base64.StdEncoding.EncodeToString(lookup.GetWritableIndexes()),
		ReadonlyIndexes: base64.StdEncoding.EncodeToString(lookup.GetReadonlyIndexes()),
	}
}

func transactionMetaToData(meta *pb.TransactionStatusMeta) *watchTransactionMetaOutput {
	if meta == nil {
		return nil
	}
	data := &watchTransactionMetaOutput{
		Err:                     transactionErrorToData(meta.GetErr()),
		Fee:                     meta.GetFee(),
		PreBalances:             meta.GetPreBalances(),
		PostBalances:            meta.GetPostBalances(),
		InnerInstructions:       innerInstructionsListToData(meta.GetInnerInstructions()),
		InnerInstructionsNone:   meta.GetInnerInstructionsNone(),
		LogMessages:             meta.GetLogMessages(),
		LogMessagesNone:         meta.GetLogMessagesNone(),
		PreTokenBalances:        tokenBalancesToData(meta.GetPreTokenBalances()),
		PostTokenBalances:       tokenBalancesToData(meta.GetPostTokenBalances()),
		Rewards:                 rewardsListToData(meta.GetRewards()),
		LoadedWritableAddresses: publicKeysToStrings(meta.GetLoadedWritableAddresses()),
		LoadedReadonlyAddresses: publicKeysToStrings(meta.GetLoadedReadonlyAddresses()),
		ReturnData:              returnDataToData(meta.GetReturnData()),
		ReturnDataNone:          meta.GetReturnDataNone(),
	}
	if meta.ComputeUnitsConsumed != nil {
		value := meta.GetComputeUnitsConsumed()
		data.ComputeUnitsConsumed = &value
	}
	if meta.CostUnits != nil {
		value := meta.GetCostUnits()
		data.CostUnits = &value
	}
	return data
}

func transactionStatusUpdateToData(update *pb.SubscribeUpdateTransactionStatus) *watchTransactionStatusOutput {
	if update == nil {
		return nil
	}
	return &watchTransactionStatusOutput{
		Slot:      update.GetSlot(),
		Signature: signatureString(update.GetSignature()),
		IsVote:    update.GetIsVote(),
		Index:     update.GetIndex(),
		Err:       transactionErrorToData(update.GetErr()),
	}
}

func blockUpdateToData(update *pb.SubscribeUpdateBlock) *watchBlockUpdateOutput {
	if update == nil {
		return nil
	}
	transactions := make([]watchTransactionInfoOutput, 0, len(update.GetTransactions()))
	for _, tx := range update.GetTransactions() {
		if item := transactionInfoToData(tx); item != nil {
			transactions = append(transactions, *item)
		}
	}
	accounts := make([]watchAccountInfoOutput, 0, len(update.GetAccounts()))
	for _, account := range update.GetAccounts() {
		if item := accountInfoToData(account); item != nil {
			accounts = append(accounts, *item)
		}
	}
	entries := make([]watchEntryOutput, 0, len(update.GetEntries()))
	for _, entry := range update.GetEntries() {
		if item := entryUpdateToData(entry); item != nil {
			entries = append(entries, *item)
		}
	}
	return &watchBlockUpdateOutput{
		Slot:                     update.GetSlot(),
		Blockhash:                update.GetBlockhash(),
		Rewards:                  rewardsToData(update.GetRewards()),
		BlockTime:                unixTimestampToData(update.GetBlockTime()),
		BlockHeight:              blockHeightToData(update.GetBlockHeight()),
		ParentSlot:               update.GetParentSlot(),
		ParentBlockhash:          update.GetParentBlockhash(),
		ExecutedTransactionCount: update.GetExecutedTransactionCount(),
		Transactions:             transactions,
		UpdatedAccountCount:      update.GetUpdatedAccountCount(),
		Accounts:                 accounts,
		EntriesCount:             update.GetEntriesCount(),
		Entries:                  entries,
	}
}

func blockMetaUpdateToData(update *pb.SubscribeUpdateBlockMeta) *watchBlockMetaOutput {
	if update == nil {
		return nil
	}
	return &watchBlockMetaOutput{
		Slot:                     update.GetSlot(),
		Blockhash:                update.GetBlockhash(),
		Rewards:                  rewardsToData(update.GetRewards()),
		BlockTime:                unixTimestampToData(update.GetBlockTime()),
		BlockHeight:              blockHeightToData(update.GetBlockHeight()),
		ParentSlot:               update.GetParentSlot(),
		ParentBlockhash:          update.GetParentBlockhash(),
		ExecutedTransactionCount: update.GetExecutedTransactionCount(),
		EntriesCount:             update.GetEntriesCount(),
	}
}

func entryUpdateToData(entry *pb.SubscribeUpdateEntry) *watchEntryOutput {
	if entry == nil {
		return nil
	}
	return &watchEntryOutput{
		Slot:                     entry.GetSlot(),
		Index:                    entry.GetIndex(),
		NumHashes:                entry.GetNumHashes(),
		Hash:                     base64.StdEncoding.EncodeToString(entry.GetHash()),
		ExecutedTransactionCount: entry.GetExecutedTransactionCount(),
		StartingTransactionIndex: entry.GetStartingTransactionIndex(),
	}
}

func innerInstructionsListToData(items []*pb.InnerInstructions) []watchInnerInstructionsOutput {
	out := make([]watchInnerInstructionsOutput, 0, len(items))
	for _, item := range items {
		out = append(out, innerInstructionsToData(item))
	}
	return out
}

func innerInstructionsToData(item *pb.InnerInstructions) watchInnerInstructionsOutput {
	if item == nil {
		return watchInnerInstructionsOutput{}
	}
	instructions := make([]watchInnerInstructionOutput, 0, len(item.GetInstructions()))
	for _, ix := range item.GetInstructions() {
		instructions = append(instructions, innerInstructionToData(ix))
	}
	return watchInnerInstructionsOutput{
		Index:        item.GetIndex(),
		Instructions: instructions,
	}
}

func innerInstructionToData(ix *pb.InnerInstruction) watchInnerInstructionOutput {
	if ix == nil {
		return watchInnerInstructionOutput{}
	}
	data := watchInnerInstructionOutput{
		ProgramIDIndex: ix.GetProgramIdIndex(),
		Accounts:       base64.StdEncoding.EncodeToString(ix.GetAccounts()),
		Data:           base64.StdEncoding.EncodeToString(ix.GetData()),
	}
	if ix.StackHeight != nil {
		value := ix.GetStackHeight()
		data.StackHeight = &value
	}
	return data
}

func tokenBalancesToData(items []*pb.TokenBalance) []watchTokenBalanceOutput {
	out := make([]watchTokenBalanceOutput, 0, len(items))
	for _, item := range items {
		out = append(out, tokenBalanceToData(item))
	}
	return out
}

func tokenBalanceToData(item *pb.TokenBalance) watchTokenBalanceOutput {
	if item == nil {
		return watchTokenBalanceOutput{}
	}
	return watchTokenBalanceOutput{
		AccountIndex:  item.GetAccountIndex(),
		Mint:          item.GetMint(),
		UITokenAmount: uiTokenAmountToData(item.GetUiTokenAmount()),
		Owner:         item.GetOwner(),
		ProgramID:     item.GetProgramId(),
	}
}

func uiTokenAmountToData(item *pb.UiTokenAmount) *watchUITokenAmountOutput {
	if item == nil {
		return nil
	}
	return &watchUITokenAmountOutput{
		UIAmount:       item.GetUiAmount(),
		Decimals:       item.GetDecimals(),
		Amount:         item.GetAmount(),
		UIAmountString: item.GetUiAmountString(),
	}
}

func transactionErrorToData(err *pb.TransactionError) *watchTransactionErrorOutput {
	if err == nil {
		return nil
	}
	return &watchTransactionErrorOutput{
		Err: base64.StdEncoding.EncodeToString(err.GetErr()),
	}
}

func rewardsToData(rewards *pb.Rewards) *watchRewardsOutput {
	if rewards == nil {
		return nil
	}
	data := &watchRewardsOutput{
		Rewards: rewardsListToData(rewards.GetRewards()),
	}
	if rewards.GetNumPartitions() != nil {
		data.NumPartitions = &watchNumPartitionsOutput{
			NumPartitions: rewards.GetNumPartitions().GetNumPartitions(),
		}
	}
	return data
}

func rewardsListToData(items []*pb.Reward) []watchRewardOutput {
	out := make([]watchRewardOutput, 0, len(items))
	for _, item := range items {
		out = append(out, rewardToData(item))
	}
	return out
}

func rewardToData(item *pb.Reward) watchRewardOutput {
	if item == nil {
		return watchRewardOutput{}
	}
	return watchRewardOutput{
		Pubkey:      item.GetPubkey(),
		Lamports:    item.GetLamports(),
		PostBalance: item.GetPostBalance(),
		RewardType:  item.GetRewardType().String(),
		Commission:  item.GetCommission(),
	}
}

func returnDataToData(item *pb.ReturnData) *watchReturnDataOutput {
	if item == nil {
		return nil
	}
	return &watchReturnDataOutput{
		ProgramID: publicKeyString(item.GetProgramId()),
		Data:      base64.StdEncoding.EncodeToString(item.GetData()),
	}
}

func unixTimestampToData(ts *pb.UnixTimestamp) *watchUnixTimestampOutput {
	if ts == nil {
		return nil
	}
	return &watchUnixTimestampOutput{
		Timestamp: ts.GetTimestamp(),
	}
}

func blockHeightToData(height *pb.BlockHeight) *watchBlockHeightOutput {
	if height == nil {
		return nil
	}
	return &watchBlockHeightOutput{
		BlockHeight: height.GetBlockHeight(),
	}
}

func publicKeysToStrings(keys [][]byte) []string {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, publicKeyString(key))
	}
	return out
}

func publicKeyString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return solana.PublicKeyFromBytes(data).String()
}

func signatureString(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	return solana.SignatureFromBytes(data).String()
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
