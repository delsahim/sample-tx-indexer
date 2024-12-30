package websocketmethods

import (
	"context"
	"time"

	"github.com/coder/websocket"
)


func ConnectWebsocket(nodeUrl string) (*websocket.Conn, error) {
	// Create websocket connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, nodeUrl, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
