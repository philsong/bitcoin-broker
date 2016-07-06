/*
  trader  strategy
*/

package strategy

import (
	"db"
	"fmt"
	"logger"
	"trade_service"
)

func _queryFund(exchange string) (err error) {
	logger.Debugln("_______queryFund begin:", exchange)
	tradeAPI, err := GetExchange(exchange)
	if err != nil {
		logger.Errorln("queryFund err:", err, exchange)
		return
	}

	account, err := tradeAPI.GetAccount()
	if err != nil {
		logger.Errorln("queryFund err:", err, exchange)
		return
	}

	var dbFund trade_service.Account
	dbFund.Exchange = exchange
	dbFund.AvailableCny = account.Available_cny
	dbFund.AvailableBtc = account.Available_btc
	dbFund.FrozenCny = account.Frozen_cny
	dbFund.FrozenBtc = account.Frozen_btc

	fundExchages[exchange] = dbFund
	funds_log := fmt.Sprintf("_queryFund: exchange=%s cny=%f btc=%f Frozen_cny=%f Frozen_btc=%f\n",
		exchange,
		dbFund.AvailableCny,
		dbFund.AvailableBtc,
		dbFund.FrozenCny,
		dbFund.FrozenBtc)
	logger.Debugln(funds_log)

	err = db.SetAccount(&dbFund)
	if err != nil {
		logger.Errorln("queryFund:", err, exchange)
		return
	}

	logger.Debugln("_______queryFund end:", exchange)
	return
}

func _queryDepth(exchange string) (err error) {
	logger.Debugln("_______queryDepth begin:", exchange)
	tradeAPI, err := GetExchange(exchange)
	if err != nil {
		logger.Errorln("_queryDepth:", err, exchange)
		return
	}

	orderbook, err := tradeAPI.GetDepth()
	if err != nil {
		logger.Errorln("_queryDepth:", err, exchange)
		return
	}

	err = db.SetDepth(exchange, &orderbook)
	if err != nil {
		logger.Errorln("_queryDepth:", err, exchange)
		return
	}

	logger.Debugln("_______queryDepth end:", exchange)
	return
}

func QueryTicker() {
	logger.Debugln("QueryTicker begin")
	btc_threshold := 30.0
	amount_config, err := db.GetAmountConfig()
	if err != nil {
		logger.Errorln(err)
		// return
	} else {
		btc_threshold = amount_config.MaxBtc
	}

	ticker := trade_service.NewTicker()
	var order db.SiteOrder
	order.ID = -1
	{
		order.Amount = btc_threshold + get_current_buy_total()
		order.TradeType = trade_service.TradeType_BUY

		markets := GetUsableExchange(order.TradeType.String(), true)
		if len(markets) == 0 {
			logger.Errorln("QueryTicker: no used market, use all:", order.TradeType.String())
			markets = GetUsableExchange(order.TradeType.String(), false)
		}

		_, err := estimateOrder(&order, markets)
		if err != nil {
			logger.Errorln("QueryTicker: estimateOrder:", err)
			return
		}

		ticker.Ask = order.EstimatePrice
	}

	{
		order.Amount = btc_threshold + get_current_sell_total()
		order.TradeType = trade_service.TradeType_SELL

		markets := GetUsableExchange(order.TradeType.String(), true)
		if len(markets) == 0 {
			logger.Errorln("QueryTicker: no used market, use all", order.TradeType.String())
			markets = GetUsableExchange(order.TradeType.String(), false)
		}

		_, err := estimateOrder(&order, markets)
		if err != nil {
			logger.Errorln(err)
			return
		}

		ticker.Bid = order.EstimatePrice
	}

	logger.Debugln(ticker.Ask, ticker.Bid)
	if ticker.Ask < ticker.Bid {
		logger.Infoln("QueryTicker adjust begin:", ticker.Ask, ticker.Bid)
		mid_price := (ticker.Ask + ticker.Bid) * 0.5
		ticker.Ask = mid_price + 0.01
		ticker.Bid = mid_price - 0.01
		logger.Infoln("QueryTicker adjust end:", ticker.Ask, ticker.Bid, mid_price)
	}

	db.SetTicker(ticker)

	logger.Debugln("QueryTicker end")

	return
}

func _queryExchangeData(exchange string) (err error) {
	_queryFund(exchange)
	_queryDepth(exchange)

	return
}
