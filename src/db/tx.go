package db

import (
	"github.com/jinzhu/gorm"
	"logger"
)

func TxBegin() (tx *gorm.DB, err error) {
	db_handle, err := GetDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	tx = db_handle.Begin()

	return tx, nil
}

func TxEnd(tx *gorm.DB, exception error) (err error) {
	if exception != nil {
		logger.Errorln(exception)
		tx.Rollback()
		return nil
	} else {
		tx.Commit()
		return nil
	}

	return nil
}

func TXWrapper(f func(tx *gorm.DB) error) {
	tx, err := TxBegin()
	if err != nil {
		logger.Errorln("TxBegin  failed", err)
		return
	}

	err = f(tx)

	err = TxEnd(tx, err)
	if err != nil {
		logger.Errorln("TxEnd  failed", err)
		return
	}

	return
}

func TXWrapperEx(f func(tx *gorm.DB, exchange string) error, exchange string) {
	tx, err := TxBegin()
	if err != nil {
		logger.Errorln("TxBegin  failed", err)
		return
	}

	err = f(tx, exchange)

	err = TxEnd(tx, err)
	if err != nil {
		logger.Errorln("TxEnd  failed", err)
		return
	}

	return
}
