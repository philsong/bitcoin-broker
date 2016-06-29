/*
  trader  strategy
*/

package strategy

import (
	"logger"
	"sync"
	"time"
)

var traders map[string]*time.Ticker
var traders_mutex = &sync.Mutex{}

func init() {
	traders = make(map[string]*time.Ticker)
}

func TaskTemplateTrader(seconds time.Duration, f func(exchange string) error, exchange string) {
	ticker := time.NewTicker(seconds * time.Second) // one second
	defer ticker.Stop()

	traders_mutex.Lock()
	if traders[exchange] != nil {
		logger.Errorln("is not nil", exchange, traders[exchange], f)
		traders_mutex.Unlock()
		return
	}

	traders[exchange] = ticker
	traders_mutex.Unlock()

	for t := range ticker.C {
		logger.Debugln(exchange, t)
		f(exchange)
	}
}

func Update_traders() {
	update_exchanges := getConfExchanges()

	for _, update_exchange := range update_exchanges {
		if _, ret := traders[update_exchange]; !ret {
			logger.Infoln("to add trader", update_exchange)
			go TaskTemplateTrader(1, ProgressFSM, update_exchange)
		}
	}
}
