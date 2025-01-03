package handlers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	solstructs "indexer_golang/solStructs"
	"math"
	"net/http"

	"github.com/mr-tron/base58"
)


func GetTransferAmount(data string) (solstructs.TokenTransferDecodedInstruction, error) {
	decodedData, err := base58.Decode(data)
    if err != nil {
        return solstructs.TokenTransferDecodedInstruction{}, err
    }

	if len(decodedData) < 9 {
        return solstructs.TokenTransferDecodedInstruction{}, fmt.Errorf("instruction data too short")
    }

	// get the discriminator

	disc := decodedData[0]
	args := solstructs.TokenTransferDecodedInstruction{}
	if disc == 3 {
		// handle normal token transfer
		args.Amount = binary.LittleEndian.Uint64(decodedData[1:9]) 
		
	} else if (disc == 12) {
		// handle the transfer checked info
		if len(decodedData) < 10 {
            return solstructs.TokenTransferDecodedInstruction{}, fmt.Errorf("TransferChecked data too short")
        }
        args.Decimals = uint8(decodedData[9])
	} else {
		
	}

	return args, nil


}

func DecodeSystemTransfer(accountList []string, instruction solstructs.InstructionStruct, rpcUrl string) (solstructs.TokenTransferDetails, error) {
	// gget the source acct
	sourceAcct := accountList[instruction.Accounts[0]]

	// get the amount
	decodedData, err := GetTransferAmount(instruction.Data)
	if err != nil {
		// return based on error
		return solstructs.TokenTransferDetails{}, err
	}
	transferAmount := float64(decodedData.Amount)

	//get the mint and decimals
	mint, decimals, err := GetTokenAccountInfo(rpcUrl, sourceAcct)
	if err != nil {
		//  handle return based on error
		return solstructs.TokenTransferDetails{}, err
	}

	// convert the amount to actual amount
	actualTransferAmount := transferAmount/math.Pow(10, float64(decimals))

	returnValue := solstructs.TokenTransferDetails{
		Amount: actualTransferAmount,
		Mint: mint,
	}


	// return values
	return returnValue, nil

}

func GetTokenAccountInfo(rpcURL, tokenAccount string) (string, uint8, error) {
    // Prepare RPC request
    payload := map[string]interface{}{
        "jsonrpc": "2.0",
        "id":      1,
        "method":  "getAccountInfo",
        "params": []interface{}{
            tokenAccount,
            map[string]string{
                "encoding": "jsonParsed",
            },
        },
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return "", 0, fmt.Errorf("error marshaling request: %v", err)
    }

    // Make HTTP request
    resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", 0, fmt.Errorf("error making request: %v", err)
    }
    defer resp.Body.Close()

    // Parse response
    var rpcResp solstructs.AccountInfoResponse
    if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
        return "", 0, fmt.Errorf("error decoding response: %v", err)
    }

    // Extract mint and decimals
    mint := rpcResp.Result.Value.Data.Parsed.Info.Mint
    decimals := rpcResp.Result.Value.Data.Parsed.Info.TokenAmount.Decimals

    return mint, decimals, nil
}