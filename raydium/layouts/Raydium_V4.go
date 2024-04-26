package layouts

import (
	"encoding/base64"
	"reflect"
	"unsafe"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"lukechampine.com/uint128"
)

type ApiPoolInfoV4 struct {
	ID                 solana.PublicKey `json:"id"`
	BaseMint           solana.PublicKey `json:"baseMint"`
	QuoteMint          solana.PublicKey `json:"quoteMint"`
	LpMint             solana.PublicKey `json:"lpMint"`
	BaseDecimals       uint64           `json:"baseDecimals"`
	QuoteDecimals      uint64           `json:"quoteDecimals"`
	LpDecimals         uint64           `json:"lpDecimals"`
	Version            uint64           `json:"version"`
	ProgramId          solana.PublicKey `json:"programId"`
	Authority          solana.PublicKey `json:"authority"`
	OpenOrders         solana.PublicKey `json:"openOrders"`
	TargetOrders       solana.PublicKey `json:"targetOrders"`
	BaseVault          solana.PublicKey `json:"baseVault"`
	QuoteVault         solana.PublicKey `json:"quoteVault"`
	WithdrawQueue      solana.PublicKey `json:"withdrawQueue"`
	LpVault            solana.PublicKey `json:"lpVault"`
	MarketVersion      uint64           `json:"marketVersion"`
	MarketProgramId    solana.PublicKey `json:"marketProgramId"`
	MarketId           solana.PublicKey `json:"marketId"`
	MarketAuthority    solana.PublicKey `json:"marketAuthority"`
	MarketBaseVault    solana.PublicKey `json:"marketBaseVault"`
	MarketQuoteVault   solana.PublicKey `json:"marketQuoteVault"`
	MarketBids         solana.PublicKey `json:"marketBids"`
	MarketAsks         solana.PublicKey `json:"marketAsks"`
	MarketEventQueue   solana.PublicKey `json:"marketEventQueue"`
	LookupTableAccount solana.PublicKey `json:"lookupTableAccount"`
}

type LIQUIDITY_STATE_LAYOUT_V4 struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64
	SwapBaseInAmount       uint128.Uint128
	SwapQuoteOutAmount     uint128.Uint128
	SwapBase2QuoteFee      uint64
	SwapQuoteInAmount      uint128.Uint128
	SwapBaseOutAmount      uint128.Uint128
	SwapQuote2BaseFee      uint64
	BaseVault              solana.PublicKey
	QuoteVault             solana.PublicKey
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	LpMint                 solana.PublicKey
	OpenOrders             solana.PublicKey
	MarketId               solana.PublicKey
	MarketProgramId        solana.PublicKey
	TargetOrders           solana.PublicKey
	WithdrawQueue          solana.PublicKey
	LpVault                solana.PublicKey
	Owner                  solana.PublicKey
	LpReserve              uint64
	Padding                [3]uint64
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) Span() uint64 {
	return uint64(unsafe.Sizeof(*l))
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) Offset(value string) uint64 {
	fieldType, found := reflect.TypeOf(*l).FieldByName(value)
	if !found {
		return 0
	}
	return uint64(fieldType.Offset)
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) DecodeBase64(data string) error {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return l.Decode(decoded)
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) Decode(data []byte) error {
	err := bin.UnmarshalBorsh(&l, data)
	return err
}
