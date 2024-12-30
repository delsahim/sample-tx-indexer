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