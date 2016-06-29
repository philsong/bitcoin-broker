package mocks

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"trade_server"
)

func Test_GetAccount(t *testing.T) {
	setup()

	// create an instance of our test object
	chbtc_mockObj := new(TradeAPI)
	okcoin_mockObj := new(TradeAPI)
	huobi_mockObj := new(TradeAPI)

	chbtc_account, okcoin_account, huobi_account := setupAccount(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)

	serverhandler := new(trade_server.TradeServiceHandler)
	exchange_accounts, err := serverhandler.GetAccount()
	assert.NoError(t, err, "should be no err")

	for i := 0; i < len(exchange_accounts); i++ {
		exchange_account := exchange_accounts[i]
		switch exchange_account.Exchange {
		case "chbtc":
			assert.Equal(t, exchange_account.AvailableCny, chbtc_account.Available_cny, "should same cny amount")
			assert.Equal(t, exchange_account.AvailableBtc, chbtc_account.Available_btc, "should same btc amount")
		case "okcoin":
			assert.Equal(t, exchange_account.AvailableCny, okcoin_account.Available_cny, "should same cny amount")
			assert.Equal(t, exchange_account.AvailableBtc, okcoin_account.Available_btc, "should same btc amount")
		case "huobi":
			assert.Equal(t, exchange_account.AvailableCny, huobi_account.Available_cny, "should same cny amount")
			assert.Equal(t, exchange_account.AvailableBtc, huobi_account.Available_btc, "should same btc amount")
		}
	}

	// assert that the expectations were met
	chbtc_mockObj.AssertExpectations(t)
	okcoin_mockObj.AssertExpectations(t)
	huobi_mockObj.AssertExpectations(t)
}
