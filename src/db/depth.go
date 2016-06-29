package db

import (
	"common"
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"logger"
	"time"
)

type Depth struct {
	gorm.Model

	Exchange  string
	Orderbook string `gorm:"type:longtext;"`
}

func SetDepth(exchange string, orderbook *common.OrderBook) (err error) {
	jsonOrderbook, err := json.Marshal(orderbook)
	if err != nil {
		logger.Errorln("Marshal fail", err)
		return
	}

	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var depth Depth
	depth.Exchange = exchange
	depth.Orderbook = string(jsonOrderbook)

	if err := db.Create(&depth).Error; err != nil {
		logger.Errorln("SetDepth err:", err)
		return err
	}

	return
}

func GetDepth(exchange string) (orderbook common.OrderBook, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var depth Depth

	if err = db.Where("exchange = ?", exchange).Last(&depth).Error; err != nil {
		logger.Errorln("GetDepth Last err:", err, exchange)
		return
	}

	now_time := time.Now()
	diff := now_time.Sub(depth.CreatedAt).Seconds()

	logger.Debugln("time:", diff, depth.CreatedAt.Format(time.RFC3339), now_time.Format(time.RFC3339))
	if diff > 15 {
		err = errors.New("last depth falls behind 15 seconds")
		logger.Errorln(exchange, err, diff, depth.CreatedAt.Format(time.RFC3339), now_time.Format(time.RFC3339))
		return
	}

	if err = json.Unmarshal([]byte(depth.Orderbook), &orderbook); err != nil {
		logger.Errorln("Unmarshal fail", err)
		return
	}

	logger.Infoln("GetDepth:", depth.ID, depth.CreatedAt)
	return
}
