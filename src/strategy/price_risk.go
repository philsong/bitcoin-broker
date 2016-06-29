/*
  trader  risk
*/

package strategy

import (
	"config"
	"db"
	"logger"
	"math"
	"strconv"
	"trade_service"
	"util"
)

const MeasurementError = 0.000001

func check_price_limit(siteOrder *db.SiteOrder) bool {
	_price_threshold, prs := config.Config["price_threshold"]
	if !prs {
		_price_threshold = "5"
	}

	price_threshold := util.ToFloat(_price_threshold)

	if siteOrder.TradeType == trade_service.TradeType_BUY {
		if siteOrder.EstimatePrice-siteOrder.Price > price_threshold {
			return false
		}
	} else {
		if siteOrder.Price-siteOrder.EstimatePrice > price_threshold {
			return false
		}
	}

	return true
}

func Check_ticker_limit(price float64, tradeType trade_service.TradeType) bool {
	if price < MeasurementError {
		return false
	}

	_ticker_compare_count, prs := config.Config["ticker_compare_count"]
	if !prs {
		_ticker_compare_count = "3"
	}

	ticker_compare_count, err := strconv.ParseInt(_ticker_compare_count, 10, 64)
	if err != nil {
		ticker_compare_count = 3
	}

	tickers, err := db.GetNTickers(int(ticker_compare_count))
	if err != nil {
		logger.Infoln("GetNTickers err:", err)
		return false
	}

	if len(tickers) < 1 {
		logger.Infoln(tickers, "tickers is empty")
		return false
	}

	// logger.Infoln(len(tickers), tickers)

	isOurTicker := false
	for i := 0; i < len(tickers); i++ {
		ticker_price := 0.0
		if tradeType == trade_service.TradeType_BUY {
			ticker_price = tickers[i].Ask
		} else if tradeType == trade_service.TradeType_SELL {
			ticker_price = tickers[i].Bid
		} else {
			logger.Infoln("invalid trade type", tradeType)
			return false
		}

		logger.Infoln("check ticker limit:", ticker_price, price, math.Abs(ticker_price-price), MeasurementError)
		if math.Abs(ticker_price-price) < MeasurementError {
			isOurTicker = true
			break
		}
	}

	return isOurTicker
}
