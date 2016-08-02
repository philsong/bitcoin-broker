/*
  trader  strategy
*/

package strategy

import (
	"config"
	"db"
	"logger"
	"sync"
	"trade_service"
	"util"
)

var buy_queue_mutex = &sync.Mutex{}
var buy_total_amount float64
var sell_queue_mutex = &sync.Mutex{}
var sell_total_amount float64

var buy_limit float64
var sell_limit float64

func initTotalReadyAmount() {
	_buy_queue_amount, prs := config.Config["buy_queue_amount"]
	if !prs {
		_buy_queue_amount = "200"
	}

	_sell_queue_amount, prs := config.Config["sell_queue_amount"]
	if !prs {
		_sell_queue_amount = "200"
	}

	buy_limit = util.ToFloat(_buy_queue_amount)
	sell_limit = util.ToFloat(_sell_queue_amount)

	buy_total, sell_total := db.GetTotalReadyNow()

	logger.Infoln("init limit:", buy_limit, sell_limit, buy_total, sell_total)

	incr_buy(buy_total)
	incr_sell(sell_total)
}

func get_current_buy_total() float64 {
	buy_queue_mutex.Lock()
	defer buy_queue_mutex.Unlock()

	if buy_total_amount < 0 {
		logger.Errorln("buy_total_amount:", buy_total_amount)
		buy_total_amount = 0
	}
	return buy_total_amount
}

func get_current_sell_total() float64 {
	sell_queue_mutex.Lock()
	defer sell_queue_mutex.Unlock()

	if sell_total_amount < 0 {
		logger.Errorln("sell_total_amount:", sell_total_amount)
		sell_total_amount = 0
	}

	return sell_total_amount
}

func get_factor(amount float64, tradeType trade_service.TradeType) float64 {
	if tradeType == trade_service.TradeType_BUY {
		return get_buy_factor(amount)
	} else {
		return get_sell_factor(amount)
	}
}
func get_buy_factor(amount float64) float64 {
	factor := buy_total_amount + amount/buy_limit
	return factor
}

func get_sell_factor(amount float64) float64 {
	factor := sell_total_amount + amount/sell_limit
	return factor
}

func is_limit_buy(amount float64) bool {
	buy_queue_mutex.Lock()
	defer buy_queue_mutex.Unlock()

	if buy_total_amount+amount > buy_limit {
		logger.Infoln("buy_total_amount:", buy_total_amount)
		return true
	}

	buy_total_amount += amount

	logger.Infoln("buy_total_amount:", buy_total_amount)

	return false
}

func is_limit_sell(amount float64) bool {
	sell_queue_mutex.Lock()
	defer sell_queue_mutex.Unlock()

	if sell_total_amount+amount > sell_limit {
		logger.Infoln("sell_total_amount:", sell_total_amount)
		return true
	}

	sell_total_amount += amount

	logger.Infoln("sell_total_amount:", sell_total_amount)
	return false
}

func incr_buy(amount float64) {
	buy_queue_mutex.Lock()
	defer buy_queue_mutex.Unlock()

	buy_total_amount += amount

	logger.Infoln("buy_total_amount:", buy_total_amount)
}

func incr_sell(amount float64) {
	sell_queue_mutex.Lock()
	defer sell_queue_mutex.Unlock()

	sell_total_amount += amount

	logger.Infoln("sell_total_amount:", sell_total_amount)
}

func decr_buy(amount float64) {
	buy_queue_mutex.Lock()
	defer buy_queue_mutex.Unlock()

	buy_total_amount -= amount

	logger.Infoln("buy_total_amount:", buy_total_amount)
}

func decr_sell(amount float64) {
	sell_queue_mutex.Lock()
	defer sell_queue_mutex.Unlock()

	sell_total_amount -= amount

	logger.Infoln("sell_total_amount:", sell_total_amount)
}
