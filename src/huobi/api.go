/*
  trader API Engine
*/

package huobi

import (
	"common"
	"config"
	"fmt"
	"logger"
	"os"
	"strconv"
	"util"
)

type Huobi struct {
	tradeAPI *HuobiTrade
	name     string
}

func NewExchange(name, access_key, secret_key string) *Huobi {
	w := new(Huobi)
	w.tradeAPI = NewHuobiTrade(access_key, secret_key)
	return w
}

func (w Huobi) GetTicker() (ticker common.Ticker, err error) {
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "btc"
	} else {
		symbol = "ltc"
	}

	return w.getTicker(symbol)
}

func (w Huobi) GetDepth() (orderBook common.OrderBook, err error) {
	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "btc"
	} else {
		symbol = "ltc"
	}

	return w.getDepth(symbol)
}

func (w Huobi) Buy(price, btc string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	return tradeAPI.doTrade("buy", symbol, price, btc)
}

func (w Huobi) Sell(price, btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	return tradeAPI.doTrade("sell", symbol, price, btc)
}

func (w Huobi) BuyMarket(cny string) (buyId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	return tradeAPI.doTrade("buy_market", symbol, "", cny)
}

func (w Huobi) SellMarket(btc string) (sellId string, result string, err error) {
	tradeAPI := w.tradeAPI
	symbol := config.Config["symbol"]
	return tradeAPI.doTrade("sell_market", symbol, "", btc)
}

func (w Huobi) GetOrder(order_id string) (order common.Order, result string, err error) {
	tradeAPI := w.tradeAPI

	symbol := config.Config["symbol"]
	if symbol == "btc_cny" {
		symbol = "1"
	} else if symbol == "ltc_cny" {
		symbol = "0"
	} else {
		panic(-1)
	}
	hbOrder, result, err := tradeAPI.Get_order(symbol, order_id)
	if err != nil {
		return
	}

	order.Id = fmt.Sprintf("%d", hbOrder.Id)

	price, err := strconv.ParseFloat(hbOrder.Processed_price, 64)
	if err != nil {
		logger.Errorln("config item order_price is not float")
		return
	}

	amount, err := strconv.ParseFloat(hbOrder.Order_amount, 64)
	if err != nil {
		logger.Errorln("config item order_amount is not float")
		return
	}

	deal_amount, err := strconv.ParseFloat(hbOrder.Processed_amount, 64)
	if err != nil {
		logger.Errorln("config item processed_amount is not float")
		return
	}

	order.Price = price

	//1限价买　2限价卖　3市价买　4市价卖
	if hbOrder.Type == 3 { //火币的市价买单，代表买入金额和成交金额，比较特殊
		if price > 1 {
			order.Amount = amount / price
			order.Deal_amount = deal_amount / price
		}
	} else if hbOrder.Type == 4 || hbOrder.Type == 2 || hbOrder.Type == 1 {
		order.Amount = amount
		order.Deal_amount = deal_amount
	} else {
		logger.Panicln("unsupport trade type!")
		os.Exit(-1)
	}

	switch hbOrder.Status {
	case 0, 1:
		order.Status = common.ORDER_STATE_PENDING
	case 2:
		order.Status = common.ORDER_STATE_SUCCESS
	case 3, 6:
		order.Status = common.ORDER_STATE_CANCELED
	default:
		order.Status = common.ORDER_STATE_UNKNOWN
	}

	return
}

func (w Huobi) CancelOrder(order_id string) (err error) {
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

func (w Huobi) GetAccount() (account common.Account, err error) {
	tradeAPI := w.tradeAPI

	userInfo, err := tradeAPI.GetAccount()

	if err != nil {
		logger.Debugln("Huobi GetAccount failed", err)
		return
	}
	account.Available_cny = util.ToFloat(userInfo.Available_cny_display)
	account.Available_btc = util.ToFloat(userInfo.Available_btc_display)
	account.Available_ltc = util.ToFloat(userInfo.Available_ltc_display)

	account.Frozen_cny = util.ToFloat(userInfo.Frozen_cny_display)
	account.Frozen_btc = util.ToFloat(userInfo.Frozen_btc_display)
	account.Frozen_ltc = util.ToFloat(userInfo.Frozen_ltc_display)

	return
}
