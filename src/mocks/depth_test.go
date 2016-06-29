package mocks

import (
	"github.com/stretchr/testify/assert"
	"strategy"
	"testing"
	// "trade_service"
	"logger"
)

func Test_GetDepth(t *testing.T) {
	setup()

	chbtc_mockObj := new(TradeAPI)
	okcoin_mockObj := new(TradeAPI)
	huobi_mockObj := new(TradeAPI)

	setupDepth(chbtc_mockObj, okcoin_mockObj, huobi_mockObj)

	exchanges := []string{"huobi", "okcoin", "chbtc"}
	asks, bids, markets, err := strategy.GetMergeDepth(exchanges)
	assert.NoError(t, err, "should be no err")
	assert.NotNil(t, asks, "should be not nil")
	assert.NotNil(t, bids, "should be not nil")

	// strategy.PrintDepthList(asks)
	// strategy.PrintDepthList(bids)

	depthCount := 0
	for e := asks.Front(); e != nil; e = e.Next() {
		sumExchangeOrder := e.Value.(*strategy.SumExchangeOrder)
		depthCount++
		switch depthCount {
		case 1:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0, "Amount should be 0")
			assert.EqualValues(t, sumExchangeOrder.Price, 0, "Price should be 0")
		case 2:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.4, "Amount should be 0.4")
			assert.EqualValues(t, sumExchangeOrder.Price, 1400.02, "Price should be 1400.02")
		case 3:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.2, "Amount should be 0.2")
			assert.EqualValues(t, sumExchangeOrder.Price, 1401.01, "Price should be 1401.01")
		case 4:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.2, "Amount should be 0.2")
			assert.EqualValues(t, sumExchangeOrder.Price, 1402.01, "Price should be 1402.01")
		}

		logger.Debugln(depthCount, sumExchangeOrder.Amount, sumExchangeOrder.Price)
		for i := 0; i < len(markets); i++ {
			exchange := markets[i]
			if sumExchangeOrder.ExchangeOrder[exchange] != nil {
				logger.Infoln(depthCount, sumExchangeOrder.ExchangeOrder[exchange])
				switch depthCount {
				case 1:
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0, "Amount should be 0")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 0, "Price should be 0")
				case 2:
					switch exchange {
					case "chbtc":
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.1")
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1400.02, "Price should be 1400.02")
					case "okcoin":
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.2, "Amount should be 0.2")
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1400.02, "Price should be 1400.02")
					case "huobi":
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.1")
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1400.02, "Price should be 1400.02")
					}
				case 3:
					assert.NotEqual(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "chbtc", "Exchange should not be chbtc")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.2")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1401.01, "Price should be 1401.01")
				case 4:
					assert.Equal(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "chbtc", "Exchange should be chbtc")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.2, "Amount should be 0.2")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1402.01, "Price should be 1402.01")
				}
			}
		}
	}

	assert.EqualValues(t, depthCount, 4, "depthCount should be 4")

	depthCount = 0
	for e := bids.Front(); e != nil; e = e.Next() {
		sumExchangeOrder := e.Value.(*strategy.SumExchangeOrder)
		depthCount++

		logger.Debugln(depthCount, sumExchangeOrder.Amount, sumExchangeOrder.Price)

		switch depthCount {
		case 1:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.1, "Amount should be 0.1")
			assert.EqualValues(t, sumExchangeOrder.Price, 1399.02, "Price should be 1399.02")
		case 2:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.2, "Amount should be 0.2")
			assert.EqualValues(t, sumExchangeOrder.Price, 1398.02, "Price should be 1398.02")
		case 3:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.6, "Amount should be 0.6")
			assert.EqualValues(t, sumExchangeOrder.Price, 1397.01, "Price should be 1397.01")
		case 4:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0.3, "Amount should be 0.3")
			assert.EqualValues(t, sumExchangeOrder.Price, 1396.01, "Price should be 1396.01")
		case 5:
			assert.EqualValues(t, sumExchangeOrder.Amount, 0, "Amount should be 0")
			assert.EqualValues(t, sumExchangeOrder.Price, 0, "Price should be 0")
		}

		for i := 0; i < len(markets); i++ {
			exchange := markets[i]
			if sumExchangeOrder.ExchangeOrder[exchange] != nil {
				logger.Infoln(depthCount, sumExchangeOrder.ExchangeOrder[exchange])
				switch depthCount {
				case 1:
					assert.Equal(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "okcoin", "Exchange should be okcoin")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.1")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1399.02, "Price should be 1399.02")
				case 2:
					assert.NotEqual(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "okcoin", "Exchange should not be okcoin")
					switch exchange {
					case "chbtc":
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.1")
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1398.02, "Price should be 1398.02")
					case "huobi":
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.1, "Amount should be 0.1")
						assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1398.02, "Price should be 1398.02")
					}
				case 3:
					assert.NotEqual(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "chbtc", "Exchange should be chbtc")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.3, "Amount should be 0.3")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1397.01, "Price should be 1397.01")
				case 4:
					assert.Equal(t, sumExchangeOrder.ExchangeOrder[exchange].Exchange, "chbtc", "Exchange should be chbtc")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0.3, "Amount should be 0.3")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 1396.01, "Price should be 1396.01")
				case 5:
					assert.Contains(t, markets, exchange, "should contains!")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Amount, 0, "Amount should be 0")
					assert.EqualValues(t, sumExchangeOrder.ExchangeOrder[exchange].Price, 0, "Price should be 0")
				}
			}
		}
	}

	assert.EqualValues(t, depthCount, 5, "depthCount should be 5")

	// assert that the expectations were met
	chbtc_mockObj.AssertExpectations(t)
	okcoin_mockObj.AssertExpectations(t)
	huobi_mockObj.AssertExpectations(t)
}
