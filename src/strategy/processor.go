/*
  trader  strategy
*/

package strategy

import (
	"common"
	"db"
	"fmt"
	"github.com/jinzhu/gorm"
	"logger"
	"time"
	"trade_service"
)

func processReady(tx *gorm.DB, tradeOrder trade_service.TradeOrder) (err error) {
	if tradeOrder.OrderStatus != trade_service.OrderStatus_READY {
		return
	}

	exchange := tradeOrder.Exchange

	tradePrice := fmt.Sprintf("%0.2f", tradeOrder.EstimatePrice)
	tradeBTC := fmt.Sprintf("%0.4f", tradeOrder.EstimateBtc)
	tradeCNY := fmt.Sprintf("%0.2f", tradeOrder.EstimateCny)
	tradeType := tradeOrder.TradeType

	if use_time_weighted_algorithm {
		orderbook, _err := db.GetDepth(exchange)
		if _err != nil {
			logger.Errorln("db.GetDepth err", exchange, err)
			return
		}

		// okcoin except special progress
		size := len(orderbook.Bids)
		if size < 1 || (size > 0 && orderbook.Bids[0].Price < 0.000001) {
			logger.Errorln("exception orderbook:", exchange, orderbook)
			return
		}
		// logger.Infoln(orderbook.Bids)
		// logger.Infoln(orderbook.Asks)
		// logger.Infoln(orderbook.Bids[0].Price, orderbook.Asks[size-1])

		market_price := (orderbook.Bids[0].Price + orderbook.Asks[size-1].Price) * 0.5
		tradePrice = fmt.Sprintf("%0.2f", market_price)
		if tradeOrder.TryTimes > 0 {
			if tradeType == trade_service.TradeType_BUY {
				tradePrice = fmt.Sprintf("%0.2f", market_price+20)
			} else {
				tradePrice = fmt.Sprintf("%0.2f", market_price-20)
			}
		}

		if exchange == "haobtc" {
			haobtc_market_price := 0
			if tradeType == trade_service.TradeType_BUY {
				haobtc_market_price = int(orderbook.Asks[size-1].Price + 1)
			} else {
				haobtc_market_price = int(orderbook.Bids[0].Price - 1)
			}
			tradePrice = fmt.Sprintf("%d", haobtc_market_price)

			if tradeOrder.TryTimes > 0 {
				if tradeType == trade_service.TradeType_BUY {
					tradePrice = fmt.Sprintf("%d", haobtc_market_price+20)
				} else {
					tradePrice = fmt.Sprintf("%d", haobtc_market_price-20)
				}
			}
		}
	}

	logger.Infoln("processReady", exchange, tradePrice, tradeBTC, tradeCNY)

	tradeAPI, err := GetExchange(exchange)
	if err != nil {
		logger.Errorln(err)
		tradeOrder.Info = err.Error()
		return
	}

	if tradeOrder.EstimateBtc < unit_min_amount || tradeOrder.EstimateCny < unit_min_amount*tradeOrder.EstimatePrice {
		tradeOrder.Memo += fmt.Sprintf("忽略BTC数量%.3f/%.3f小于%.3f的订单", tradeOrder.EstimateBtc, tradeOrder.EstimateCny, unit_min_amount)
		tradeOrder.OrderID = "-1"
		tradeOrder.OrderStatus = trade_service.OrderStatus_ORDERED
	} else {
		var orderID string
		var result string

		funds_log := fmt.Sprintf("tradeAPI begin: id=%d exchange=%s tradePrice=%s tradeBTC=%s tradeCNY=%s\n",
			tradeOrder.ID,
			exchange,
			tradePrice,
			tradeBTC,
			tradeCNY)
		logger.Infoln(funds_log)

		if tradeOrder.TryTimes > 0 {
			// if tradeType == trade_service.TradeType_BUY {
			// 	orderID, result, err = tradeAPI.BuyMarket(tradeCNY)
			// } else {
			// 	orderID, result, err = tradeAPI.SellMarket(tradeBTC)
			// }
			if tradeType == trade_service.TradeType_BUY {
				orderID, result, err = tradeAPI.Buy(tradePrice, tradeBTC)
			} else {
				orderID, result, err = tradeAPI.Sell(tradePrice, tradeBTC)
			}

			tradeOrder.Memo += fmt.Sprintf("market:%s;", tradePrice)
		} else {
			if tradeType == trade_service.TradeType_BUY {
				orderID, result, err = tradeAPI.Buy(tradePrice, tradeBTC)
			} else {
				orderID, result, err = tradeAPI.Sell(tradePrice, tradeBTC)
			}
			tradeOrder.Memo += fmt.Sprintf("limit:%s;", tradePrice)
		}

		if err != nil {
			logger.Infoln("tradeAPI error:", err, tradeOrder.ID, orderID, result)
			tradeOrder.Info = err.Error()
			tradeOrder.OrderStatus = trade_service.OrderStatus_READY
		} else {
			funds_log = fmt.Sprintf("tradeAPI end: id=%d exchange=%s orderID=%s tradePrice=%s tradeBTC=%s tradeCNY=%s result=%s\n",
				tradeOrder.ID,
				exchange,
				orderID,
				tradePrice,
				tradeBTC,
				tradeCNY,
				result)
			logger.Infoln(funds_log)

			if result != "" {
				tradeOrder.Info = result
				logger.Infoln("trade action failed, change to nextExchange", result)

				newExhanges := nextExchange(tradeOrder.ID, tradeOrder.TradeType.String(), exchange)
				if len(newExhanges) != 0 {
					tradeOrder.Memo += fmt.Sprintf("%s->%s;", tradeOrder.Exchange, newExhanges[0])
					tradeOrder.Exchange = newExhanges[0]
					tradeOrder.OrderStatus = trade_service.OrderStatus_READY
				} else {
					tradeOrder.OrderStatus = trade_service.OrderStatus_ERROR
				}
			} else {
				tradeOrder.Info = ""
				tradeOrder.OrderID = orderID
				tradeOrder.OrderStatus = trade_service.OrderStatus_ORDERED
				tradeOrder.TryTimes++
				tradeOrder.UpdateAt = time.Now().Format(time.RFC3339)
			}
		}
	}

	logger.Infoln("UpdateTradeOrder begin:", tradeOrder)
	err = db.UpdateTradeOrder(tx, tradeOrder)
	if err != nil {
		logger.Errorln("progressReady UpdateTradeOrder failed", err, tradeOrder)
		return
	}
	logger.Infoln("UpdateTradeOrder end:", tradeOrder)

	return
}

func processOrdered(tx *gorm.DB, tradeOrder trade_service.TradeOrder) (err error) {
	exchange := tradeOrder.Exchange

	tradeAPI, err := GetExchange(exchange)
	if err != nil {
		logger.Errorln(err)
		return
	}

	tradeType := tradeOrder.TradeType

	if tradeOrder.OrderID == "-1" { //ignore MINOR AMOUNT ORDER
		tradeOrder.PriceMargin = 0.0
		tradeOrder.DealPrice = tradeOrder.Price
		tradeOrder.DealBtc = tradeOrder.EstimateBtc
		tradeOrder.DealCny = tradeOrder.EstimateCny

		tradeOrder.OrderStatus = trade_service.OrderStatus_CANCELED
	} else if tradeOrder.OrderID == "" {
		tradeOrder.OrderStatus = trade_service.OrderStatus_ERROR //EXCEPT
	} else {
		funds_log := fmt.Sprintf("GetOrder begin: id=%d exchange=%s\n",
			tradeOrder.ID,
			tradeOrder.Exchange)
		logger.Infoln(funds_log)
		logger.Infoln("GetOrder begin ex:", tradeOrder)

		order, result, err := tradeAPI.GetOrder(tradeOrder.OrderID)
		if err != nil {
			logger.Infoln("GetOrder error:", err, order, result, tradeOrder)

			logger.Errorln(err)
			tradeOrder.Info = err.Error()
			tradeOrder.OrderStatus = trade_service.OrderStatus_ORDERED //GO ON CHECK IN NEXT CYCLE
		} else {
			funds_log := fmt.Sprintf("GetOrder end: id=%d exchange=%s\n",
				tradeOrder.ID,
				tradeOrder.Exchange)
			logger.Infoln(funds_log)

			tradeOrder.Info = result
			tradeOrder.DealCny = order.Deal_amount * order.Price
			tradeOrder.DealBtc = order.Deal_amount
			tradeOrder.DealPrice = order.Price

			logger.Infoln("GetOrder end ex:", order, result, tradeOrder)

			//MAP:CONVERT STATUS
			if order.Status == common.ORDER_STATE_SUCCESS {
				tradeOrder.OrderStatus = trade_service.OrderStatus_SUCCESS

				if tradeType == trade_service.TradeType_BUY {
					tradeOrder.PriceMargin = tradeOrder.Price - tradeOrder.DealPrice
					decr_buy(tradeOrder.DealCny)
				} else {
					tradeOrder.PriceMargin = tradeOrder.DealPrice - tradeOrder.Price
					decr_sell(tradeOrder.DealBtc)
				}
			} else if order.Status == common.ORDER_STATE_ERROR {
				tradeOrder.OrderStatus = trade_service.OrderStatus_ERROR
			} else if order.Status == common.ORDER_STATE_CANCELED {
				tradeOrder.OrderStatus = trade_service.OrderStatus_CANCELED
				if tradeType == trade_service.TradeType_BUY {
					if tradeOrder.DealCny > MeasurementError {
						tradeOrder.PriceMargin = tradeOrder.Price - tradeOrder.DealPrice
					}
					decr_buy(tradeOrder.DealCny)
				} else {
					if tradeOrder.DealBtc > MeasurementError {
						tradeOrder.PriceMargin = tradeOrder.DealPrice - tradeOrder.Price
					}
					decr_sell(tradeOrder.DealBtc)
				}

				//fixed: 注意：坑！ok的买单，amount返回为空
				if order.Amount < 0.000001 &&
					tradeOrder.TradeType == trade_service.TradeType_BUY &&
					(tradeOrder.Exchange == "okcoin" || tradeOrder.Exchange == "haobtc") {
					order.Amount = tradeOrder.EstimateBtc
				}

				left_btc := order.Amount - order.Deal_amount
				var left_amount float64
				if tradeType == trade_service.TradeType_BUY {
					left_amount = tradeOrder.EstimateCny - tradeOrder.DealCny
				} else {
					left_amount = left_btc
				}
				logger.Infoln("left to progress:", left_btc, left_amount)

				if (tradeType == trade_service.TradeType_BUY && left_btc < unit_min_amount) ||
					(tradeType == trade_service.TradeType_SELL && left_btc < unit_min_amount) {
					tradeOrder.Memo += fmt.Sprintf("忽略BTC数量%.3f小于%.3f的订单", left_btc, unit_min_amount)
					tradeOrder.PriceMargin = 0.0
					tradeOrder.OrderStatus = trade_service.OrderStatus_SUCCESS

					tradeOrder.DealBtc = tradeOrder.EstimateBtc
					tradeOrder.DealCny = tradeOrder.EstimateCny
				} else {
					newTradeOrder := tradeOrder //延续/覆盖原有订单的信息
					//更新新信息
					left_cny := 0.0
					if tradeType == trade_service.TradeType_BUY {
						left_cny = tradeOrder.EstimateCny - tradeOrder.DealCny
					} else {
						left_cny = left_btc * tradeOrder.Price
					}

					newTradeOrder.Memo += fmt.Sprintf("re-order:%d;", tradeOrder.ID)
					newTradeOrder.EstimateBtc = left_btc
					newTradeOrder.EstimateCny = left_cny
					newTradeOrder.OrderStatus = trade_service.OrderStatus_READY
					newTradeOrder.TryTimes++
					newTradeOrder.Info = ""
					newTradeOrder.OrderID = ""
					logger.Infoln("re-order newTradeOrder as market order begin", newTradeOrder)
					_, err = db.InsertTradeOrder(tx, &newTradeOrder)
					if err != nil {
						logger.Errorln("re-order newTradeOrder failed", err, newTradeOrder)
						return err
					} else {
						logger.Infoln("re-order newTradeOrder end", newTradeOrder)
					}
				}
			} else {
				tradeOrder.OrderStatus = trade_service.OrderStatus_ORDERED

			}
		}
	}

	logger.Infoln("UpdateTradeOrder begin:", tradeOrder)
	err = db.UpdateTradeOrder(tx, tradeOrder)
	if err != nil {
		logger.Errorln("processOrdered UpdateTradeOrder failed", err, tradeOrder)
		return
	}
	logger.Infoln("UpdateTradeOrder end:", tradeOrder)

	return
}

func processTimeout(tx *gorm.DB, tradeOrder trade_service.TradeOrder) (err error) {
	exchange := tradeOrder.Exchange

	tradeAPI, err := GetExchange(exchange)
	if err != nil {
		logger.Errorln(err)
		return
	}

	logger.Infoln("processTimeout:CancelOrder", tradeOrder.OrderID)
	err = tradeAPI.CancelOrder(tradeOrder.OrderID)
	if err != nil {
		logger.Infoln("processTimeout:CancelOrder error", tradeOrder.OrderID, err)
		tradeOrder.Info = err.Error()
		err = db.UpdateTradeOrder(tx, tradeOrder)
		if err != nil {
			tradeOrder.Info = err.Error()
			logger.Errorln("processTimeout:UpdateTradeOrder failed", err, tradeOrder)
			return
		}
		//even cancel err, we go on check order for next step
	} else {
		logger.Infoln("processTimeout:CancelOrder end", tradeOrder.OrderID)
	}

	logger.Infoln("processTimeout:processOrdered", tradeOrder)
	err = processOrdered(tx, tradeOrder)
	logger.Infoln("processTimeout:processOrdered result:", err)
	return
}

func ProcessReady(tx *gorm.DB, exchange string) (err error) {
	if is_queueing(tx, exchange) {
		return
	}

	tradeOrders, err := db.GetFirstTradeOrder(tx, exchange, trade_service.OrderStatus_READY)
	if err != nil {
		logger.Errorln("GetFirstTradeOrder READY failed...")
		return
	}

	if len(tradeOrders) == 0 {
		return
	}

	logger.Infoln("processReady tradeOrders", tradeOrders)
	for key, _ := range tradeOrders {
		tradeOrder := tradeOrders[key]
		err = processReady(tx, tradeOrder)
		if err != nil {
			return
		}
	}

	return
}

func ProcessOrdered(tx *gorm.DB, exchange string) (err error) {
	tradeOrders, err := db.GetTradeOrders(tx, exchange, trade_service.OrderStatus_ORDERED)
	if err != nil {
		logger.Errorln("GetTradeOrders ORDERED failed...")
		return
	}

	if len(tradeOrders) == 0 {
		return
	}

	logger.Infoln("processOrdered tradeOrders", tradeOrders)
	for key, _ := range tradeOrders {
		tradeOrder := tradeOrders[key]
		err = processOrdered(tx, tradeOrder)
		if err != nil {
			return
		}
	}

	return
}

func ProcessTimeout(tx *gorm.DB, exchange string) (err error) {
	tradeOrders, err := db.GetTimeoutOrder(tx, exchange)
	if err != nil {
		logger.Errorln("GetTimeoutOrder failed...")
		return
	}

	if len(tradeOrders) == 0 {
		return
	}

	logger.Infoln("processTimeout tradeOrders", tradeOrders)
	//cancel ,query, reorder
	for key, _ := range tradeOrders {
		tradeOrder := tradeOrders[key]
		err = processTimeout(tx, tradeOrder)
		if err != nil {
			return
		}
	}

	return
}

func is_queueing(tx *gorm.DB, exchange string) bool {
	tradeOrders, err := db.GetTradeOrders(tx, exchange, trade_service.OrderStatus_ORDERED)
	if err != nil {
		logger.Errorln("GetTradeOrders ORDERED failed...")
		return true
	}

	if len(tradeOrders) > 0 {
		logger.Infoln("order queueing: ", exchange, len(tradeOrders))
		return true
	}

	return false
}
