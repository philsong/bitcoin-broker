/*
  trader  strategy
*/

package strategy

import (
	"logger"
	"sync"
	"time"
)

var tickers map[string]*time.Ticker
var tickers_mutex = &sync.Mutex{}

func init() {
	tickers = make(map[string]*time.Ticker)
}

func TaskTemplateTicker(seconds time.Duration, f func(exchange string) error, exchange string) {
	ticker := time.NewTicker(seconds * time.Second) // one second
	defer ticker.Stop()

	tickers_mutex.Lock()
	if tickers[exchange] != nil {
		logger.Errorln("is not nil", exchange, tickers[exchange], f)
		tickers_mutex.Unlock()
		return
	}

	tickers[exchange] = ticker
	tickers_mutex.Unlock()

	for t := range ticker.C {
		logger.Debugln(exchange, t)
		f(exchange)
	}
}

func Update_tickers() {
	update_exchanges := getConfExchanges()

	for current_exchange, _ := range tickers {
		is_remove := true
		for _, update_exchange := range update_exchanges {
			if current_exchange == update_exchange {
				is_remove = false
				break
			}
		}

		if is_remove {
			logger.Infoln("to del ticker", current_exchange)
			tickers_mutex.Lock()
			if _, ret := tickers[current_exchange]; ret {
				tickers[current_exchange].Stop()
				delete(tickers, current_exchange)
			}
			tickers_mutex.Unlock()
		}
	}

	for _, update_exchange := range update_exchanges {
		if _, ret := tickers[update_exchange]; !ret {
			logger.Infoln("to add ticker", update_exchange)
			go TaskTemplateTicker(1, _queryExchangeData, update_exchange)
		}
	}
}
