package liquidity

import (
	bin "github.com/gagliardetto/binary"
	"github.com/scribesavant/raydium-swap-go/raydium/utils"
)

type SimulateStruct struct {
	Instruction  uint8
	SimulateType uint8
}

func (l *SimulateStruct) Encode() ([]byte, error) {
	return bin.MarshalBin(&l)
}

type LiquidityPoolInfo struct {
	Status        int    `json:"status"`
	BaseDecimals  uint64 `json:"coin_decimals"`
	QuoteDecimals uint64 `json:"pc_decimals"`
	LpDecimals    uint64 `json:"lp_decimals"`
	QuoteReserve  uint64 `json:"pool_pc_amount"`
	BaseReserve   uint64 `json:"pool_coin_amount"`
	LpSupply      int    `json:"pool_lp_supply"`
	PoolOpenTime  int64  `json:"pool_open_time"`
	AmmID         string `json:"amm_id"`
}

type AmountsOut struct {
	AmountIn     *utils.TokenAmount
	AmountOut    *utils.TokenAmount
	MinAmountOut *utils.TokenAmount
}
