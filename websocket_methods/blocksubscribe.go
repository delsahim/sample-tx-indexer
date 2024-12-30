package websocketmethods

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)


type SolanaRPCRequest struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func BlockSubscrice(conn *websocket.Conn, commitLevel, programID string) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Subscription request
	subscriptionParams := SolanaRPCRequest{
		JsonRPC: "2.0",
		ID:      1,
		Method:  "blockSubscribe",
		Params: []interface{}{
			map[string]string{
				"mentionsAccountOrProgram": programID,
			},
			map[string]interface{}{
				"commitment":                     commitLevel,
				"encoding":                       "json",
				"showRewards":                    false,
				"transactionDetails":             "full",
				"maxSupportedTransactionVersion": 0,
			},
		},
	}
	
	
	if err := wsjson.Write(ctx, conn, subscriptionParams); err != nil {
		return err
	}

	return nil
} 