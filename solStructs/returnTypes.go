package solstructs


type SwapHandlerResponse struct {
	TokenMintA string
	TokenMintB string
	TokenAmountA float64
	TokenAmountB float64
	TransactionHash []string
	AtoB bool
	SignerWallet string
}

type OrcaDecodedInstruction struct {
	Amount                uint64
    OtherAmountThreshold uint64
    SqrtPriceLimit       Uint128 // You'll need to handle uint128
    AmountSpecifiedIsInput bool
    AToB                  bool
}

type TokenTransferDecodedInstruction struct {
    Amount uint64
    Decimals uint8
}

type Uint128 struct {
    Low  uint64
    High uint64
}

type TokenTransferDetails struct {
	Amount float64
	Mint string
}
