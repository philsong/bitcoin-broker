/*
  trader API Engine
*/

package haobtc

import (
	. "common"
	. "config"
	"encoding/json"
	"fmt"
	"logger"
	"util"
)

func (w *Haobtc) getTicker(symbol string) (ticker Ticker, err error) {
	ticker_url := fmt.Sprintf(Config[w.name+"_ticker_url"], symbol)
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

func (w *Haobtc) getDepth(symbol string) (orderBook OrderBook, err error) {
	depth_url := fmt.Sprintf(Config[w.name+"_depth_url"], DEPTH)

	logger.Debugln("Haobtc", depth_url)
	body, err := util.HttpGet(depth_url)
	if err != nil {
		logger.Errorln(err, depth_url)
		return
	}

	logger.Debugln("Haobtc", depth_url, body)
	return w.analyzeOrderBook(body)
}

type _OKOrderBook struct {
	Asks [DEPTH]interface{}
	Bids [DEPTH]interface{}
}

func convert2struct(_okOrderBook _OKOrderBook) (okOrderBook OrderBook) {
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

func (w *Haobtc) analyzeOrderBook(content string) (orderBook OrderBook, err error) {
	// init to false
	var _okOrderBook _OKOrderBook
	if err = json.Unmarshal([]byte(content), &_okOrderBook); err != nil {
		logger.Errorln(err)
		return
	}

	okOrderBook := convert2struct(_okOrderBook)

	for i := 0; i < DEPTH; i++ {
		orderBook.Asks[i].Price = okOrderBook.Asks[len(_okOrderBook.Asks)-i-1].Price
		orderBook.Asks[i].Amount = okOrderBook.Asks[len(_okOrderBook.Asks)-i-1].Amount
		orderBook.Bids[i].Price = okOrderBook.Bids[i].Price
		orderBook.Bids[i].Amount = okOrderBook.Bids[i].Amount
	}

	return
}
