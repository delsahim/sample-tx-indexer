# Solana Transaction Indexer

A high-performance indexer for tracking and decoding Solana DEX transactions.

## Overview
This indexer monitors Solana blockchain transactions in real-time, specifically focusing on decentralized exchange (DEX) interactions. It utilizes WebSocket connections to receive block notifications and processes transactions that interact with specified programs.

## Architecture

### WebSocket Connection
- Maintains a real-time connection to Solana network
- Receives block notifications for targeted program interactions
- Ensures minimal latency for transaction processing

### Transaction Processing

#### DEX Transaction Sources
1. **Direct Interactions**
   - Located in the main instruction array
   - Identified by the `programId` field in instructions
   - Direct user interactions with DEX programs

2. **Aggregator-Initiated Transactions**
   - Found in the `innerInstructions` array
   - Typically complex transactions routed through DEX aggregators
   - May contain multiple nested swaps

### Swap Transaction Decoding

#### Token Program Transfers
Each swap transaction typically involves two token program transfers (`TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA`):
1. User → Pool Account: Initial token transfer
2. Pool Account → User: Received token transfer

#### Instruction Decoding Process
1. **Base58 Decoding**
   - All instruction data is base58 encoded
   - Must be decoded before processing

2. **Data Structure**
   - Instruction discriminator: Identifies the specific program instruction
   - Arguments: Transaction-specific parameters defined in the IDL

## Helper Functions

### Transfer Processing
1. `GetTransferAmount`
   - Decodes token program transfer data
   - Returns raw transfer amount
   - Input: Base58 encoded instruction data
   - Output: Transfer amount (uint64)

2. `DecodeTransfer`
   - Comprehensive transfer instruction decoder
   - Returns both UI amount and token mint address
   - Handles decimal precision conversion
   - Input: Transfer instruction
   - Output: `{uiAmount: float64, tokenMint: string}`

### DEX-Specific Functions
1. `DecodeOrcaSwapData`
   - Specialized decoder for Orca DEX instructions
   - Extracts swap-specific parameters
   - Input: Orca instruction data
   - Output: Decoded swap parameters

## Setup and Usage

[Add setup instructions, configuration details, and usage examples here]

## Contributing

[Add contribution guidelines here]

## License

[Add license information here]