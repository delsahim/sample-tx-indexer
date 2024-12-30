package parsers

import (
	"encoding/json"
	solstructs "indexer_golang/solStructs"
)


func BlockMessageToBlockStruct(message []byte) (solstructs.SolanaBlockSubscribe, error) {
	var jsonRPCMessage solstructs.SolanaBlockSubscribe
	if err := json.Unmarshal(message, &jsonRPCMessage); err != nil {
		return solstructs.SolanaBlockSubscribe{}, err
	}
	return jsonRPCMessage, nil
}