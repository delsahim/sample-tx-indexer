package solstructs

// fix the transaction entity of the solana block struvt
type SolanaBlockSubscribe struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Context struct {
				Slot int `json:"slot"`
			} `json:"context"`
			Value struct {
				Block *struct {
					BlockHeight       int    `json:"blockHeight"`
					BlockTime         int    `json:"blockTime"`
					Blockhash         string `json:"blockhash"`
					ParentSlot        int    `json:"parentSlot"`
					PreviousBlockhash string `json:"previousBlockhash"`
					Transactions      []CompleteTransactionStruct  `json:"transactions"`
				} `json:"block"`
				Err *struct {
					InstructionError []interface{} `json:"InstructionError"`
				} `json:"err"`
				Slot int `json:"slot"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

type TokenBalance struct {
	AccountIndex int    `json:"accountIndex"`
	Mint         string `json:"mint"`
	Owner        string `json:"owner"`
	UIAmount     struct {
		UIAmount    *float64 `json:"uiAmount"`
		Decimals    int      `json:"decimals"`
		Amount      string   `json:"amount"`
		UIAmountStr string   `json:"uiAmountString"`
	} `json:"uiTokenAmount"`
}

type CompleteTransactionStruct struct {
		Transaction TransactionStruct `json:"transaction"`
		Meta      TransactionMeta   `json:"meta"`
		Version *interface{} `json:"version"`
}


type TransactionStruct struct {
	Signatures []string `json:"signatures"`
	Message    struct {
		Header struct {
			NumRequiredSignatures       int `json:"numRequiredSignatures"`
			NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
			NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
		} `json:"header"`
		AccountKeys []string `json:"accountKeys"`
		RecentBlockhash string `json:"recentBlockhash"`
		Instructions []InstructionStruct  `json:"instructions"`
		AddressTableLookups []struct {
			AccountKey      string `json:"accountKey"`
			WritableIndexes []int  `json:"writableIndexes"`
			ReadonlyIndexes []int  `json:"readonlyIndexes"`
		} `json:"addressTableLookups"`
	} `json:"message"`
} 

type TransactionMeta struct {
	Err         interface{} `json:"err"` // Using interface{} as it can be null
	Status      struct {
		Ok interface{} `json:"Ok"` // Using interface{} as it can be null
	} `json:"status"`
	Fee          int `json:"fee"`
	PreBalances  []int `json:"preBalances"`
	PostBalances []int `json:"postBalances"`
	InnerInstructions []struct {
		Index        int `json:"index"`
		Instructions []InstructionStruct `json:"instructions"`
	} `json:"innerInstructions"`
	LogMessages []string `json:"logMessages"`
	PreTokenBalances []TokenBalance `json:"preTokenBalances"`
	PostTokenBalances []TokenBalance `json:"postTokenBalances"`
	Rewards interface{} `json:"rewards"` // Using interface{} as it can be null
	LoadedAddresses struct {
		Writable  []string `json:"writable"`
		Readonly  []string `json:"readonly"`
	} `json:"loadedAddresses"`
	ComputeUnitsConsumed int `json:"computeUnitsConsumed"`
}

type InstructionStruct struct {
	ProgramIdIndex int     `json:"programIdIndex"`
	Accounts      []int    `json:"accounts"`
	Data         string    `json:"data"`
	StackHeight  *int     `json:"stackHeight"` // Using pointer for nullable int
}