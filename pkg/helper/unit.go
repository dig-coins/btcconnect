package helper

import "github.com/shopspring/decimal"

const (
	SatoshiPerBitcoin = 1e8
)

func UnitBTC2SatoshiBTC(v float64) int64 {
	return decimal.NewFromFloat(v).Mul(decimal.NewFromInt(SatoshiPerBitcoin)).IntPart()
}
