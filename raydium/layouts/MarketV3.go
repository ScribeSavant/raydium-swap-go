package layouts

import (
	"encoding/base64"
	"reflect"
	"unsafe"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type MarketStateLayoutV3 struct {
	AccountFlag            [5]byte
	Padding                [8]byte
	OwnAddress             solana.PublicKey
	VaultSignerNonce       uint64
	BaseMint               solana.PublicKey
	QuoteMint              solana.PublicKey
	BaseVault              solana.PublicKey
	BaseDepositsTotal      uint64
	BaseFeesAccrued        uint64
	QuoteVault             solana.PublicKey
	QuoteDepositsTotal     uint64
	QuoteFeesAccrued       uint64
	QuoteDustThreshold     uint64
	RequestQueue           solana.PublicKey
	EventQueue             solana.PublicKey
	Bids                   solana.PublicKey
	Asks                   solana.PublicKey
	BaseLotSize            uint64
	QuoteLotSize           uint64
	FeeRateBps             uint64
	ReferrerRebatesAccrued uint64
	PaddingEnd             [7]byte
}

func (l MarketStateLayoutV3) Span() uint64 {
	return uint64(unsafe.Sizeof(l)) - 4 // Blame golang
}

func (l *MarketStateLayoutV3) DecodeBase64(data string) error {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return l.Decode(decoded)
}

func (l *MarketStateLayoutV3) Decode(data []byte) error {
	err := bin.UnmarshalBorsh(&l, data)
	return err
}

func (l *MarketStateLayoutV3) Offset(value string) uint64 {
	fieldType, found := reflect.TypeOf(*l).FieldByName(value)
	if !found {
		return 0
	}
	return uint64(fieldType.Offset)
}

// Serum market
