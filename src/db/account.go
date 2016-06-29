package db

import (
	"github.com/jinzhu/gorm"
	"logger"
	"time"
	"trade_service"
)

type Account struct {
	gorm.Model
	trade_service.Account
}

func SetAccount(account *trade_service.Account) (err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var db_account Account
	if err := db.Where("exchange = ?", account.Exchange).First(&db_account).Error; err != nil {
		logger.Errorln("not found:", err)
		db_account.Exchange = account.Exchange
		db_account.AvailableCny = account.AvailableCny
		db_account.AvailableBtc = account.AvailableBtc
		db_account.FrozenCny = account.FrozenCny
		db_account.FrozenBtc = account.FrozenBtc

		if err := db.Create(&db_account).Error; err != nil {
			logger.Errorln(err)
			return err
		}
	} else {
		db_account.CreatedAt = time.Now()
		db_account.AvailableCny = account.AvailableCny
		db_account.AvailableBtc = account.AvailableBtc
		db_account.FrozenCny = account.FrozenCny
		db_account.FrozenBtc = account.FrozenBtc
		if err := db.Save(&db_account).Error; err != nil {
			logger.Errorln(err)
			return err
		}
	}

	logger.Debugln("set account ok", account)
	return
}

func GetAccount() (accounts []*trade_service.Account, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var db_accounts []Account

	if err = db.Order("id").Find(&db_accounts).Error; err != nil {
		logger.Errorln("GetAccount err:", err)
		return
	}

	for _, db_account := range db_accounts {
		account := new(trade_service.Account)
		account.Exchange = db_account.Exchange
		account.AvailableCny = db_account.AvailableCny
		account.AvailableBtc = db_account.AvailableBtc
		account.FrozenCny = db_account.FrozenCny
		account.FrozenBtc = db_account.FrozenBtc
		account.PauseTrade = db_account.PauseTrade

		accounts = append(accounts, account)
	}

	return
}
