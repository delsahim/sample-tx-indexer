package main

import (
	"context"
	"encoding/json"
	"indexer_golang/handlers"
	"indexer_golang/parsers"
	"indexer_golang/utils"
	websocketmethods "indexer_golang/websocket_methods"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coder/websocket/wsjson"
)



func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	const (
		NODE_WSS_URL = "wss://solemn-fluent-glitter.solana-mainnet.quiknode.pro/a4b0c2d7fa048c4818a5f20dd20018d16ebdc4d3"
		ALL_NODE_WSS = "wss://solana-rpc.publicnode.com"
		NODE_HTTP_URL = "https://solemn-fluent-glitter.solana-mainnet.quiknode.pro/a4b0c2d7fa048c4818a5f20dd20018d16ebdc4d3"
		ORCA_PROGRAM_ID = "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"
		COMMIT_LVL_FINALIZED = "finalized"
		BLOCK_FILE_NAME = "blockData"
		ORCA_IDL_FILENAME = "whirlpoolIDL.json"
	) 

	log.Println("Program Started ..........")

	// websocket connection
	log.Println("Establishing Websocket Connection ...............")
	conn, err := websocketmethods.ConnectWebsocket(NODE_WSS_URL)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.CloseNow()

	conn.SetReadLimit(10 * 1024 * 1024) // 10 MB

	log.Println("Connected to WebSocket ..........")

	// subscribe to a log
	log.Println("Sending subscription message")
	err = websocketmethods.BlockSubscrice(conn, COMMIT_LVL_FINALIZED, ORCA_PROGRAM_ID)
	if err != nil {
		log.Fatalf("Failed to send subscription request: %v", err)
	}
	log.Println("Subscription request sent")


	// Channel to handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)


	// deal with respnonse
	go func() {
		
		for {
			// notify that a new block has been received and read it 
			log.Println("Received new block notification ...........")
			var rawResponse json.RawMessage
			err := wsjson.Read(ctx, conn, &rawResponse)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}
			log.Println("Read new block notification ...........")


			// convert the message to json
			jsonMessage, err := parsers.MessageToJson(rawResponse)
			if err != nil {
				log.Printf("Failed to parse JSON: %v", err)
			}


			// check for subscription confirmation message
			if _, isResult := jsonMessage["result"]; isResult {
				log.Println("Received subscription confirmation, ignoring this message.")
				continue
			}

			block, err := parsers.BlockMessageToBlockStruct(rawResponse)
			if err != nil {
				log.Fatalf("Failed to parse block: %v", err)
			}

			for _, singleTransaction := range block.Params.Result.Value.Block.Transactions {
				// check for transaction error
				if singleTransaction.Meta.Err != nil {
					log.Println("Failed Transaction skipped")
					continue
				}

				transactionSignature := singleTransaction.Transaction.Signatures
				completeAccountKeys := append(singleTransaction.Transaction.Message.AccountKeys,singleTransaction.Meta.LoadedAddresses.Writable...)
				completeAccountKeys = append(completeAccountKeys, singleTransaction.Meta.LoadedAddresses.Readonly...)
				for _,accountKey := range singleTransaction.Transaction.Message.AccountKeys {
					if accountKey == ORCA_PROGRAM_ID {
						// find the instruction needed
						for _, instruction := range singleTransaction.Transaction.Message.Instructions {
							// deal with orca swap
							
							if singleTransaction.Transaction.Message.AccountKeys[instruction.ProgramIdIndex] == ORCA_PROGRAM_ID {
								// verify if it is a swap instruction
								// indexValue := instructionIndex
								// tokenAuth := completeAccountKeys[instruction.Accounts[1]]
								disc, err := utils.GetInstructionDiscriminator(instruction.Data)
								if err != nil {
									log.Println("Unable to get discriminator")
									continue
								}
								if disc == "f8c69e91e17587c8" {
									log.Println("Swap Found")
									log.Println("Inner transaction programs")
									//get the tokenA and tokenb account
									userTokenAccountA :=instruction.Accounts[3]
									userTokenAccountB := instruction.Accounts[5]
									liquidityTokenVaultA := instruction.Accounts[4]
									liquidityTokenVaultB := instruction.Accounts[6]
									userWallet := completeAccountKeys[instruction.Accounts[1]]
									// create the swap variables
									var (
										tokenMintA string
										tokenMintB string
										userPreBalanceA float64
										userPostBalanceA float64
										userPreBalanceB  float64
										userPostBalanceB float64
										lpPreBalanceA  float64
										lpPostBalanceA float64
										lpPreBalanceB  float64
										lpPostBalanceB float64
									)

									// get the pre token details 
									tokenMintA, tokenMintB, userPreBalanceA, userPreBalanceB, lpPreBalanceA, lpPreBalanceB =handlers.GetSwapTokenBalace(
																singleTransaction.Meta.PreTokenBalances,
																userTokenAccountA,
																userTokenAccountB,
																liquidityTokenVaultA,
																liquidityTokenVaultB,
															)

									// get the post token details 
									tokenMintA, tokenMintB, userPostBalanceA, userPostBalanceB, lpPostBalanceA, lpPostBalanceB =handlers.GetSwapTokenBalace(
										singleTransaction.Meta.PostTokenBalances,
										userTokenAccountA,
										userTokenAccountB,
										liquidityTokenVaultA,
										liquidityTokenVaultB,
									)
									aToB := handlers.GetAtoB(userPreBalanceA, userPostBalanceA)
									// determine buy and sell
									userTokenChangeA := math.Abs(userPostBalanceA - userPreBalanceA)
									userTokenChangeB := math.Abs(userPostBalanceB - userPreBalanceB)
									log.Printf("Transaction Hash %v \n",transactionSignature)
									log.Printf("A to b %v",aToB)
									log.Printf("User {%s}, swapped {%v} of token A {%s} for {%v} of token B {%s} \n",userWallet, userTokenChangeA, tokenMintA, userTokenChangeB, tokenMintB)
									log.Printf("Liquidity pool for token A; PreBalance {%v}  PostBalance {%v}",lpPreBalanceA,lpPostBalanceA)
									log.Printf("Liquidity pool for token B; PreBalance {%v}  PostBalance {%v}",lpPreBalanceB,lpPostBalanceB)									
								}
							}
						}
					}
				}
			}

			log.Println("parsed successfully ......")

		}
	}()


	// deal with keyboard interruption
	<-interrupt
	log.Println("Interrupt received, shutting down .......")
}