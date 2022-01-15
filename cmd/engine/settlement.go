package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"option-dance/core"
	"option-dance/service/mtg"
	"strings"
	"time"
)

var (
	date         string
	act          string
	deliveryType string
)

const (
	ActionRun    = "run"
	ActionGather = "gather"
)

var SettlementCmd = &cobra.Command{
	Use:   "settlement",
	Short: "settlement in optiondance",
	Run: func(cmd *cobra.Command, args []string) {
		deliveryType = strings.ToUpper(deliveryType)
		if deliveryType != core.DeliveryTypeCash && deliveryType != core.DeliveryTypePhysical {
			log.Panicf("wrong delivery type: %s", deliveryType)
		}
		switch act {
		case ActionRun:
			zap.L().Info("run settlement task via cmd ,settlement date:" + date)
			tx, err := mtg.SettlementTx(cmd.Context(), date, deliveryType)
			if err != nil {
				zap.L().Error("settlement error", zap.Error(err))
				return
			}
			zap.L().Info("settlement tx send", zap.Any("tx", tx))
			break
		case ActionGather:
			zap.L().Info("gather settlement info via cmd, settlement date:" + date)
			positionz := ProvidePositionService()
			gather, err := positionz.SettlementGather(cmd.Context(), "", date, deliveryType, false, time.Now())
			if err != nil {
				zap.L().Error("gather error", zap.Error(err))
				return
			}
			positionz.LogGatherResult(gather)
			break
		default:
			break
		}
	},
}

func init() {
	SettlementCmd.PersistentFlags().StringVarP(&date, "date", "d", "", "settlement date, eg: 1970-01-01")
	SettlementCmd.PersistentFlags().StringVarP(&act, "action", "a", "gather", "settlement action, gather or run")
	SettlementCmd.PersistentFlags().StringVar(&deliveryType, "deliveryType", "cash", "settlement deliveryType, cash or physical")
}
