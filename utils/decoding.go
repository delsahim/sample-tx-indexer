package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/mr-tron/base58"
)

func GetInstructionDiscriminator(instructionData string) (string, error){
	decodedData, err := base58.Decode(instructionData)
    if err != nil {
        return "", err
    }

	// Get first 8 bytes and convert to hex
    if len(decodedData) < 8 {
        return "", fmt.Errorf("decoded data too short, expected at least 8 bytes")
    }
    
    discriminator := hex.EncodeToString(decodedData[:8])
    return discriminator, nil
}

func GetOneByteDiscriminator(instructionData string) (string, error) {
    // Decode base58 data
    decodedData, err := base58.Decode(instructionData)
    if err != nil {
        return "", err
    }

    // Check if we have at least 1 byte
    if len(decodedData) < 1 {
        return "", fmt.Errorf("decoded data too short, expected at least 1 byte")
    }
    
    // Get first byte and convert to hex
    discriminator := hex.EncodeToString(decodedData[:1])
    return discriminator, nil
}