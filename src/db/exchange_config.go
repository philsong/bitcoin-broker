package db

import (
	"logger"
	"trade_service"
)

func SetExchangeConfig(exchange_configs []*trade_service.ExchangeConfig) (err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	tx := db.Begin()

	var db_exchange_config trade_service.ExchangeConfig
	if err := tx.Delete(&db_exchange_config).Error; err != nil {
		logger.Errorln("_setExchangeConfig delete all err:", err)
		tx.Rollback()
		return err
	}

	for i := 0; i < len(exchange_configs); i++ {
		exchange_config := exchange_configs[i]
		if err := tx.Create(&exchange_config).Error; err != nil {
			logger.Errorln("_setExchangeConfig save err:", err)
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func GetExchangeConfig(exchange string) (exchange_config trade_service.ExchangeConfig, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err = db.Where("exchange = ?", exchange).First(&exchange_config).Error; err != nil {
		logger.Errorln("GetExchangeConfig Find err:", err, exchange)
		return
	}

	return
}

func GetExchangeConfigs() (exchange_configs []trade_service.ExchangeConfig, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err = db.Find(&exchange_configs).Error; err != nil {
		logger.Errorln("GetExchangeConfigs err:", err)
		return
	}

	return
}
