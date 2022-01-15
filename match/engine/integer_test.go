package engine

import (
	"github.com/MixinNetwork/go-number"
	"github.com/shopspring/decimal"
	"testing"
)

func TestInteger(t *testing.T) {
	integer := number.NewInteger(100, 4)
	t.Log(integer.Value())
	t.Log(integer.Decimal())

	remainingAmount := number.FromString("2")
	totalAmount := number.FromString("3")
	remainingRatio := remainingAmount.Div(totalAmount)
	//sell put

	transferAmount := number.FromString("15000").Mul(remainingRatio).Round(8).String()
	t.Log(transferAmount)

	r, _ := decimal.NewFromString("2")
	to, _ := decimal.NewFromString("3")
	ratio := r.Div(to)
	m, _ := decimal.NewFromString("15000")
	mul := m.Mul(ratio).String()
	t.Log(mul)
}
