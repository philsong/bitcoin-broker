/*
  trader API Engine
*/

package chbtc

import (
	. "common"
	. "config"
	"encoding/json"
	"errors"
	"logger"
	"util"
)

type ChTicker struct {
	Ticker ChTickerPrice
}

type ChTickerPrice struct {
	Buy  string
	High string
	Last string
	Low  string
	Sell string
	Vol  string
}

func (w *Chbtc) getTicker(symbol string) (ticker ChTicker, err error) {
	ticker_url := Config["ch_ticker_url"]
	body, err := util.HttpGet(ticker_url)
	if err != nil {
		logger.Errorln(err)
		return
	}

	if err = json.Unmarshal([]byte(body), &ticker); err != nil {
		logger.Infoln(err)
		return
	}

	return
}

func (w *Chbtc) getDepth(symbol string) (orderBook OrderBook, err error) {
	depth_url := Config["ch_depth_url"]
	body, err := util.HttpGet(depth_url)
	if err != nil {
		logger.Errorln(err)
		return
	}

	defaultstruct := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &defaultstruct)
	if err != nil {
		logger.Infoln("defaultstruct", defaultstruct)
		return
	}

	asks := defaultstruct["asks"].([]interface{})
	bids := defaultstruct["bids"].([]interface{})

	for i, ask := range asks {
		_ask := ask.([]interface{})
		price, ret := _ask[0].(float64)
		if !ret {
			err = errors.New("data wrong")
			return orderBook, err
		}
		amount, ret := _ask[1].(float64)
		if !ret {
			err = errors.New("data wrong")
			return orderBook, err
		}
		order := MarketOrder{
			Price:  price,
			Amount: amount,
		}
		if i < DEPTH {
			orderBook.Asks[i] = order
		} else {
			break
		}
	}

	for i, bid := range bids {
		_bid := bid.([]interface{})
		price, ret := _bid[0].(float64)
		if !ret {
			err = errors.New("data wrong")
			return orderBook, err
		}
		amount, ret := _bid[1].(float64)
		if !ret {
			err = errors.New("data wrong")
			return orderBook, err
		}
		order := MarketOrder{
			Price:  price,
			Amount: amount,
		}

		if i < DEPTH {
			orderBook.Bids[i] = order
		} else {
			break
		}
	}

	return
}
