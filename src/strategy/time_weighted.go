/*
  trader  strategy
*/

package strategy

import (
	"db"
	"fmt"
	"github.com/jinzhu/gorm"
	"logger"
	"trade_service"
)

func processTimeWeighted(tx *gorm.DB, tradeOrder trade_service.TradeOrder) (err error) {
	exchange := tradeOrder.Exchange

	tradePrice := fmt.Sprintf("%0.2f", tradeOrder.EstimatePrice)
	tradeBTC := fmt.Sprintf("%0.4f", tradeOrder.EstimateBtc)
	tradeCNY := fmt.Sprintf("%0.2f", tradeOrder.EstimateCny)
	tradeType := tradeOrder.TradeType
	logger.Infoln("processTimeWeighted", exchange, tradeType, tradePrice, tradeBTC, tradeCNY)

	cur_unit_amount := unit_max_amount
	logger.Infoln("rm factor,cur_unit_amount:", cur_unit_amount)

	if tradeOrder.EstimateBtc < 2*cur_unit_amount {
		tradeOrder.OrderStatus = trade_service.OrderStatus_READY
	} else {
		tradeOrder.OrderStatus = trade_service.OrderStatus_SPLIT

		count := (int)(tradeOrder.EstimateBtc / cur_unit_amount)
		avg_price := tradeOrder.EstimateCny / tradeOrder.EstimateBtc

		logger.Infoln("processTimeWeighted", count, avg_price, cur_unit_amount)
		for i := 0; i < count; i++ {
			sub_tradeOrder := tradeOrder //延续/覆盖原有订单的信息
			//更新新信息
			btc_amount := 0.0
			if count == i+1 {
				btc_amount = tradeOrder.EstimateBtc - cur_unit_amount*float64(i)
			} else {
				btc_amount = cur_unit_amount
			}

			sub_tradeOrder.EstimatePrice = avg_price
			sub_tradeOrder.EstimateBtc = btc_amount
			sub_tradeOrder.EstimateCny = btc_amount * avg_price
			sub_tradeOrder.OrderStatus = trade_service.OrderStatus_READY
			logger.Infoln("InsertTradeOrder sub_tradeOrder begin", sub_tradeOrder)
			_, err = db.InsertTradeOrder(tx, &sub_tradeOrder)
			logger.Infoln("InsertTradeOrder sub_tradeOrder end", err, sub_tradeOrder)
			if err != nil {
				logger.Errorln("InsertTradeOrder sub_tradeOrder failed", err, sub_tradeOrder)
				return
			}
		}
	}

	logger.Infoln("processTimeWeighted UpdateTradeOrder begin:", tradeOrder)
	err = db.UpdateTradeOrder(tx, tradeOrder)
	if err != nil {
		logger.Errorln("processTimeWeighted UpdateTradeOrder failed", err, tradeOrder)
		return
	}
	logger.Infoln("processTimeWeighted UpdateTradeOrder end:", tradeOrder)

	return
}

func ProcessTimeWeighted(tx *gorm.DB) (err error) {
	tradeOrders, err := db.GetAllTradeOrders(tx, trade_service.OrderStatus_TIME_WEIGHTED)
	if err != nil {
		logger.Errorln("GetTradeOrders OrderStatus_TIME_WEIGHTED failed...")
		return
	}

	if len(tradeOrders) == 0 {
		return
	}

	logger.Infoln("processTimeWeighted tradeOrders", len(tradeOrders), tradeOrders)
	for key, _ := range tradeOrders {
		tradeOrder := tradeOrders[key]
		err = processTimeWeighted(tx, tradeOrder)
		if err != nil {
			return
		}
	}

	return
}
