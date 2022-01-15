package core

import (
	"fmt"
	"gorm.io/gorm"
	"option-dance/cmd/config"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator() Migrator {
	return Migrator{db: config.Db}
}

func (m *Migrator) Run() error {
	if err := m.db.AutoMigrate(
		Action{},
		OptionMarket{},
		Order{},
		Position{},
		Property{},
		RawTransaction{},
		Trade{},
		Transfer{},
		User{},
		UTXO{},
		WaitList{},
		MixinMessage{},
		MixinUser{},
		Message{},
		DeliveryPrice{},
		DbMarket{},
		DbMarketTracker{},
	); err != nil {
		return err
	}
	fmt.Println("migrate success")
	return nil
}
