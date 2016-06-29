/*
  trader API Engine
*/

package huobi

import (
	. "common"
	. "config"
	"encoding/json"
	"errors"
	"fmt"
	"logger"
	"util"
)

func (w *Huobi) getTicker(symbol string) (ticker Ticker, err error) {
	ticker_url := fmt.Sprintf(Config["hb_ticker_url"], symbol)
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

func (w *Huobi) getDepth(symbol string) (orderBook OrderBook, err error) {
	depth_url := fmt.Sprintf(Config["hb_depth_url"], symbol, DEPTH)
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
		orderBook.Asks[len(asks)-i-1] = order
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
		orderBook.Bids[i] = order
	}

	return
}
