package mocks

import (
	"common"
	"logger"
	"strategy"
	"testing"
	"trade_server"
	"trade_service"
)

func Test_Sell(t *testing.T) {
	setup()

	chbtc_mockObj := new(TradeAPI)
	okcoin_mockObj := new(TradeAPI)
	huobi_mockObj := new(TradeAPI)

	setupDepth(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)
	setupAccount(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)

	// create an instance of our test object
	// setup expectations
	{
		orderID := "19"
		chbtc_mockObj.On("Sell", "1396.01", "0.2000").Return(orderID, "", nil)
		okcoin_mockObj.On("SellMarket", "0.4000").Return(orderID, "", nil)
		huobi_mockObj.On("SellMarket", "0.4000").Return(orderID, "", nil)

		var order common.Order
		order.Id = orderID
		order.Price = 1396.01
		order.Amount = 0.2000
		order.Deal_amount = 0.2000
		order.Status = common.ORDER_STATE_SUCCESS
		chbtc_mockObj.On("GetOrder", orderID).Return(order, "", nil)
		order.Status = common.ORDER_STATE_SUCCESS
		okcoin_mockObj.On("GetOrder", orderID).Return(order, "", nil)
		order.Status = common.ORDER_STATE_SUCCESS
		huobi_mockObj.On("GetOrder", orderID).Return(order, "", nil)

		// chbtc_mockObj.On("CancelOrder", orderID).Return(nil)
		// okcoin_mockObj.On("CancelOrder", buyOrderID).Return(nil)
		// huobi_mockObj.On("CancelOrder", buyOrderID).Return(nil)
	}

	common.SetMockTradeAPI("chbtc", chbtc_mockObj)
	common.SetMockTradeAPI("okcoin", okcoin_mockObj)
	common.SetMockTradeAPI("huobi", huobi_mockObj)

	// call the code we are testing
	var sellOrder trade_service.Trade
	sellOrder.UID = "10"
	sellOrder.Amount = 1
	serverhandler := new(trade_server.TradeServiceHandler)
	siteOrder, err := serverhandler.Sell(&sellOrder)
	logger.Infoln(siteOrder, err)

	strategy.ProgressReady()
	strategy.ProgressOrdered()

	// assert that the expectations were met
	chbtc_mockObj.AssertExpectations(t)
	okcoin_mockObj.AssertExpectations(t)
	huobi_mockObj.AssertExpectations(t)
}
