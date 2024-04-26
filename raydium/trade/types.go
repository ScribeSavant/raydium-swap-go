package trade

import (
	"bytes"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scribesavant/raydium-swap-go/raydium/constants"
	"github.com/scribesavant/raydium-swap-go/raydium/layouts"
)

type SwapV4Instruction struct {
	Instruction      uint8
	AmountIn         uint64
	MinimumOutAmount uint64
}

func (inst *SwapV4Instruction) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func NewSwapV4Instruction(
	connection *rpc.Client,
	poolKeys *layouts.ApiPoolInfoV4,
	amountIn uint64,
	minAmountOut uint64,
	tokenAccountIn solana.PublicKey,
	tokenAccountOut solana.PublicKey,
	signer solana.PrivateKey) (solana.Instruction, error) {
	inst := &SwapV4Instruction{
		Instruction:      9,
		AmountIn:         amountIn,
		MinimumOutAmount: minAmountOut,
	}
	keys := solana.AccountMetaSlice{}

	keys.Append(solana.Meta(solana.TokenProgramID))
	keys.Append(solana.Meta(poolKeys.ID).WRITE())
	keys.Append(solana.Meta(poolKeys.Authority))
	keys.Append(solana.Meta(poolKeys.OpenOrders).WRITE())

	if poolKeys.Version == 4 {
		keys.Append(solana.Meta(poolKeys.TargetOrders).WRITE())
	}

	keys.Append(solana.Meta(poolKeys.BaseVault).WRITE())
	keys.Append(solana.Meta(poolKeys.QuoteVault).WRITE())

	if poolKeys.Version == 5 {
		keys.Append(solana.Meta(constants.ModelDataPubkey).WRITE())
	}
	keys.Append(solana.Meta(poolKeys.MarketProgramId))
	keys.Append(solana.Meta(poolKeys.MarketId).WRITE())
	keys.Append(solana.Meta(poolKeys.MarketBids).WRITE())
	keys.Append(solana.Meta(poolKeys.MarketAsks).WRITE())
	keys.Append(solana.Meta(poolKeys.MarketEventQueue).WRITE())
	keys.Append(solana.Meta(poolKeys.BaseVault).WRITE())
	keys.Append(solana.Meta(poolKeys.QuoteVault).WRITE())
	keys.Append(solana.Meta(poolKeys.MarketAuthority))
	keys.Append(solana.Meta(tokenAccountIn).WRITE())
	keys.Append(solana.Meta(tokenAccountOut).WRITE())
	keys.Append(solana.Meta(signer.PublicKey()).SIGNER())

	data, err := inst.Data()
	if err != nil {
		return nil, err
	}
	intstr := solana.NewInstruction(
		constants.RAYDIUM_V4_PROGRAM_ID,
		keys,
		data,
	)

	return intstr, nil
}
