/*
  trader API Engine
*/

package chbtc

import (
	"common"
	"config"
	"logger"
	"util"
)

type Chbtc struct {
	tradeAPI *ChbtcTrade
	name     string
}

func NewExchange(name, access_key, secret_key string) *Chbtc {
	w := new(Chbtc)
	w.tradeAPI = NewChbtcTrade(access_key, secret_key)
	return w
}

func (w Chbtc) GetTicker() (ticker common.Ticker, err error) {
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "btc"
	} else {
		symbol = "ltc"
		panic(-1)
	}

	chTicker, err := w.getTicker(symbol)
	if err != nil {
		return
	}

	ticker.Date = ""
	ticker.Ticker.Buy = util.ToFloat(chTicker.Ticker.Buy)
	ticker.Ticker.High = util.ToFloat(chTicker.Ticker.High)
	ticker.Ticker.Last = util.ToFloat(chTicker.Ticker.Last)
	ticker.Ticker.Low = util.ToFloat(chTicker.Ticker.Low)
	ticker.Ticker.Sell = util.ToFloat(chTicker.Ticker.Sell)
	ticker.Ticker.Vol = util.ToFloat(chTicker.Ticker.Vol)
	return
}

func (w Chbtc) GetDepth() (orderBook common.OrderBook, err error) {
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "btc"
	} else {
		symbol = "ltc"
		panic(-1)
	}

	return w.getDepth(symbol)
}

func (w Chbtc) Buy(price, btc string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	// symbol := config.Config["symbol"]
	return tradeAPI.doTrade("buy", price, btc)
}

func (w Chbtc) Sell(price, btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI

	// symbol := config.Config["symbol"]
	return tradeAPI.doTrade("sell", price, btc)
}

func (w Chbtc) BuyMarket(amount string) (buyId string, result string, err error) {
	panic("exchange does not support")
}

func (w Chbtc) SellMarket(amount string) (sellId string, result string, err error) {
	panic("exchange does not support")
}

func (w Chbtc) GetOrder(order_id string) (order common.Order, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "1"
	} else if symbol == "ltc_cny" {
		symbol = "0"
	} else {
		panic(-1)
	}

	return tradeAPI.Get_order(symbol, order_id)
}

func (w Chbtc) CancelOrder(order_id string) (err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "1"
	} else if symbol == "ltc_cny" {
		symbol = "0"
	} else {
		panic(-1)
	}
	return tradeAPI.Cancel_order(symbol, order_id)
}

func (w Chbtc) GetAccount() (account common.Account, err error) {
	tradeAPI := w.tradeAPI

	account, err = tradeAPI.GetAccount()
	if err != nil {
		logger.Debugln("Chbtc GetAccount failed")
		return
	}

	return
}
