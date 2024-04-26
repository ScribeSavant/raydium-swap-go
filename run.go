package main

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scribesavant/raydium-swap-go/raydium"
	"github.com/scribesavant/raydium-swap-go/raydium/utils"
)

func main() {

	executeTransaction := false

	connection := rpc.New(os.Getenv("RPC_URL"))
	raydium := raydium.New(connection, os.Getenv("WALLET_PRIVATE_KEY"))

	inputToken := utils.NewToken("SOL", "So11111111111111111111111111111111111111112", 9)
	outputToken := utils.NewToken("RAY", "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R", 6)

	slippage := utils.NewPercent(1, 100) // 1% slippage (for 0.5 set second parameter to "1000" example: utils.NewPercent(5, 1000) )

	amount := utils.NewTokenAmount(inputToken, 0.1) // 0.1 sol

	poolKeys, err := raydium.Pool.GetPoolKeys(inputToken.Mint, outputToken.Mint)

	if err != nil {
		panic(err)
	}

	amountsOut, err := raydium.Liquidity.GetAmountsOut(poolKeys, amount, slippage)
	if err != nil {
		panic(err)
	}

	tx, err := raydium.Trade.MakeSwapTransaction(poolKeys, amountsOut.AmountIn, amountsOut.MinAmountOut)

	if err != nil {
		panic(err)
	}

	if !executeTransaction {
		simRes, err := connection.SimulateTransaction(context.Background(), tx)

		if err != nil {
			spew.Dump(err.Error())
			return
		}

		spew.Dump(simRes)
	} else {
		signature, err := connection.SendTransactionWithOpts(context.Background(), tx, rpc.TransactionOpts{SkipPreflight: true})

		if err != nil {
			panic(err)
		}
		fmt.Println("Transaction successfully sent: ", signature)
	}
}
