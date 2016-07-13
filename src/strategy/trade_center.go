/*
  trader  strategy
*/

package strategy

import (
	"config"
	"db"
	"logger"
	"time"
	"util"
)

const use_time_weighted_algorithm = true
const enable_match = false

var unit_min_amount, unit_max_amount float64
var tickerFSM map[string]int

func init() {
	tickerFSM = make(map[string]int)
}

func TaskTemplate(seconds time.Duration, f func()) {
	ticker := time.NewTicker(seconds * time.Second) // one second
	defer ticker.Stop()

	for _ = range ticker.C {
		f()
	}
}

func before_enter_task() {
	db.Init_sqlstr(config.Config["sqlconn"])

	unit_min_amount = util.ToFloat(config.Config["unit_min_amount"])
	if unit_min_amount < 0.01 {
		unit_min_amount = 0.01
	}

	unit_max_amount = util.ToFloat(config.Config["unit_max_amount"])
	if unit_max_amount < unit_min_amount {
		unit_max_amount = 1
	}

	logger.Infoln("unit_min_amount,unit_max_amount:", unit_min_amount, unit_max_amount)

	initTotalReadyAmount()
}

func TradeCenter() {

	before_enter_task()
	// tickers, err := db.GetNTickers(3)
	// logger.Infoln(tickers, err)
	// return

	markets := getConfExchanges()
	for i := 0; i < len(markets); i++ {
		exchange := markets[i]

		tickerFSM[exchange] = 0

		_queryExchangeData(exchange)
		go TaskTemplateTicker(2, _queryDepth, exchange)
		go TaskTemplateTrader(2, ProgressFSM, exchange)
	}

	// go TaskTemplate(5, QueryTicker)
	go watchdogTicker()

	if use_time_weighted_algorithm {
		go TaskTemplate(1, ProcessDispathMatch)
	}
}

func ProcessDispathMatch() {
	db.TXWrapper(ProcessTimeWeighted)

	if enable_match {
		db.TXWrapper(ProcessMatchTx)
	}
}

func ProgressFSM(exchange string) (err error) {
	tickerFSM[exchange]++
	if tickerFSM[exchange] >= 10 {
		_queryFund(exchange)
		tickerFSM[exchange] = 0
	}

	db.TXWrapperEx(ProcessReady, exchange)

	db.TXWrapperEx(ProcessOrdered, exchange)
	db.TXWrapperEx(ProcessTimeout, exchange)

	return
}

// monitor via recover
type message struct {
	normal bool                   // true means exit normal, otherwise
	state  map[string]interface{} // goroutine state
}

func worker(mess chan message) {
	defer func() {
		exit_message := message{state: make(map[string]interface{})}
		err := recover()
		if err != nil {
			logger.Errorln("worker recover err:", err)
			exit_message.normal = false
		} else {
			exit_message.normal = true
		}
		mess <- exit_message
	}()

	TaskTemplate(3, QueryTicker)
}

func supervisor(mess chan message) {
	m := <-mess
	switch m.normal {
	case true:
		logger.Errorln("exit normal, nothing serious!")
	case false:
		logger.Errorln("exit abnormal, something went wrong")
	}
}

func watchdogTicker() {
	mess := make(chan message, 1)

	go worker(mess)

	supervisor(mess)
}
