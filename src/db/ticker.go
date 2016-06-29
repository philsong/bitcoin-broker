package db

import (
	"github.com/jinzhu/gorm"
	"logger"
	"trade_service"
)

type Ticker struct {
	gorm.Model
	Ask float64 `gorm:"type:decimal(65,2);"`
	Bid float64 `gorm:"type:decimal(65,2);"`
}

var g_cacheTicker *trade_service.Ticker

func SetTicker(ticker *trade_service.Ticker) (err error) {
	g_cacheTicker = ticker

	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var dbTicker Ticker
	dbTicker.Bid = ticker.GetBid()
	dbTicker.Ask = ticker.GetAsk()

	if err := db.Save(&dbTicker).Error; err != nil {
		logger.Errorln("SetTicker err:", err)
		return err
	}

	logger.Infoln("SetTicker ok", ticker)
	return
}

func GetTicker() (ticker *trade_service.Ticker, err error) {
	if g_cacheTicker != nil {
		return g_cacheTicker, nil
	}

	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var dbTicker Ticker

	if err = db.Last(&dbTicker).Error; err != nil {
		logger.Errorln("GetTicker Last err:", err)
		return
	}

	ticker = trade_service.NewTicker()
	ticker.Ask = dbTicker.Ask
	ticker.Bid = dbTicker.Bid

	logger.Infoln("GetTicker ok", ticker)
	g_cacheTicker = ticker

	return
}

func GetNTickers(ticker_compare_count int) (tickers []Ticker, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err = db.Order("id desc").Limit(ticker_compare_count).Find(&tickers).Error; err != nil {
		logger.Errorln("GetNTicker err:", err)
		return
	}

	logger.Infoln("GetNTickers ok", tickers)

	return
}
