/*
  trader  strategy
*/

package strategy

import (
	"db"
	"github.com/jinzhu/gorm"
	"logger"
	"strings"
	"trade_service"
)

func pushOrder(tx *gorm.DB, siteOrder *db.SiteOrder, tradeOrders map[string]*trade_service.TradeOrder) (err error) {
	logger.Infoln("InsertOrder begin", siteOrder)
	_, err = db.InsertOrder(tx, siteOrder)
	logger.Infoln("InsertOrder end", siteOrder)
	//Error 1062: Duplicate entry
	if err != nil && strings.Contains(err.Error(), "1062") {
		logger.Errorln("InsertOrder  Duplicate:", err, siteOrder)
		//here do'not need reponse NewTradeException tradeResult_
		return nil
	}

	if err != nil {
		logger.Errorln("InsertOrder  failed:", err, siteOrder)
		return
	}

	logger.Infoln("InsertTradeOrders begin", len(tradeOrders), tradeOrders)
	for exchange, _ := range tradeOrders {
		tradeOrder := tradeOrders[exchange]
		tradeOrder.Price = siteOrder.Price
		tradeOrder.SiteOrderID = siteOrder.ID
		logger.Infoln("InsertTradeOrder begin", tradeOrder)
		_, err = db.InsertTradeOrder(tx, tradeOrder)
		logger.Infoln("InsertTradeOrder end", err, tradeOrder)
		if err != nil {
			logger.Errorln("InsertTradeOrder failed", err, tradeOrder)
			return
		}
	}
	logger.Infoln("InsertTradeOrders end", len(tradeOrders), tradeOrders)

	return
}

func PushOrder(siteOrder *db.SiteOrder, tradeOrders map[string]*trade_service.TradeOrder) (err error) {
	tx, err := db.TxBegin()
	if err != nil {
		logger.Errorln("TxBegin  failed", err)
		return
	}

	err = pushOrder(tx, siteOrder, tradeOrders)

	err = db.TxEnd(tx, err)
	if err != nil {
		logger.Errorln("TxEnd  failed", err)
		return
	}

	return
}
