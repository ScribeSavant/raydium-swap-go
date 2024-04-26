package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"math/big"

	"github.com/gagliardetto/solana-go"
)

func FromWei(amount string, decimals int64) string {
	parsedDes := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	amountInt := new(big.Float)
	amountInt.SetString(amount)
	amountInt.Quo(amountInt, new(big.Float).SetInt(parsedDes))
	return amountInt.Text('f', -1)
}

func ToWei(amount string, decimals int64) string {
	parsedDes := new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil)
	amountInt := new(big.Float)
	amountInt.SetString(amount)
	amountInt.Mul(amountInt, new(big.Float).SetInt(parsedDes))
	amountIntInt, _ := amountInt.Int(nil)
	return amountIntInt.String()
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

type Percent struct {
	Numerator   uint64
	Denominator uint64
}

func NewPercent(numerator uint64, denominator uint64) *Percent {
	return &Percent{
		Numerator:   numerator,
		Denominator: denominator,
	}
}

func (p *Percent) ToDecimal() float64 {
	return float64(p.Numerator) / float64(p.Denominator) * 100
}

type Token struct {
	Name     string
	Mint     string
	Decimals uint64
}

func NewToken(name string, mint string, decimals uint64) *Token {
	return &Token{
		Name:     name,
		Mint:     mint,
		Decimals: decimals,
	}
}

func (t *Token) PublicKey() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(t.Mint)
}

type TokenAmount struct {
	Token
	Amount float64
}

func NewTokenAmount(token *Token, amount float64) *TokenAmount {
	return &TokenAmount{
		Token:  *token,
		Amount: amount,
	}
}

func (t *TokenAmount) Mul(amount float64) *TokenAmount {
	return &TokenAmount{
		Token:  t.Token,
		Amount: t.Amount * amount,
	}
}

func (t *TokenAmount) Div(amount float64) *TokenAmount {
	return &TokenAmount{
		Token:  t.Token,
		Amount: t.Amount / amount,
	}
}

func (t *TokenAmount) Add(amount *TokenAmount) *TokenAmount {
	return &TokenAmount{
		Token:  t.Token,
		Amount: t.Amount + amount.Amount,
	}
}

func (t *TokenAmount) Sub(amount float64) *TokenAmount {
	return &TokenAmount{
		Token:  t.Token,
		Amount: t.Amount - amount,
	}
}

func GetAssociatedAuthority(programID solana.PublicKey, marketID solana.PublicKey) (solana.PublicKey, uint8, error) {
	seeds := [][]byte{marketID.Bytes()}
	var nonce uint8 = 0

	for nonce < 100 {
		seedsWithNonce := append(seeds, int8ToBuf(nonce))
		seedsWithNonce = append(seedsWithNonce, make([]byte, 7)) // Buffer.alloc(7)

		publicKey, err := solana.CreateProgramAddress(seedsWithNonce, programID)
		if err != nil {
			nonce++
			continue
		}

		return publicKey, nonce, nil
	}

	return solana.PublicKey{}, 0, errors.New("unable to find a viable program address nonce")
}

func int8ToBuf(value uint8) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, value)
	return buf.Bytes()
}
