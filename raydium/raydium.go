package raydium

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/scribesavant/raydium-swap-go/raydium/liquidity"
	"github.com/scribesavant/raydium-swap-go/raydium/pool"
	"github.com/scribesavant/raydium-swap-go/raydium/trade"
)

type Raydium struct {
	connection *rpc.Client
	Pool       *pool.Pool
	Liquidity  *liquidity.Liquidity
	Signer     solana.PrivateKey
	Trade      *trade.Trade
}

func New(connection *rpc.Client, privKey string) *Raydium {
	signer := solana.MustPrivateKeyFromBase58(privKey)
	r := &Raydium{
		connection: connection,
		Pool:       pool.New(connection),
		Liquidity:  liquidity.New(connection),
		Signer:     signer,
		Trade:      trade.New(connection, &signer),
	}

	return r
}
