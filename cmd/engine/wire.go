//+build wireinject

package main

import (
	"github.com/google/wire"
	"option-dance/match"
)

func InitExchange() (*match.Exchange, error) {
	wire.Build(
		exchangeSet,
		match.NewExchange,
	)
	return &match.Exchange{}, nil
}
