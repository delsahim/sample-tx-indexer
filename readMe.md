# SAMPLE INDEXER REPO

## HOW THE INDEXER WORKS 
The Indexer works by using a websocket connection to receive block notifications for solana transactions that interact with the selected programs 
### KEY POINTS
1. For dex transactions there are two places where the required data can be found 
    - For direct interactions; the data is in the instruction array and can be identified by the programId field of any instruction 
    - For transactions initiated by dex aggregators; the instructions are found in the innerInstructions array

2. Decoding A swap transaction; Generally when a swap occurs, there are two calls to the token program "transfer" instruction (TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA)
    - The user sends the token to be swapped into the pool acct
    - The pool acct sends the equivalent into the user other token acct

3. Decoding Instruction; all instruction data are base58 encoded, and you have to first decode them, the instruction generally consists of discriminator and args;
    - the discriminator acts as a pointer to the instruction called in the program 
    - the remaining data is decoded using the idl

### HELPER FUNCTIONS
1. GetTransferAmount: Decodes a tokent program transfer data and reurns the amount 
2. Decodes a transfer: instruction and returns the UIAmount and Token Mint
3. DecodeOrcaSwapData:: Decodes an orca instruction data