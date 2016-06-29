package trade_server

import (
	"db"
	"logger"
	"strategy"
	"trade_service"
)

type TradeServiceHandler struct {
}

func (this *TradeServiceHandler) Ping() (err error) {
	return nil
}

func (this *TradeServiceHandler) ConfigKeys(exchange_configs []*trade_service.ExchangeConfig) (err error) {
	logger.Infoln("-->ConfigKeys begin:")

	if err = db.SetExchangeConfig(exchange_configs); err == nil {
		strategy.Update_tickers()
		strategy.Update_traders()
	}

	logger.Infoln("-->ConfigKeys end:", err)
	return
}

func (this *TradeServiceHandler) ConfigAmount(amount_config *trade_service.AmountConfig) (err error) {
	logger.Infoln("-->ConfigAmount begin:")
	err = db.SetAmountConfig(amount_config)
	logger.Infoln("-->ConfigAmount end:", err)
	return
}

func (this *TradeServiceHandler) CheckPrice(price float64, trade_type trade_service.TradeType) (err error) {
	if !strategy.Check_ticker_limit(price, trade_type) {
		tradeResult_ := trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_PRICE_NOT_SYNC
		logger.Infoln(price, trade_type, tradeResult_)
		return tradeResult_
	}

	return nil
}

// Parameters:
//  - Cny
func (this *TradeServiceHandler) Buy(buyOrder *trade_service.Trade) (err error) {
	logger.Infoln("-->Buy begin:", buyOrder)

	tradeResult_ := strategy.Buy(buyOrder)

	logger.Infoln("-->Buy end:", buyOrder, tradeResult_)
	return tradeResult_
}

// Parameters:
//  - Btc
func (this *TradeServiceHandler) Sell(sellOrder *trade_service.Trade) (err error) {
	logger.Infoln("-->Sell begin:", sellOrder)

	tradeResult_ := strategy.Sell(sellOrder)

	logger.Infoln("-->Sell end:", sellOrder, tradeResult_)
	return tradeResult_
}

func (this *TradeServiceHandler) GetAccount() (r []*trade_service.Account, err error) {
	logger.Infoln("-->GetAccount begin:")
	r, err = db.GetAccount()
	if r == nil {
		r = make([]*trade_service.Account, 0)
	}

	logger.Infoln("-->GetAccount end:", r, err)
	return
}

func (this *TradeServiceHandler) GetTicker() (r *trade_service.Ticker, err error) {
	// logger.Infoln("-->GetTicker begin:")
	r, err = db.GetTicker()
	// logger.Infoln("-->GetTicker end:", r, err)
	return
}

func (this *TradeServiceHandler) GetAlertOrders() (err error) {
	logger.Infoln("-->GetAlertOrders begin:")
	tradeOrders, err := db.GetAlertOrders()
	if err == nil && len(tradeOrders) > 0 {
		tradeResult_ := trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_EXIST_ERROR_ORDERS
		logger.Infoln("-->GetAlertOrders end:", tradeResult_)
		return tradeResult_
	}

	logger.Infoln("-->GetAlertOrders end:", err)
	return
}

func (this *TradeServiceHandler) GetExchangeStatus() (r *trade_service.ExchangeStatus, err error) {
	logger.Infoln("-->GetExchangeStatus begin:")

	r = trade_service.NewExchangeStatus()
	markets := strategy.GetUsableExchange(trade_service.TradeType_BUY.String(), true)
	if len(markets) == 0 {
		r.Canbuy = false
	} else {
		r.Canbuy = true
	}

	markets = strategy.GetUsableExchange(trade_service.TradeType_SELL.String(), true)
	if len(markets) == 0 {
		r.Cansell = false
	} else {
		r.Cansell = true
	}

	logger.Infoln("-->GetExchangeStatus end:", err, r)
	return
}
