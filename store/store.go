package store

import (
	"errors"
	"gorm.io/gorm"
)

func CheckErr(err error) error {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}
