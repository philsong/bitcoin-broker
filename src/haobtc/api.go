/*
  trader API Engine
*/

package haobtc

import (
	"common"
	"config"
	"fmt"
	"logger"
	"util"
)

type Haobtc struct {
	tradeAPI *HaobtcTrade
	name     string
}

func NewExchange(name, partner, secret_key string) *Haobtc {
	w := new(Haobtc)
	w.name = name
	w.tradeAPI = NewHaobtcTrade(name, partner, secret_key)
	return w
}

func (w Haobtc) GetTicker() (ticker common.Ticker, err error) {
	symbol := config.Config["symbol"]

	okTicker, err := w.getTicker(symbol)
	if err != nil {
		return
	}

	ticker.Date = okTicker.Date
	ticker.Ticker.Buy = util.ToFloat(okTicker.Ticker.Buy)
	ticker.Ticker.High = util.ToFloat(okTicker.Ticker.High)
	ticker.Ticker.Last = util.ToFloat(okTicker.Ticker.Last)
	ticker.Ticker.Low = util.ToFloat(okTicker.Ticker.Low)
	ticker.Ticker.Sell = util.ToFloat(okTicker.Ticker.Sell)
	ticker.Ticker.Vol = util.ToFloat(okTicker.Ticker.Vol)
	return
}

//get depth
func (w Haobtc) GetDepth() (orderBook common.OrderBook, err error) {
	symbol := config.Config["symbol"]
	return w.getDepth(symbol)
}

func (w Haobtc) Buy(price, btc string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.BuyBTC(price, btc)
	} else {
		panic(-1)
	}

	return
}

func (w Haobtc) Sell(price, btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.SellBTC(price, btc)
	} else {
		panic(-1)
	}

	return
}

func (w Haobtc) BuyMarket(cny string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.BuyMarketBTC(cny)
	} else {
		panic(-1)
	}

	return
}

func (w Haobtc) SellMarket(btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.SellMarketBTC(btc)
	} else {
		panic(-1)
	}

	return
}

func (w Haobtc) GetOrder(order_id string) (order common.Order, result string, err error) {
	symbol := config.Config["symbol"]
	tradeAPI := w.tradeAPI

	haobtc_order, result, err := tradeAPI.Get_order(symbol, order_id)
	if err != nil {
		return
	}

	order.Id = fmt.Sprintf("%d", haobtc_order.Order_id)
	order.Price = haobtc_order.Avg_price
	order.Amount = haobtc_order.Amount
	order.Deal_amount = haobtc_order.Deal_size

	if haobtc_order.Type == "MARKET" && haobtc_order.Side == "BUY" { //haobtc的市价买单，代表买入金额，比较特殊
		order.Amount = 0
	}

	switch haobtc_order.Status {
	case "PENDING", "SUBMIT", "OPEN":
		order.Status = common.ORDER_STATE_PENDING
	case "CLOSE":
		order.Status = common.ORDER_STATE_SUCCESS
	case "CANCELED":
		order.Status = common.ORDER_STATE_CANCELED // treat canceled status as a error since there is not enough fund
	default:
		order.Status = common.ORDER_STATE_UNKNOWN
	}

	return
}

func (w Haobtc) CancelOrder(order_id string) (err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	return tradeAPI.Cancel_order(symbol, order_id)
}

func (w Haobtc) GetAccount() (account common.Account, err error) {
	tradeAPI := w.tradeAPI

	userInfo, err := tradeAPI.GetAccount()
	if err != nil {
		logger.Debugln("haobtc GetAccount failed")
		return
	}

	account.Available_cny = userInfo.Exchange_cny
	account.Available_btc = userInfo.Exchange_btc

	account.Frozen_cny = userInfo.Exchange_frozen_cny
	account.Frozen_btc = userInfo.Exchange_frozen_btc

	return
}
