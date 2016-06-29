/*
  trader  strategy
*/

package strategy

import (
	"common"
	"container/list"
	"db"
	"fmt"
	"logger"
	"trade_service"
)

type ExchangeOrder struct {
	Exchange string
	Amount   float64
	Price    float64
}

type SumExchangeOrder struct {
	Amount        float64
	Price         float64
	ExchangeOrder map[string]*ExchangeOrder //key is exchange
}

func mergeDepth(exchange, depth_type string, marketOrders []common.MarketOrder, depthList *list.List) (err error) {
	size := len(marketOrders)

	for i := 0; i < size; i++ {
		marketOrder := marketOrders[i]

		e := depthList.Front()
		for nil != e {
			if e.Value.(*SumExchangeOrder).Price == marketOrder.Price {
				//combine to the depth list item, just needd update the list item, so just break in the end.
				sumExchangeOrder := e.Value.(*SumExchangeOrder)

				if sumExchangeOrder.ExchangeOrder == nil {
					//cannot arrive here, just in case
					logger.Errorln("mergeDepth exception.")
					sumExchangeOrder.ExchangeOrder = make(map[string]*ExchangeOrder)
				}

				if sumExchangeOrder.ExchangeOrder[exchange] == nil {
					sumExchangeOrder.ExchangeOrder[exchange] = new(ExchangeOrder)
				}

				sumExchangeOrder.ExchangeOrder[exchange].Exchange = exchange
				sumExchangeOrder.ExchangeOrder[exchange].Price = marketOrder.Price
				sumExchangeOrder.ExchangeOrder[exchange].Amount = marketOrder.Amount
				sumExchangeOrder.Amount += marketOrder.Amount
				sumExchangeOrder.Price = marketOrder.Price

				break
			} else if (depth_type == "asks" && marketOrder.Price < e.Value.(*SumExchangeOrder).Price) ||
				(depth_type == "bids" && marketOrder.Price > e.Value.(*SumExchangeOrder).Price) {

				sumExchangeOrder := new(SumExchangeOrder)
				sumExchangeOrder.ExchangeOrder = make(map[string]*ExchangeOrder)
				sumExchangeOrder.ExchangeOrder[exchange] = new(ExchangeOrder)
				sumExchangeOrder.ExchangeOrder[exchange].Exchange = exchange
				sumExchangeOrder.ExchangeOrder[exchange].Price = marketOrder.Price
				sumExchangeOrder.ExchangeOrder[exchange].Amount = marketOrder.Amount
				sumExchangeOrder.Amount += marketOrder.Amount
				sumExchangeOrder.Price = marketOrder.Price

				depthList.InsertBefore(sumExchangeOrder, e)

				break
			}
			e = e.Next()
		}

		//the biggest,put @v on the back of the list
		if nil == e {
			sumExchangeOrder := new(SumExchangeOrder)
			sumExchangeOrder.ExchangeOrder = make(map[string]*ExchangeOrder)
			sumExchangeOrder.ExchangeOrder[exchange] = new(ExchangeOrder)
			sumExchangeOrder.ExchangeOrder[exchange].Exchange = exchange
			sumExchangeOrder.ExchangeOrder[exchange].Price = marketOrder.Price
			sumExchangeOrder.ExchangeOrder[exchange].Amount = marketOrder.Amount
			sumExchangeOrder.Amount += marketOrder.Amount
			sumExchangeOrder.Price = marketOrder.Price

			depthList.PushBack(sumExchangeOrder)
		}
	}

	return
}

func GetMergeDepth(markets []string) (asks, bids *list.List, newMarkets []string, err error) {
	asks = list.New()
	bids = list.New()

	errcount := 0
	for i := 0; i < len(markets); i++ {
		exchange := markets[i]
		orderbook, _err := db.GetDepth(exchange)
		if _err != nil {
			errcount++
			err = _err
			logger.Errorln("GetMergeDepth err:", exchange, err, len(markets), errcount)
			continue
		}

		// okcoin except special progress
		size := len(orderbook.Asks)
		if size > 0 && orderbook.Asks[common.DEPTH-1].Price < 0.000001 {
			logger.Errorln("GetMergeDepth exception orderbook:", exchange, orderbook)
			continue
		}

		newMarkets = append(newMarkets, exchange)
		mergeDepth(exchange, "asks", orderbook.Asks[:], asks)
		mergeDepth(exchange, "bids", orderbook.Bids[:], bids)
	}

	if errcount < len(markets) {
		err = nil
	} else {
		logger.Errorln("GetMergeDepth failed, restart", len(markets), errcount)
		return
		// os.Exit(-1) //triger supervisor restart to fix the broken net connection.
	}

	return
}

func PrintDepthList(depthList *list.List, markets []string) {
	return

	logger.Infoln("analyzeDepth depthList", depthList)
	depthCount := 0
	for e := depthList.Front(); e != nil; e = e.Next() {
		sumExchangeOrder := e.Value.(*SumExchangeOrder)
		depthCount++
		logger.Infoln(depthCount, sumExchangeOrder.Amount, sumExchangeOrder.Price)
		for i := 0; i < len(markets); i++ {
			exchange := markets[i]
			if sumExchangeOrder.ExchangeOrder[exchange] != nil {
				logger.Infoln(depthCount, sumExchangeOrder.ExchangeOrder[exchange])
			}
		}
	}

	logger.Infoln("depthCount:", depthCount)
}

func analyzeAskDepth(ncny float64, markets []string) (tradeOrders map[string]*trade_service.TradeOrder, err error) {
	depthList, _, newMarkets, err := GetMergeDepth(markets)
	if err != nil {
		logger.Errorln(err)
		return
	}

	PrintDepthList(depthList, newMarkets)

	sum_cny := 0.0

	tradeOrders = make(map[string]*trade_service.TradeOrder)
	for e := depthList.Front(); e != nil; e = e.Next() {
		sumExchangeOrder := e.Value.(*SumExchangeOrder)
		ask_price := sumExchangeOrder.Price
		ask_vol := sumExchangeOrder.Amount

		if sum_cny+ask_price*ask_vol > ncny {
			for i := 0; i < len(newMarkets); i++ {
				exchange := newMarkets[i]
				logger.Debugln(i, exchange)

				if sumExchangeOrder.ExchangeOrder[exchange] == nil {
					continue
				}

				if tradeOrders[exchange] == nil {
					tradeOrders[exchange] = new(trade_service.TradeOrder)
				}

				sub_vol := sumExchangeOrder.ExchangeOrder[exchange].Amount
				if sum_cny+ask_price*sub_vol > ncny {
					left_vol := (ncny - sum_cny) / ask_price

					tradeOrders[exchange].EstimateBtc += left_vol
					tradeOrders[exchange].EstimatePrice = ask_price
					tradeOrders[exchange].EstimateCny += left_vol * ask_price
					sum_cny += left_vol * ask_price
					break
				} else {
					tradeOrders[exchange].EstimateBtc += sub_vol
					tradeOrders[exchange].EstimatePrice = ask_price
					tradeOrders[exchange].EstimateCny += sub_vol * ask_price
					sum_cny += sub_vol * ask_price
				}
			}

			break
		} else { //<=
			for i := 0; i < len(newMarkets); i++ {
				exchange := newMarkets[i]

				if sumExchangeOrder.ExchangeOrder[exchange] == nil {
					continue
				}

				if tradeOrders[exchange] == nil {
					tradeOrders[exchange] = new(trade_service.TradeOrder)
				}

				tradeOrders[exchange].EstimateBtc += sumExchangeOrder.ExchangeOrder[exchange].Amount
				tradeOrders[exchange].EstimatePrice = ask_price
				tradeOrders[exchange].EstimateCny += sumExchangeOrder.ExchangeOrder[exchange].Amount * ask_price

				sum_cny += sumExchangeOrder.ExchangeOrder[exchange].Amount * ask_price
			}
		}
	}

	return
}

func analyzeBidDepth(nbtc float64, markets []string) (tradeOrders map[string]*trade_service.TradeOrder, err error) {
	_, depthList, newMarkets, err := GetMergeDepth(markets)
	if err != nil {
		logger.Errorln(err)
		return
	}

	PrintDepthList(depthList, newMarkets)

	sum_btc := 0.0

	tradeOrders = make(map[string]*trade_service.TradeOrder)
	for e := depthList.Front(); e != nil; e = e.Next() {
		sumExchangeOrder := e.Value.(*SumExchangeOrder)
		bid_price := sumExchangeOrder.Price
		bid_vol := sumExchangeOrder.Amount

		if sum_btc+bid_vol > nbtc {
			for i := 0; i < len(newMarkets); i++ {
				exchange := newMarkets[i]
				//logger.Infoln(i, exchange)

				if sumExchangeOrder.ExchangeOrder[exchange] == nil {
					continue
				}

				if tradeOrders[exchange] == nil {
					tradeOrders[exchange] = new(trade_service.TradeOrder)
				}

				sub_vol := sumExchangeOrder.ExchangeOrder[exchange].Amount
				if sum_btc+sub_vol > nbtc {
					left_vol := (nbtc - sum_btc)

					tradeOrders[exchange].EstimateBtc += left_vol
					tradeOrders[exchange].EstimatePrice = bid_price
					tradeOrders[exchange].EstimateCny += left_vol * bid_price
					sum_btc += left_vol
					break
				} else {
					tradeOrders[exchange].EstimateBtc += sub_vol
					tradeOrders[exchange].EstimatePrice = bid_price
					tradeOrders[exchange].EstimateCny += sub_vol * bid_price
					sum_btc += sub_vol
				}
			}

			break
		} else { //<=
			for i := 0; i < len(newMarkets); i++ {
				exchange := newMarkets[i]

				if sumExchangeOrder.ExchangeOrder[exchange] == nil {
					continue
				}

				if tradeOrders[exchange] == nil {
					tradeOrders[exchange] = new(trade_service.TradeOrder)
				}

				tradeOrders[exchange].EstimateBtc += sumExchangeOrder.ExchangeOrder[exchange].Amount
				tradeOrders[exchange].EstimatePrice = bid_price
				tradeOrders[exchange].EstimateCny += sumExchangeOrder.ExchangeOrder[exchange].Amount * bid_price

				sum_btc += sumExchangeOrder.ExchangeOrder[exchange].Amount
			}
		}
	}

	return
}

func estimateOrder(siteOrder *db.SiteOrder, markets []string) (tradeOrders map[string]*trade_service.TradeOrder, tradeResult_ *trade_service.TradeException) {
	amount := siteOrder.Amount
	tradeType := siteOrder.TradeType

	logger.Infoln("estimateOrder markets:", tradeType, amount, markets)
	if len(markets) == 0 {
		tradeResult_ = trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_NO_USABLE_FUND

		logger.Errorln(tradeType.String(), tradeResult_)
		return
	}

	var err error
	if tradeType == trade_service.TradeType_BUY {
		tradeOrders, err = analyzeAskDepth(amount, markets)
	} else {
		tradeOrders, err = analyzeBidDepth(amount, markets)
	}

	if err != nil {
		logger.Errorln(err)
		tradeResult_ = trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_NO_USABLE_DEPTH
		return
	}

	estimate_btc := 0.0
	estimate_cny := 0.0
	estimate_price := 0.0
	for exchange, _ := range tradeOrders {
		tradeOrder := tradeOrders[exchange]
		tradeOrder.SiteOrderID = siteOrder.ID
		tradeOrder.Exchange = exchange
		tradeOrder.TradeType = tradeType

		if use_time_weighted_algorithm {
			tradeOrder.OrderStatus = trade_service.OrderStatus_TIME_WEIGHTED
		} else {
			tradeOrder.OrderStatus = trade_service.OrderStatus_READY
		}
		estimate_cny += tradeOrder.EstimateCny
		estimate_btc += tradeOrder.EstimateBtc

		tradePrice := fmt.Sprintf("%0.2f", tradeOrder.EstimatePrice)
		tradeBTC := fmt.Sprintf("%0.4f", tradeOrder.EstimateBtc)
		tradeCNY := fmt.Sprintf("%0.2f", tradeOrder.EstimateCny)

		logger.Infoln("estimateOrder trade:", exchange, tradePrice, tradeBTC, tradeCNY)
	}

	if len(tradeOrders) > 1 {
		toExchange := ""
		for exchange, _ := range tradeOrders {
			tradeOrder := tradeOrders[exchange]
			if tradeOrder.EstimateBtc > unit_min_amount {
				toExchange = exchange
				break
			}
		}

		if toExchange != "" {
			for exchange, _ := range tradeOrders {
				tradeOrder := tradeOrders[exchange]
				if tradeOrder.EstimateBtc <= unit_min_amount {
					if toExchange != exchange {
						logger.Debugln("combine trade begin", tradeType.String(), tradeOrders[toExchange].EstimatePrice, tradeOrders)
						logger.Debugln("to add trade", toExchange, tradeOrders[toExchange])
						logger.Debugln("to del trade", exchange, tradeOrder)
						tradeOrder.Exchange = toExchange
						tradeOrders[toExchange].EstimateBtc += tradeOrders[exchange].EstimateBtc
						if tradeType == trade_service.TradeType_BUY {
							if tradeOrders[toExchange].EstimatePrice < tradeOrders[exchange].EstimatePrice {
								tradeOrders[toExchange].EstimatePrice = tradeOrders[exchange].EstimatePrice
							}
						} else {
							if tradeOrders[toExchange].EstimatePrice > tradeOrders[exchange].EstimatePrice {
								tradeOrders[toExchange].EstimatePrice = tradeOrders[exchange].EstimatePrice
							}
						}

						tradeOrders[toExchange].EstimateCny += tradeOrders[exchange].EstimateCny

						delete(tradeOrders, exchange)

						logger.Debugln("combine trade end", tradeOrders[toExchange].EstimatePrice, tradeOrders)
					}
				}
			}
		}
	}

	if tradeType == trade_service.TradeType_BUY {
		if estimate_cny+0.01 < amount {
			tradeResult_ = trade_service.NewTradeException()
			tradeResult_.Reason = trade_service.EX_DEPTH_INSUFFICIENT

			logger.Errorln(tradeResult_, estimate_cny, amount)
			return
		}
	} else {
		if estimate_btc+0.01 < amount {
			tradeResult_ = trade_service.NewTradeException()
			tradeResult_.Reason = trade_service.EX_DEPTH_INSUFFICIENT

			logger.Errorln(tradeResult_, estimate_btc, amount)
			return
		}
	}

	if estimate_btc > 0.01 {
		estimate_price = estimate_cny / estimate_btc
	} else {
		tradeResult_ = trade_service.NewTradeException()
		tradeResult_.Reason = trade_service.EX_INTERNAL_ERROR

		logger.Errorln(tradeResult_, estimate_btc, siteOrder)
		return
	}

	logger.Infoln("estimateOrder result:", amount, siteOrder.TradeType, estimate_cny, estimate_btc, estimate_price)

	siteOrder.EstimatePrice = estimate_price
	siteOrder.EstimateCny = estimate_cny
	siteOrder.EstimateBtc = estimate_btc

	if use_time_weighted_algorithm {
		siteOrder.OrderStatus = trade_service.OrderStatus_TIME_WEIGHTED
	} else {
		siteOrder.OrderStatus = trade_service.OrderStatus_READY
	}

	tradeResult_ = nil

	return
}
