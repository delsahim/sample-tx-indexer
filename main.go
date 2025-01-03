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
				// check the main instructions
				for _,accountKey := range singleTransaction.Transaction.Message.AccountKeys {
					if accountKey == ORCA_PROGRAM_ID {
						// find the instruction needed
						for _, instruction := range singleTransaction.Transaction.Message.Instructions {
							// deal with orca swap
							
							if singleTransaction.Transaction.Message.AccountKeys[instruction.ProgramIdIndex] == ORCA_PROGRAM_ID {
								// verify if it is a swap instruction
								// indexValue := instructionIndex
								// tokenAuth := completeAccountKeys[instruction.Accounts[1]]
								disc, err := utils.GetEightByteDiscriminator(instruction.Data)
								if err != nil {
									log.Println("Unable to get discriminator")
									continue
								}
								if disc == utils.GetAnchorDiscriminatorFromInstructionName("swap") {
									log.Println("Swap Found")
									swapArgs, err := handlers.DecodeOrcaSwapData(instruction.Data)
									if err != nil {
										log.Println("error in retrieving data")
										
									}
									log.Printf("swap data args: %v",swapArgs)
									//get the tokenA and tokenb account
									userTokenAccountA :=instruction.Accounts[3]
									userTokenAccountB := instruction.Accounts[5]
									liquidityTokenVaultA := instruction.Accounts[4]
									liquidityTokenVaultB := instruction.Accounts[6]
									userWallet := completeAccountKeys[instruction.Accounts[1]]
									pairPoolAccount := completeAccountKeys[instruction.Accounts[2]]
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
									log.Printf("A to b %v\n",aToB)
									log.Printf("Pair Pool: %s\n", pairPoolAccount)
									log.Printf("User {%s}, swapped {%v} of token A {%s} for {%v} of token B {%s} \n",userWallet, userTokenChangeA, tokenMintA, userTokenChangeB, tokenMintB)
									log.Printf("Liquidity pool for token A; PreBalance {%v}  PostBalance {%v}",lpPreBalanceA,lpPostBalanceA)
									log.Printf("Liquidity pool for token B; PreBalance {%v}  PostBalance {%v}",lpPreBalanceB,lpPostBalanceB)									
								} else if (disc == utils.GetAnchorDiscriminatorFromInstructionName("initializePool")) {
									log.Println("initializing a new pool")
								} else if (disc == utils.GetAnchorDiscriminatorFromInstructionName("collectProtocolFees")) {
									log.Println("decrese liquidity found")
								} else {
									log.Println("unidentified instruction")
								}
							}
						}
					}
				}

				// check the inner instructions
				
				for innerInstructionRange, innerInstruction := range singleTransaction.Meta.InnerInstructions {
					for range2, instruction := range innerInstruction.Instructions {
						// check for the orca program in the instruction id
						if completeAccountKeys[instruction.ProgramIdIndex] == ORCA_PROGRAM_ID {

							log.Println("Inner transaction found")

							// log.Printf("Instruction Data: %v and discriminator: %v",instruction.Data,)
							

							// decode the inner transaction first to confirm that it is a right transaction
							disc, err := utils.GetEightByteDiscriminator(instruction.Data)
							if err != nil {
								log.Println("Unable to get discriminator")
								continue
							}
							if disc == utils.GetAnchorDiscriminatorFromInstructionName("swap") {
								// handle the swap discriminator
								log.Println("innner swap fpund, rinting tokenkeg details")
								log.Printf("Transaction Hash: %v", transactionSignature)
								tokenInstruction1 := singleTransaction.Meta.InnerInstructions[innerInstructionRange].Instructions[range2+1]
								
								tokenInstruction2 := singleTransaction.Meta.InnerInstructions[innerInstructionRange].Instructions[range2+2]
								
								tokenDecoded1, err := handlers.DecodeSystemTransfer(completeAccountKeys,tokenInstruction1,NODE_HTTP_URL)
								if err != nil {
									log.Println("Error in getting token details")
								}
								log.Printf("Decoded token: %v", tokenDecoded1)
								tokenDecoded2, err := handlers.DecodeSystemTransfer(completeAccountKeys,tokenInstruction2,NODE_HTTP_URL)
								if err != nil {
									log.Println("Error in getting token details")
								}
								log.Printf("Decoded token: %v", tokenDecoded2)
							}

							log.Println("inner tx handled")




							



							// handle the inner insrruction
							// log.Println("Orca inner transaction found")
							// disc, err := utils.GetEightByteDiscriminator(instruction.Data)
							// 	if err != nil {
							// 		log.Println("Unable to get discriminator")
							// 		continue
							// 	}
							// 	if disc == utils.GetAnchorDiscriminatorFromInstructionName("swap") {
							// 		log.Println("Swap Found")
							// 		log.Println("Inner transaction programs")
							// 		//get the tokenA and tokenb account
							// 		userTokenAccountA :=instruction.Accounts[3]
							// 		userTokenAccountB := instruction.Accounts[5]
							// 		liquidityTokenVaultA := instruction.Accounts[4]
							// 		liquidityTokenVaultB := instruction.Accounts[6]
							// 		userWallet := completeAccountKeys[instruction.Accounts[1]]
							// 		pairPoolAccount := completeAccountKeys[instruction.Accounts[2]]
							// 		// create the swap variables
							// 		var (
							// 			tokenMintA string
							// 			tokenMintB string
							// 			userPreBalanceA float64
							// 			userPostBalanceA float64
							// 			userPreBalanceB  float64
							// 			userPostBalanceB float64
							// 			lpPreBalanceA  float64
							// 			lpPostBalanceA float64
							// 			lpPreBalanceB  float64
							// 			lpPostBalanceB float64
							// 		)

							// 		// get the pre token details 
							// 		tokenMintA, tokenMintB, userPreBalanceA, userPreBalanceB, lpPreBalanceA, lpPreBalanceB =handlers.GetSwapTokenBalace(
							// 									singleTransaction.Meta.PreTokenBalances,
							// 									userTokenAccountA,
							// 									userTokenAccountB,
							// 									liquidityTokenVaultA,
							// 									liquidityTokenVaultB,
							// 								)

							// 		// get the post token details 
							// 		tokenMintA, tokenMintB, userPostBalanceA, userPostBalanceB, lpPostBalanceA, lpPostBalanceB =handlers.GetSwapTokenBalace(
							// 			singleTransaction.Meta.PostTokenBalances,
							// 			userTokenAccountA,
							// 			userTokenAccountB,
							// 			liquidityTokenVaultA,
							// 			liquidityTokenVaultB,
							// 		)
							// 		aToB := handlers.GetAtoB(userPreBalanceA, userPostBalanceA)
							// 		// determine buy and sell
							// 		userTokenChangeA := math.Abs(userPostBalanceA - userPreBalanceA)
							// 		userTokenChangeB := math.Abs(userPostBalanceB - userPreBalanceB)
							// 		log.Printf("Inner Transaction Hash %v \n",transactionSignature)
							// 		log.Printf("A to b %v\n",aToB)
							// 		log.Printf("Pair Pool: %s\n", pairPoolAccount)
							// 		log.Printf("User {%s}, swapped {%v} of token A {%s} for {%v} of token B {%s} \n",userWallet, userTokenChangeA, tokenMintA, userTokenChangeB, tokenMintB)
							// 		log.Printf("Liquidity pool for token A; PreBalance {%v}  PostBalance {%v}",lpPreBalanceA,lpPostBalanceA)
							// 		log.Printf("Liquidity pool for token B; PreBalance {%v}  PostBalance {%v}",lpPreBalanceB,lpPostBalanceB)									
							// 	} else if (disc == utils.GetAnchorDiscriminatorFromInstructionName("initializePool")) {
							// 		log.Println("initializing a new pool")
							// 	} else if (disc == utils.GetAnchorDiscriminatorFromInstructionName("collectProtocolFees")) {
							// 		log.Println("decrese liquidity found")
							// 	} else {
							// 		log.Println("unidentified instruction")
							// 	}
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