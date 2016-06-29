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

func ProcessMatchTx(tx *gorm.DB) (err error) {
	tradeOrders, err := db.GetAllTradeOrders(tx, trade_service.OrderStatus_READY)
	if err != nil {
		logger.Errorln("GetTradeOrders READY failed...")
		return
	}

	if len(tradeOrders) < 2 {
		return
	}

	var buy_queue, sell_queue []trade_service.TradeOrder

	for i := 0; i < len(tradeOrders); i++ {
		tradeOrder := tradeOrders[i]
		tradeType := tradeOrder.TradeType
		if tradeType == trade_service.TradeType_BUY {
			buy_queue = append(buy_queue, tradeOrder)
		} else {
			sell_queue = append(sell_queue, tradeOrder)
		}
	}

	if len(buy_queue) == 0 || len(sell_queue) == 0 {
		return
	}

	logger.Infoln("ProcessMatch tradeOrders:")
	logger.Infoln(len(buy_queue), buy_queue)
	logger.Infoln(len(sell_queue), sell_queue)

	for i := 0; i < len(buy_queue); i++ {
		buy_tradeOrder := buy_queue[i]
		for j := 0; j < len(sell_queue); j++ {
			sell_tradeOrder := sell_queue[j]
			if buy_tradeOrder.EstimateBtc == sell_tradeOrder.EstimateBtc {
				logger.Infoln("ProcessMatch match:", buy_tradeOrder, sell_tradeOrder)
				sell_queue = append(sell_queue[:j], sell_queue[j+1:]...)
				j--

				buy_tradeOrder.OrderStatus = trade_service.OrderStatus_MATCH
				buy_tradeOrder.MatchID = sell_tradeOrder.ID
				buy_tradeOrder.Memo += fmt.Sprintf("match %d:%d;", buy_tradeOrder.ID, sell_tradeOrder.ID)
				buy_tradeOrder.DealPrice = sell_tradeOrder.Price //only one side to match the pair price.
				buy_tradeOrder.DealBtc = buy_tradeOrder.EstimateBtc
				buy_tradeOrder.DealCny = buy_tradeOrder.EstimateCny
				buy_tradeOrder.PriceMargin = buy_tradeOrder.Price - buy_tradeOrder.DealPrice

				sell_tradeOrder.OrderStatus = trade_service.OrderStatus_MATCH
				sell_tradeOrder.MatchID = buy_tradeOrder.ID
				sell_tradeOrder.Memo += fmt.Sprintf("match %d:%d;", sell_tradeOrder.ID, buy_tradeOrder.ID)
				sell_tradeOrder.DealPrice = sell_tradeOrder.Price
				sell_tradeOrder.DealBtc = sell_tradeOrder.EstimateBtc
				sell_tradeOrder.DealCny = sell_tradeOrder.EstimateCny

				decr_buy(buy_tradeOrder.DealCny)
				decr_sell(sell_tradeOrder.DealBtc)

				err = db.UpdateTradeOrder(tx, buy_tradeOrder)
				if err != nil {
					logger.Errorln("ProcessMatch UpdateTradeOrder buy failed", err, buy_tradeOrder)
					return
				}

				err = db.UpdateTradeOrder(tx, sell_tradeOrder)
				if err != nil {
					logger.Errorln("ProcessMatch UpdateTradeOrder sell failed", err, sell_tradeOrder)
					return
				}

				break
			}
		}
	}

	return
}
