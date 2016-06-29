/*
  trader API Engine
*/

package okcoin

import (
	. "common"
	. "config"
	"encoding/json"
	"fmt"
	"logger"
	"util"
)

type OKTicker struct {
	Date   string
	Ticker OKTickerPrice
}

type OKTickerPrice struct {
	Buy  string
	High string
	Last string
	Low  string
	Sell string
	Vol  string
}

func (w *Okcoin) getTicker(symbol string) (ticker OKTicker, err error) {
	ticker_url := fmt.Sprintf(Config["ok_ticker_url"], symbol)
	body, err := util.HttpGet(ticker_url)
	if err != nil {
		logger.Errorln(err)
		return
	}

	logger.Infoln(body, err)

	if err = json.Unmarshal([]byte(body), &ticker); err != nil {
		logger.Infoln(err)
		return
	}

	return
}

func (w *Okcoin) getDepth(symbol string) (orderBook OrderBook, err error) {
	depth_url := fmt.Sprintf(Config["ok_depth_url"], symbol, DEPTH)

	logger.Debugln("okcoin", depth_url)
	body, err := util.HttpGet(depth_url)
	if err != nil {
		logger.Errorln(err, depth_url)
		return
	}

	logger.Debugln("okcoin", depth_url, body)
	return w.analyzeOrderBook(body)
}

type OKMarketOrder struct {
	Price  float64 // 价格
	Amount float64 // 委单量
}

type _OKOrderBook struct {
	Asks [DEPTH]interface{}
	Bids [DEPTH]interface{}
}

type OKOrderBook struct {
	Asks [DEPTH]OKMarketOrder
	Bids [DEPTH]OKMarketOrder
}

func convert2struct(_okOrderBook _OKOrderBook) (okOrderBook OKOrderBook) {
	for k, v := range _okOrderBook.Asks {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					okOrderBook.Asks[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					okOrderBook.Asks[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}

	for k, v := range _okOrderBook.Bids {
		switch vt := v.(type) {
		case []interface{}:
			for ik, iv := range vt {
				switch ik {
				case 0:
					okOrderBook.Bids[k].Price = util.InterfaceToFloat64(iv)
				case 1:
					okOrderBook.Bids[k].Amount = util.InterfaceToFloat64(iv)
				}
			}
		}
	}
	return
}

func (w *Okcoin) analyzeOrderBook(content string) (orderBook OrderBook, err error) {
	// init to false
	var _okOrderBook _OKOrderBook
	if err = json.Unmarshal([]byte(content), &_okOrderBook); err != nil {
		logger.Infoln(err)
		return
	}

	okOrderBook := convert2struct(_okOrderBook)

	for i := 0; i < DEPTH; i++ {
		orderBook.Asks[i].Price = okOrderBook.Asks[len(_okOrderBook.Asks)-DEPTH+i].Price
		orderBook.Asks[i].Amount = okOrderBook.Asks[len(_okOrderBook.Asks)-DEPTH+i].Amount
		orderBook.Bids[i].Price = okOrderBook.Bids[i].Price
		orderBook.Bids[i].Amount = okOrderBook.Bids[i].Amount
	}

	return
}
