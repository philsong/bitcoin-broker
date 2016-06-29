package db

import (
	"github.com/jinzhu/gorm"
	"logger"
	"time"
	"trade_service"
)

type AmountConfig struct {
	gorm.Model
	trade_service.AmountConfig
}

var g_amount_config *trade_service.AmountConfig

func SetAmountConfig(amount_config *trade_service.AmountConfig) (err error) {
	g_amount_config = amount_config
	return setAmountConfig(amount_config)
}

func GetAmountConfig() (amount_config *trade_service.AmountConfig, err error) {
	if g_amount_config != nil {
		return g_amount_config, nil
	}

	amount_config, err = getAmountConfig()
	if err == nil {
		g_amount_config = amount_config
	}

	return
}

func setAmountConfig(amount_config *trade_service.AmountConfig) (err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	tx := db.Begin()

	var amountConfig AmountConfig

	if err := tx.First(&amountConfig).Error; err != nil {
		logger.Errorln("setAmountConfig amount record does not exist:", err)
		// return err
	}

	amountConfig.CreatedAt = time.Now()
	amountConfig.MaxCny = amount_config.MaxCny
	amountConfig.MaxBtc = amount_config.MaxBtc

	if err := tx.Save(&amountConfig).Error; err != nil {
		tx.Rollback()
		logger.Errorln(err)
		return err
	}

	tx.Commit()

	logger.Infoln("set amountConfig ok", amountConfig)

	return
}

func getAmountConfig() (amount_config *trade_service.AmountConfig, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var amountConfig AmountConfig

	if err = db.First(&amountConfig).Error; err != nil {
		logger.Errorln("getAmountConfig amount record does not exist:", err)
		return
	}

	amount_config = trade_service.NewAmountConfig()
	amount_config.MaxCny = amountConfig.MaxCny
	amount_config.MaxBtc = amountConfig.MaxBtc

	logger.Infoln("get amountConfig ok", amount_config)

	return
}
