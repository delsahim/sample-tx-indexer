package handlers

import (
	"encoding/binary"
	solstructs "indexer_golang/solStructs"

	"github.com/mr-tron/base58"
)



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

func DecodeOrcaSwapData(data string) (solstructs.OrcaDecodedInstruction,error) {
	// decode the data 
	decodedData, err := base58.Decode(data)
	if err != nil {
		return solstructs.OrcaDecodedInstruction{}, nil
	}

	decodedData  = decodedData[8:]
	args := solstructs.OrcaDecodedInstruction{}

	args.Amount = binary.LittleEndian.Uint64(decodedData[:8])
	decodedData  = decodedData[8:]

	args.OtherAmountThreshold = binary.LittleEndian.Uint64(decodedData[:8])
	decodedData  = decodedData[8:]

	args.SqrtPriceLimit = readUint128(decodedData[:16])
    decodedData = decodedData[16:]


	if len(decodedData) >= 2 {
		args.AmountSpecifiedIsInput = decodedData[0] != 0
    	args.AToB = decodedData[1] != 0
	}

	return args, nil
}

func readUint128(data []byte) solstructs.Uint128 {
    return solstructs.Uint128{
        Low:  binary.LittleEndian.Uint64(data[:8]),
        High: binary.LittleEndian.Uint64(data[8:16]),
    }
}




