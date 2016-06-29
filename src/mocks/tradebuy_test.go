package mocks

import (
	"common"
	"logger"
	"strategy"
	"testing"
	"trade_server"
	"trade_service"
)

func Test_Buy(t *testing.T) {
	setup()

	chbtc_mockObj := new(TradeAPI)
	okcoin_mockObj := new(TradeAPI)
	huobi_mockObj := new(TradeAPI)

	setupDepth(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)
	setupAccount(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)

	// create an instance of our test object
	// setup expectations
	{
		buyOrderID := "19"
		chbtc_mockObj.On("Buy", "1402.01", "0.1426").Return(buyOrderID, "", nil)
		okcoin_mockObj.On("BuyMarket", "420.11").Return(buyOrderID, "", nil)
		huobi_mockObj.On("BuyMarket", "280.10").Return(buyOrderID, "", nil)

		var order common.Order
		order.Id = buyOrderID
		order.Price = 1402.01
		order.Amount = 0.1425
		order.Deal_amount = 0.1425
		order.Status = common.ORDER_STATE_SUCCESS
		chbtc_mockObj.On("GetOrder", buyOrderID).Return(order, "", nil)
		order.Status = common.ORDER_STATE_SUCCESS
		okcoin_mockObj.On("GetOrder", buyOrderID).Return(order, "", nil)
		order.Status = common.ORDER_STATE_SUCCESS
		huobi_mockObj.On("GetOrder", buyOrderID).Return(order, "", nil)

		// chbtc_mockObj.On("CancelOrder", buyOrderID).Return(nil)
		// okcoin_mockObj.On("CancelOrder", buyOrderID).Return(nil)
		// huobi_mockObj.On("CancelOrder", buyOrderID).Return(nil)
	}

	common.SetMockTradeAPI("chbtc", chbtc_mockObj)
	common.SetMockTradeAPI("okcoin", okcoin_mockObj)
	common.SetMockTradeAPI("huobi", huobi_mockObj)

	// call the code we are testing
	var buyOrder trade_service.Trade
	buyOrder.UID = "10"
	buyOrder.Amount = 900
	serverhandler := new(trade_server.TradeServiceHandler)
	siteOrder, err := serverhandler.Buy(&buyOrder)
	logger.Infoln(siteOrder, err)

	strategy.ProgressReady()
	strategy.ProgressOrdered()

	// assert that the expectations were met
	chbtc_mockObj.AssertExpectations(t)
	okcoin_mockObj.AssertExpectations(t)
	huobi_mockObj.AssertExpectations(t)
}
