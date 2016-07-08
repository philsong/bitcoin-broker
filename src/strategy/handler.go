/*
  trader  strategy
*/

package strategy

import (
	"db"
	"logger"
	"trade_service"
)

func GetOrderByClientID(client_id string) (orders []db.SiteOrder, err error) {
	tx, err := db.TxBegin()
	if err != nil {
		logger.Errorln("TxBegin  failed", err)
		return
	}

	orders, err = db.GetOrderByClientID(tx, client_id)

	err = db.TxEnd(tx, err)
	if err != nil {
		logger.Errorln("TxEnd  failed", err)
		return
	}

	return
}

func check_client_order(order *trade_service.Trade) (tradeResult_ *trade_service.TradeException, exist bool) {
	orders, err := GetOrderByClientID(order.GetClientID())
	if err != nil {
		tradeResult_ := trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_INTERNAL_ERROR
		logger.Infoln(order, tradeResult_)
		return tradeResult_, false
	}

	if len(orders) > 0 {
		exist = true
		logger.Infoln("order already exist:", order, len(orders), orders)
		return nil, true
	}

	return nil, false
}

func handleTrade(trade *trade_service.Trade, tradeType trade_service.TradeType) (tradeResult_ *trade_service.TradeException) {
	siteOrder := &db.SiteOrder{}
	siteOrder.ClientID = trade.GetClientID()
	siteOrder.Amount = trade.GetAmount()
	siteOrder.Price = trade.GetPrice()
	siteOrder.TradeType = tradeType

	logger.Infoln("handleTrade begin", siteOrder)

	markets := GetUsableExchange(tradeType.String(), true)
	if len(markets) == 0 {
		tradeResult_ = trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_NO_USABLE_FUND

		logger.Errorln(tradeType.String(), tradeResult_)
		return
	}

	tradeOrders, tradeResult_ := estimateOrder(siteOrder, markets)
	if tradeResult_ != nil {
		logger.Infoln(siteOrder, tradeResult_)
		return
	}

	// if !check_price_limit(siteOrder) {
	// 	tradeResult_ = trade_service.NewTradeException()
	// 	tradeResult_.Reason = trade_service.EX_PRICE_OUT_OF_SCOPE
	// 	logger.Infoln(siteOrder, tradeResult_)
	// 	return
	// }

	err := PushOrder(siteOrder, tradeOrders)
	if err != nil {
		tradeResult_ = trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_INTERNAL_ERROR
	}

	logger.Infoln("handleTrade end", tradeResult_, siteOrder)

	return
}

func Buy(buyOrder *trade_service.Trade) *trade_service.TradeException {
	tradeResult_, exist := check_client_order(buyOrder)
	if tradeResult_ != nil {
		return tradeResult_
	}

	if exist {
		return nil
	}

	// if !Check_ticker_limit(buyOrder.Price, trade_service.TradeType_BUY) {
	// 	tradeResult_ := trade_service.NewTradeException()
	// 	tradeResult_.Reason = trade_service.EX_PRICE_NOT_SYNC
	// 	logger.Infoln(buyOrder, tradeResult_)
	// 	return tradeResult_
	// }

	ret := is_limit_buy(buyOrder.Amount)
	if ret {
		tradeResult_ := trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_TRADE_QUEUE_FULL
		logger.Infoln(tradeResult_)
		return tradeResult_
	}

	tradeResult_ = handleTrade(buyOrder, trade_service.TradeType_BUY)
	if tradeResult_ != nil {
		decr_buy(buyOrder.Amount)
	}

	return tradeResult_
}

func Sell(sellOrder *trade_service.Trade) *trade_service.TradeException {
	tradeResult_, exist := check_client_order(sellOrder)
	if tradeResult_ != nil {
		return tradeResult_
	}

	if exist {
		return nil
	}

	// if !Check_ticker_limit(sellOrder.Price, trade_service.TradeType_SELL) {
	// 	tradeResult_ := trade_service.NewTradeException()
	// 	tradeResult_.Reason = trade_service.EX_PRICE_NOT_SYNC
	// 	logger.Infoln(sellOrder, tradeResult_)
	// 	return tradeResult_
	// }

	ret := is_limit_sell(sellOrder.Amount)
	if ret {
		tradeResult_ := trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_TRADE_QUEUE_FULL
		logger.Infoln(tradeResult_)
		return tradeResult_
	}

	tradeResult_ = handleTrade(sellOrder, trade_service.TradeType_SELL)
	if tradeResult_ != nil {
		decr_sell(sellOrder.Amount)
	}

	return tradeResult_
}
