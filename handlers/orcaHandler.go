package handlers

import solstructs "indexer_golang/solStructs"



func HandleOrcaSwapTransaction(inputTransaction solstructs.CompleteTransactionStruct) () {

}

func GetSwapTokenBalace(TokenBalances []solstructs.TokenBalance,userTokenAccountA, userTokenAccountB, vaultTokenAccountA, vaultTokenAccountB int) (string, string, float64, float64, float64, float64) {
	var (
		tokenMintA string
		tokenMintB string
		userBalanceA float64
		userBalanceB float64
		lpBalanceA  float64
		lpBalanceB  float64
	)

	for _, tokenBalance := range TokenBalances {
		if tokenBalance.AccountIndex == userTokenAccountA {
			tokenMintA = tokenBalance.Mint
			userBalanceA = *tokenBalance.UIAmount.UIAmount
		}

		if tokenBalance.AccountIndex == userTokenAccountB {
			tokenMintB = tokenBalance.Mint
			userBalanceB = *tokenBalance.UIAmount.UIAmount
		}

		if tokenBalance.AccountIndex == vaultTokenAccountA {
			lpBalanceA = *tokenBalance.UIAmount.UIAmount
		}
		if tokenBalance.AccountIndex == vaultTokenAccountB {
			lpBalanceB = *tokenBalance.UIAmount.UIAmount
		}
	}

	return tokenMintA, tokenMintB, userBalanceA, userBalanceB, lpBalanceA, lpBalanceB
}


func GetAtoB(pretokenA, postTokenA float64) (bool) {
	if (pretokenA - postTokenA) > 0 {
		return true
	}
	return false
}




