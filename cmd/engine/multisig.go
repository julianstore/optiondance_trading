package main

import (
	"github.com/spf13/cobra"
	"log"
	"option-dance/core"
	"option-dance/service/mtg"
)

const (
	Unlock = "unlock"
)

var (
	multisigAction string
	spentBy        string
	tx             string
	currency       string
)

func init() {
	MultisigCmd.PersistentFlags().StringVarP(&spentBy, "spentby", "s", "", "spent by")
	MultisigCmd.PersistentFlags().StringVarP(&tx, "tx", "x", "", "signedTx")
	MultisigCmd.PersistentFlags().StringVarP(&currency, "asset", "a", "", "asset id")
	MultisigCmd.PersistentFlags().StringVarP(&multisigAction, "type", "t", "unlock", "multisig unlock")
}

var MultisigCmd = &cobra.Command{
	Use:   "multisig ",
	Short: "multisig utxo action",
	Run: func(cmd *cobra.Command, args []string) {
		switch multisigAction {
		case Unlock:
			log.Println("multisig unlock")
			assetId, _ := core.GetAssetIdByCurrency(currency)
			ctx := cmd.Context()
			if tx == "" {
				utxoStore := ProvideUtxoStore()
				utxo, err := utxoStore.FindSpentBy(ctx, assetId, spentBy)
				if err != nil {
					log.Panicln(err)
				}
				tx = utxo.SignedTx
			}
			err := mtg.UnlockMutisig(ctx, tx)
			if err != nil {
				log.Printf("unlock success")
			}
			break
		default:
			break
		}
	},
}
