/*
  trader  strategy
*/

package strategy

import (
	"common"
	"config"
	"db"
	"errors"
	"fmt"
	"haobtc"
	"huobi"
	"logger"
	"okcoin"
	"sync"
	"trade_service"
)

var fundExchages map[string]trade_service.Account
var orderExchages map[string][]string
var mutex = &sync.RWMutex{}

func init() {
	orderExchages = make(map[string][]string)
	fundExchages = make(map[string]trade_service.Account)
}

func getConfExchanges() (markets []string) {
	if config.Env == config.Test {
		markets = []string{"huobi", "okcoin", "haobtc"}
		return
	}

	exchange_configs, err := db.GetExchangeConfigs()
	if err != nil {
		return
	}

	for _, value := range exchange_configs {
		markets = append(markets, value.Exchange)
	}

	return
}

func GetUsableExchange(method string, check_fund bool) (markets []string) {
	all_accounts, err := db.GetAccount()
	if err != nil {
		logger.Errorln(err)
		return
	}

	var accounts []*trade_service.Account
	for i := 0; i < len(all_accounts); i++ {
		if !(all_accounts[i].GetPauseTrade()) && (all_accounts[i].GetExchange()) != "haobtc" {
			accounts = append(accounts, all_accounts[i])
		}
	}

	if len(accounts) == 0 {
		logger.Errorln("No Usable Exchange in DB:", accounts)
		return
	}

	if !check_fund {
		for i := 0; i < len(accounts); i++ {
			markets = append(markets, accounts[i].Exchange)
		}
		return
	}

	cny_threshold := 200000.0
	btc_threshold := 30.0
	amount_config, err := db.GetAmountConfig()
	if err != nil {
		logger.Errorln(err)
		// return
	} else {
		cny_threshold = amount_config.MaxCny
		btc_threshold = amount_config.MaxBtc
	}

	// logger.Errorln(cny_threshold, btc_threshold)

	for i := 0; i < len(accounts); i++ {
		if (method == "BUY" && accounts[i].GetAvailableCny() > cny_threshold) ||
			(method == "SELL" && accounts[i].GetAvailableBtc() > btc_threshold) {
			markets = append(markets, accounts[i].Exchange)
		}
	}

	if len(markets) == 0 {
		logger.Errorln("GetUsableExchange: No UsableExchange:", method, cny_threshold, btc_threshold, accounts)
	}

	return
}

func _performDel(markets []string, exchange string) []string {
	for i := 0; i < len(markets); i++ {
		if markets[i] == exchange {
			l := len(markets) - 1
			markets[i] = markets[l]
			markets = markets[:l]
		}
	}

	return markets
}

func nextExchange(id int64, method, exchange string) []string {
	key := fmt.Sprintf("%d%s", id, method)
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := orderExchages[key]
	// logger.Infoln(exists, exchange, key, orderExchages[key])
	if !exists {
		exchanges := GetUsableExchange(method, false)
		orderExchages[key] = make([]string, len(exchanges))
		copy(orderExchages[key], exchanges)
	}

	newExchanges := _performDel(orderExchages[key], exchange)

	orderExchages[key] = make([]string, len(newExchanges))
	copy(orderExchages[key], newExchanges)

	// logger.Infoln(exchange, key, orderExchages[key])
	return orderExchages[key]
}

func GetExchange(exchange string) (tradeAPI common.TradeAPI, err error) {
	if config.Env == config.Test {
		tradeAPI = common.GetMockTradeAPI(exchange)
		return
	}

	exchange_config, err := db.GetExchangeConfig(exchange)
	if err != nil {
		logger.Errorln(err)
		return
	}

	switch exchange {
	case "huobi":
		return huobi.NewExchange(exchange, exchange_config.AccessKey, exchange_config.SecretKey), nil
	case "okcoin":
		return okcoin.NewExchange(exchange, exchange_config.AccessKey, exchange_config.SecretKey), nil
	case "haobtc":
		return haobtc.NewExchange(exchange, exchange_config.AccessKey, exchange_config.SecretKey), nil
	default:
		err = errors.New("unknow exchange name")
		logger.Errorln("unknow exchange name?")
		return
	}

	return
}
