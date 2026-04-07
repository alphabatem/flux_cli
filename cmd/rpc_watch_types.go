package cmd

type watchPingOutput struct {
	CreatedAt string `json:"createdAt,omitempty"`
}

type watchPongOutput struct {
	CreatedAt string `json:"createdAt,omitempty"`
	ID        int32  `json:"id"`
}

type watchSignatureConfirmOutput struct {
	Signature string      `json:"signature"`
	Confirmed bool        `json:"confirmed"`
	Update    interface{} `json:"update,omitempty"`
}

type watchAccountUpdateOutput struct {
	CreatedAt    string `json:"createdAt,omitempty"`
	Slot         uint64 `json:"slot"`
	IsStartup    bool   `json:"isStartup"`
	Pubkey       string `json:"pubkey"`
	Lamports     uint64 `json:"lamports"`
	Owner        string `json:"owner"`
	Executable   bool   `json:"executable"`
	RentEpoch    uint64 `json:"rentEpoch"`
	Data         string `json:"data"`
	WriteVersion uint64 `json:"writeVersion"`
	TxnSignature string `json:"txnSignature,omitempty"`
}

type watchAccountInfoOutput struct {
	Pubkey       string `json:"pubkey"`
	Lamports     uint64 `json:"lamports"`
	Owner        string `json:"owner"`
	Executable   bool   `json:"executable"`
	RentEpoch    uint64 `json:"rentEpoch"`
	Data         string `json:"data"`
	WriteVersion uint64 `json:"writeVersion"`
	TxnSignature string `json:"txnSignature,omitempty"`
}

type watchSlotUpdateOutput struct {
	CreatedAt string  `json:"createdAt,omitempty"`
	Slot      uint64  `json:"slot"`
	Parent    *uint64 `json:"parent,omitempty"`
	Status    string  `json:"status"`
	DeadError string  `json:"deadError,omitempty"`
}

type watchTransactionUpdateOutput struct {
	CreatedAt   string                      `json:"createdAt,omitempty"`
	Slot        uint64                      `json:"slot"`
	Signature   string                      `json:"signature"`
	IsVote      bool                        `json:"isVote"`
	Transaction *watchTransactionBodyOutput `json:"transaction,omitempty"`
	Index       uint64                      `json:"index"`
}

type watchTransactionBodyOutput struct {
	Signatures []string            `json:"signatures,omitempty"`
	Message    *watchMessageOutput `json:"message,omitempty"`
}

type watchMessageOutput struct {
	Header              *watchMessageHeaderOutput        `json:"header,omitempty"`
	AccountKeys         []string                         `json:"accountKeys,omitempty"`
	RecentBlockhash     string                           `json:"recentBlockhash"`
	Instructions        []watchCompiledInstructionOutput `json:"instructions,omitempty"`
	Versioned           bool                             `json:"versioned"`
	AddressTableLookups []watchAddressTableLookupOutput  `json:"addressTableLookups,omitempty"`
}

type watchMessageHeaderOutput struct {
	NumRequiredSignatures       uint32 `json:"numRequiredSignatures"`
	NumReadonlySignedAccounts   uint32 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts uint32 `json:"numReadonlyUnsignedAccounts"`
}

type watchCompiledInstructionOutput struct {
	ProgramIDIndex uint32 `json:"programIdIndex"`
	Accounts       string `json:"accounts"`
	Data           string `json:"data"`
}

type watchAddressTableLookupOutput struct {
	AccountKey      string `json:"accountKey"`
	WritableIndexes string `json:"writableIndexes"`
	ReadonlyIndexes string `json:"readonlyIndexes"`
}

type watchTransactionErrorOutput struct {
	Err string `json:"err"`
}

type watchTransactionStatusOutput struct {
	CreatedAt string                       `json:"createdAt,omitempty"`
	Slot      uint64                       `json:"slot"`
	Signature string                       `json:"signature"`
	IsVote    bool                         `json:"isVote"`
	Index     uint64                       `json:"index"`
	Err       *watchTransactionErrorOutput `json:"err,omitempty"`
}

type watchRewardOutput struct {
	Pubkey      string `json:"pubkey"`
	Lamports    int64  `json:"lamports"`
	PostBalance uint64 `json:"postBalance"`
	RewardType  string `json:"rewardType"`
	Commission  string `json:"commission"`
}

type watchRewardsOutput struct {
	Rewards       []watchRewardOutput       `json:"rewards,omitempty"`
	NumPartitions *watchNumPartitionsOutput `json:"numPartitions,omitempty"`
}

type watchNumPartitionsOutput struct {
	NumPartitions uint64 `json:"numPartitions"`
}

type watchUnixTimestampOutput struct {
	Timestamp int64 `json:"timestamp"`
}

type watchBlockHeightOutput struct {
	BlockHeight uint64 `json:"blockHeight"`
}

type watchBlockUpdateOutput struct {
	CreatedAt                string                         `json:"createdAt,omitempty"`
	Slot                     uint64                         `json:"slot"`
	Blockhash                string                         `json:"blockhash"`
	Rewards                  *watchRewardsOutput            `json:"rewards,omitempty"`
	BlockTime                *watchUnixTimestampOutput      `json:"blockTime,omitempty"`
	BlockHeight              *watchBlockHeightOutput        `json:"blockHeight,omitempty"`
	ParentSlot               uint64                         `json:"parentSlot"`
	ParentBlockhash          string                         `json:"parentBlockhash"`
	ExecutedTransactionCount uint64                         `json:"executedTransactionCount"`
	Transactions             []watchTransactionUpdateOutput `json:"transactions,omitempty"`
	UpdatedAccountCount      uint64                         `json:"updatedAccountCount"`
	Accounts                 []watchAccountInfoOutput       `json:"accounts,omitempty"`
	EntriesCount             uint64                         `json:"entriesCount"`
	Entries                  []watchEntryOutput             `json:"entries,omitempty"`
}

type watchBlockMetaOutput struct {
	CreatedAt                string                    `json:"createdAt,omitempty"`
	Slot                     uint64                    `json:"slot"`
	Blockhash                string                    `json:"blockhash"`
	Rewards                  *watchRewardsOutput       `json:"rewards,omitempty"`
	BlockTime                *watchUnixTimestampOutput `json:"blockTime,omitempty"`
	BlockHeight              *watchBlockHeightOutput   `json:"blockHeight,omitempty"`
	ParentSlot               uint64                    `json:"parentSlot"`
	ParentBlockhash          string                    `json:"parentBlockhash"`
	ExecutedTransactionCount uint64                    `json:"executedTransactionCount"`
	EntriesCount             uint64                    `json:"entriesCount"`
}

type watchEntryOutput struct {
	CreatedAt                string `json:"createdAt,omitempty"`
	Slot                     uint64 `json:"slot"`
	Index                    uint64 `json:"index"`
	NumHashes                uint64 `json:"numHashes"`
	Hash                     string `json:"hash"`
	ExecutedTransactionCount uint64 `json:"executedTransactionCount"`
	StartingTransactionIndex uint64 `json:"startingTransactionIndex"`
}
