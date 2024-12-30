package parsers

import (
	"encoding/json"
)

func MessageToJson(logMessage []byte) (map[string]interface{}, error) {
	var parsedMessage map[string]interface{}
	if err := json.Unmarshal(logMessage, &parsedMessage); err != nil {
		return nil, err
	}
	return parsedMessage, nil
}