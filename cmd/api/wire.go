//+build wireinject

package main

import (
	"github.com/google/wire"
	"option-dance/handler/api"
)

func InitApiServer() (api.ServerApi, error) {
	wire.Build(
		apiServerSet,
		apiDbMarketSet,
		api.NewServer,
	)
	return api.ServerApi{}, nil
}
