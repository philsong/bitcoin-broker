package mocks

import (
	"common"
	"config"
	"db"
	"strategy"
)

func setup() {
	config.Env = config.Test
	config.Root = "/Users/phil/dev/work/trader"
	config.LoadConfig()
	db.Init_sqlstr("root:root@tcp(127.0.0.1:3306)/trader_test")
	//strategy.Task()
}

func setupDepth(chbtc_mockObj, okcoin_mockObj, huobi_mockObj *TradeAPI) {
	// setup expectations
	var chbtc_orderBook common.OrderBook
	chbtc_orderBook.Asks[common.DEPTH-1] = common.MarketOrder{Price: 1402.01, Amount: 0.2}
	chbtc_orderBook.Asks[common.DEPTH-2] = common.MarketOrder{Price: 1400.02, Amount: 0.1}
	chbtc_orderBook.Bids[0] = common.MarketOrder{Price: 1398.02, Amount: 0.1}
	chbtc_orderBook.Bids[1] = common.MarketOrder{Price: 1396.01, Amount: 0.3}
	chbtc_mockObj.On("GetDepth").Return(chbtc_orderBook, nil)

	var okcoin_orderBook common.OrderBook
	okcoin_orderBook.Asks[common.DEPTH-1] = common.MarketOrder{Price: 1401.01, Amount: 0.1}
	okcoin_orderBook.Asks[common.DEPTH-2] = common.MarketOrder{Price: 1400.02, Amount: 0.2}
	okcoin_orderBook.Bids[0] = common.MarketOrder{Price: 1399.02, Amount: 0.1}
	okcoin_orderBook.Bids[1] = common.MarketOrder{Price: 1397.01, Amount: 0.3}
	okcoin_mockObj.On("GetDepth").Return(okcoin_orderBook, nil)

	var huobi_orderBook common.OrderBook
	huobi_orderBook.Asks[common.DEPTH-1] = common.MarketOrder{Price: 1401.01, Amount: 0.1}
	huobi_orderBook.Asks[common.DEPTH-2] = common.MarketOrder{Price: 1400.02, Amount: 0.1}
	huobi_orderBook.Bids[0] = common.MarketOrder{Price: 1398.02, Amount: 0.1}
	huobi_orderBook.Bids[1] = common.MarketOrder{Price: 1397.01, Amount: 0.3}
	huobi_mockObj.On("GetDepth").Return(huobi_orderBook, nil)

	common.SetMockTradeAPI("chbtc", chbtc_mockObj)
	common.SetMockTradeAPI("okcoin", okcoin_mockObj)
	common.SetMockTradeAPI("huobi", huobi_mockObj)

	strategy.QueryDepth()

	return
}

func setupAccount(chbtc_mockObj, okcoin_mockObj, huobi_mockObj *TradeAPI) (chbtc_account, okcoin_account, huobi_account common.Account) {
	// setup expectations

	chbtc_account.Available_cny = 10000
	chbtc_account.Available_btc = 10
	chbtc_mockObj.On("GetAccount").Return(chbtc_account, nil)

	okcoin_account.Available_cny = 20000
	okcoin_account.Available_btc = 20
	okcoin_mockObj.On("GetAccount").Return(okcoin_account, nil)

	huobi_account.Available_cny = 30000
	huobi_account.Available_btc = 30
	huobi_mockObj.On("GetAccount").Return(huobi_account, nil)

	common.SetMockTradeAPI("chbtc", chbtc_mockObj)
	common.SetMockTradeAPI("okcoin", okcoin_mockObj)
	common.SetMockTradeAPI("huobi", huobi_mockObj)

	// call the code we are testing
	strategy.QueryFund()

	return
}
