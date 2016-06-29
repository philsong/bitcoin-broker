/*
  trader API Engine
*/

package okcoin

import (
	"common"
	"config"
	"fmt"
	"logger"
	"util"
)

import s "strings"

type Okcoin struct {
	tradeAPI *OkcoinTrade
	name     string
}

func NewExchange(name, partner, secret_key string) *Okcoin {
	w := new(Okcoin)
	w.tradeAPI = NewOkcoinTrade(partner, secret_key)
	return w
}

func (w Okcoin) GetTicker() (ticker common.Ticker, err error) {
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
func (w Okcoin) GetDepth() (orderBook common.OrderBook, err error) {
	symbol := config.Config["symbol"]
	return w.getDepth(symbol)
}

func (w Okcoin) Buy(price, btc string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.BuyBTC(price, btc)
	} else {
		panic(-1)
	}

	return
}

func (w Okcoin) Sell(price, btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.SellBTC(price, btc)
	} else {
		panic(-1)
	}

	return
}

func (w Okcoin) BuyMarket(cny string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.BuyMarketBTC(cny)
	} else {
		panic(-1)
	}

	return
}

func (w Okcoin) SellMarket(btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		return tradeAPI.SellMarketBTC(btc)
	} else {
		panic(-1)
	}

	return
}

func (w Okcoin) GetOrder(order_id string) (order common.Order, result string, err error) {
	symbol := config.Config["symbol"]
	tradeAPI := w.tradeAPI

	ok_orderTable, result, err := tradeAPI.Get_order(symbol, order_id)
	if err != nil {
		return
	}

	order.Id = fmt.Sprintf("%d", ok_orderTable.Orders[0].Order_id)
	order.Price = ok_orderTable.Orders[0].Avg_price
	order.Amount = ok_orderTable.Orders[0].Amount //注意：坑！ok的买单，amount返回为空
	order.Deal_amount = ok_orderTable.Orders[0].Deal_amount

	switch ok_orderTable.Orders[0].Status {
	case 0, 1:
		order.Status = common.ORDER_STATE_PENDING
	case 2:
		order.Status = common.ORDER_STATE_SUCCESS
	case -1:
		order.Status = common.ORDER_STATE_CANCELED // treat canceled status as a error since there is not enough fund
	default:
		order.Status = common.ORDER_STATE_UNKNOWN
	}

	return
}

func (w Okcoin) CancelOrder(order_id string) (err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	return tradeAPI.Cancel_order(symbol, order_id)
}

func (w Okcoin) GetAccount() (account common.Account, err error) {
	tradeAPI := w.tradeAPI

	userInfo, err := tradeAPI.GetAccount()
	if err != nil {
		logger.Debugln("okcoin GetAccount failed")
		return
	}

	cnystr := userInfo.Info.Funds.Free.CNY
	cnystr = s.Replace(cnystr, ",", "", -1)

	account.Available_cny = util.ToFloat(cnystr)

	account.Available_btc = util.ToFloat(userInfo.Info.Funds.Free.BTC)
	account.Available_ltc = util.ToFloat(userInfo.Info.Funds.Free.LTC)

	account.Frozen_cny = util.ToFloat(userInfo.Info.Funds.Freezed.CNY)
	account.Frozen_btc = util.ToFloat(userInfo.Info.Funds.Freezed.BTC)
	account.Frozen_ltc = util.ToFloat(userInfo.Info.Funds.Freezed.LTC)

	return
}
