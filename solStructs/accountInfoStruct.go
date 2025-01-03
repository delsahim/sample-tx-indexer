package solstructs

// RPC response structures
type AccountInfoResponse struct {
    JsonRPC string         `json:"jsonrpc"`
    Result  AccountResult  `json:"result"`
    ID      int           `json:"id"`
}

type AccountResult struct {
    Context Context     `json:"context"`
    Value   AccountInfo `json:"value"`
}

type Context struct {
    Slot uint64 `json:"slot"`
}

type AccountInfo struct {
    Data       TokenAccountData `json:"data"`
    Executable bool            `json:"executable"`
    Lamports   uint64          `json:"lamports"`
    Owner      string          `json:"owner"`
    RentEpoch  uint64          `json:"rentEpoch"`
}

type TokenAccountData struct {
    Program string        `json:"program"`
    Parsed  ParsedAccount `json:"parsed"`
}

type ParsedAccount struct {
    Type string           `json:"type"`
    Info TokenAccountInfo `json:"info"`
}

type TokenAccountInfo struct {
    Mint        string      `json:"mint"`
    Owner       string      `json:"owner"`
    TokenAmount TokenAmount `json:"tokenAmount"`
}

type TokenAmount struct {
    Amount   string  `json:"amount"`
    Decimals uint8   `json:"decimals"`
    UiAmount float64 `json:"uiAmount"`
}

