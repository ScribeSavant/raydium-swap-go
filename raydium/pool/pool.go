package pool

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scribesavant/raydium-swap-go/raydium/constants"
	"github.com/scribesavant/raydium-swap-go/raydium/layouts"
	"github.com/scribesavant/raydium-swap-go/raydium/utils"
)

type Pool struct {
	connection *rpc.Client
}

func New(connection *rpc.Client) *Pool {
	return &Pool{
		connection: connection,
	}
}

func (p *Pool) getProgramAccounts(mint1 string, mint2 string) (rpc.GetProgramAccountsResult, error) {
	var layout layouts.LIQUIDITY_STATE_LAYOUT_V4

	return p.connection.GetProgramAccountsWithOpts(context.Background(), constants.RAYDIUM_V4_PROGRAM_ID, &rpc.GetProgramAccountsOpts{
		Filters: []rpc.RPCFilter{
			{
				DataSize: layout.Span(),
			},
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: layout.Offset("BaseMint"),
					Bytes:  solana.MustPublicKeyFromBase58(mint1).Bytes(),
				},
			},
			{
				Memcmp: &rpc.RPCFilterMemcmp{
					Offset: layout.Offset("QuoteMint"),
					Bytes:  solana.MustPublicKeyFromBase58(mint2).Bytes(),
				},
			},
		},
	})
}

func (p *Pool) GetPoolKeys(mint1 string, mint2 string) (*layouts.ApiPoolInfoV4, error) {
	var layout layouts.LIQUIDITY_STATE_LAYOUT_V4
	var marketLayout layouts.MarketStateLayoutV3

	programAccounts, err := p.getProgramAccounts(mint1, mint2)

	if err != nil {
		return &layouts.ApiPoolInfoV4{}, errors.New("Pool not found")
	}
	if len(programAccounts) == 0 {
		return &layouts.ApiPoolInfoV4{}, errors.New("Pool not found")
	}

	programAccount := programAccounts[0]

	layout.Decode(programAccount.Account.Data.GetBinary())
	marketAccount, err := p.connection.GetAccountInfo(context.Background(), layout.MarketId)

	if err != nil {
		return &layouts.ApiPoolInfoV4{}, err
	}

	marketLayout.Decode(marketAccount.Value.Data.GetBinary())

	authority, _, err := solana.FindProgramAddress([][]byte{{97, 109, 109, 32, 97, 117, 116, 104, 111, 114, 105, 116, 121}}, constants.RAYDIUM_V4_PROGRAM_ID)

	if err != nil {
		return &layouts.ApiPoolInfoV4{}, nil
	}

	marketAuthority, _, err := utils.GetAssociatedAuthority(marketAccount.Value.Owner, marketLayout.OwnAddress)

	if err != nil {
		return &layouts.ApiPoolInfoV4{}, err
	}

	return &layouts.ApiPoolInfoV4{
		ID:                 programAccount.Pubkey,
		BaseMint:           layout.BaseMint,
		QuoteMint:          layout.QuoteMint,
		LpMint:             layout.LpMint,
		BaseDecimals:       layout.BaseDecimal,
		QuoteDecimals:      layout.QuoteDecimal,
		LpDecimals:         layout.BaseDecimal,
		Version:            4,
		ProgramId:          constants.RAYDIUM_V4_PROGRAM_ID,
		OpenOrders:         layout.OpenOrders,
		TargetOrders:       layout.TargetOrders,
		BaseVault:          layout.BaseVault,
		QuoteVault:         layout.QuoteVault,
		MarketVersion:      3,
		Authority:          authority,
		MarketProgramId:    marketAccount.Value.Owner,
		MarketId:           marketLayout.OwnAddress,
		MarketAuthority:    marketAuthority,
		MarketBaseVault:    marketLayout.BaseVault,
		MarketQuoteVault:   marketLayout.QuoteVault,
		MarketBids:         marketLayout.Bids,
		MarketAsks:         marketLayout.Asks,
		MarketEventQueue:   marketLayout.EventQueue,
		WithdrawQueue:      layout.WithdrawQueue,
		LpVault:            layout.LpVault,
		LookupTableAccount: solana.PublicKey{},
	}, nil
}
