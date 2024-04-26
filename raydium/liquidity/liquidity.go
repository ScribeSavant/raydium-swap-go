package liquidity

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scribesavant/raydium-swap-go/raydium/layouts"
	"github.com/scribesavant/raydium-swap-go/raydium/utils"
)

type Liquidity struct {
	connection *rpc.Client
}

func New(connection *rpc.Client) *Liquidity {
	return &Liquidity{
		connection: connection,
	}
}

func (l *Liquidity) FetchInfo(poolKeys *layouts.ApiPoolInfoV4) (*LiquidityPoolInfo, error) {
	var LiquidityPoolInfo LiquidityPoolInfo
	instructions := l.makeSimulatePoolInfoInstruction(poolKeys)

	logs, err := l.simulateAmountsOut(instructions)
	if err != nil {
		return nil, err
	}

	if len(logs.Value.Logs) == 0 {
		return nil, fmt.Errorf("pool info unavailable")
	}

	for _, log := range logs.Value.Logs {
		if strings.Contains(log, "GetPoolData") {
			jsonLog := l.parseLog2Json(log, "GetPoolData")
			json.Unmarshal([]byte(jsonLog), &LiquidityPoolInfo)
			return &LiquidityPoolInfo, nil
		}
	}
	return nil, fmt.Errorf("pool info unavailable")
}

func (l *Liquidity) GetAmountsOut(poolKeys *layouts.ApiPoolInfoV4, amountIn *utils.TokenAmount, slippage *utils.Percent) (*AmountsOut, error) {
	poolInfo, err := l.FetchInfo(poolKeys)

	if err != nil {
		return &AmountsOut{}, err
	}

	reserves := []uint64{poolInfo.BaseReserve, poolInfo.QuoteReserve}
	tokens := []utils.Token{
		*utils.NewToken("", poolKeys.BaseMint.String(), poolInfo.BaseDecimals),
		*utils.NewToken("", poolKeys.QuoteMint.String(), poolInfo.QuoteDecimals),
	}

	if amountIn.Mint != poolKeys.BaseMint.String() { // Reverse reserves for quote swap
		for i, j := 0, len(reserves)-1; i < j; i, j = i+1, j-1 {
			reserves[i], reserves[j] = reserves[j], reserves[i]
		}

		for i, j := 0, len(tokens)-1; i < j; i, j = i+1, j-1 {
			tokens[i], tokens[j] = tokens[j], tokens[i]
		}

	}

	reserverIn, reserveOut := big.NewInt(int64(reserves[0])), big.NewInt(int64(reserves[1]))
	inTok, outTok := tokens[0], tokens[1]

	amountIn = utils.NewTokenAmount(&inTok, amountIn.Amount*math.Pow(10, float64(inTok.Decimals)))

	denominator := reserverIn.Add(reserverIn, big.NewInt(int64(amountIn.Amount)))
	_amountOut, _ := reserveOut.Mul(reserveOut, big.NewInt(int64(amountIn.Amount))).Div(reserveOut, denominator).Float64()
	amountOut := utils.NewTokenAmount(&outTok, _amountOut)
	minAmountOut := utils.NewTokenAmount(&outTok, float64(uint64(amountOut.Amount)*uint64(float64(slippage.Denominator)-float64(slippage.Numerator))/slippage.Denominator))

	return &AmountsOut{
		AmountIn:     amountIn,
		AmountOut:    amountOut,
		MinAmountOut: minAmountOut,
	}, nil
}

func (l *Liquidity) parseLog2Json(log string, keyword string) string {
	jsonData := strings.Split(log, keyword+": ")[1]
	return jsonData
}

func (l *Liquidity) makeSimulatePoolInfoInstruction(poolKeys *layouts.ApiPoolInfoV4) []solana.Instruction {
	layout := &SimulateStruct{
		Instruction:  12,
		SimulateType: 0,
	}
	data, err := layout.Encode()

	if err != nil {
		panic(err)
	}

	keys := solana.AccountMetaSlice{}
	keys.Append(solana.Meta(poolKeys.ID))
	keys.Append(solana.Meta(poolKeys.Authority))
	keys.Append(solana.Meta(poolKeys.OpenOrders))
	keys.Append(solana.Meta(poolKeys.BaseVault))
	keys.Append(solana.Meta(poolKeys.QuoteVault))
	keys.Append(solana.Meta(poolKeys.LpMint))
	keys.Append(solana.Meta(poolKeys.MarketId))
	keys.Append(solana.Meta(poolKeys.MarketEventQueue))

	return []solana.Instruction{
		solana.NewInstruction(
			poolKeys.ProgramId,
			keys,
			data,
		),
	}
}

func (l *Liquidity) simulateAmountsOut(instructions []solana.Instruction) (*rpc.SimulateTransactionResponse, error) {
	feePayer := solana.MustPublicKeyFromBase58("RaydiumSimuLateTransaction11111111111111111")
	recent, err := l.connection.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return &rpc.SimulateTransactionResponse{}, err
	}
	tx, _ := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(feePayer),
	)
	tx.Signatures = make([]solana.Signature, 1)
	tx.Signatures[0] = solana.MustSignatureFromBase58("1111111111111111111111111111111111111111111111111111111111111111") // If you know better way to do this, feel free to change
	return l.connection.SimulateTransaction(context.Background(), tx)
}
