package chbtc

import (
	"config"
	"logger"
	"testing"
	"time"
)

const (
	Chbtc_access_key string = ""
	Chbtc_secret_key string = ""
)

func setup() {
	config.Env = config.Test
	config.Root = "/Users/phil/dev/work/haobtc/trader"
	config.LoadConfig()

	//strategy.Task()
}

func Test_API(t *testing.T) {
	setup()

	api := NewExchange(Chbtc_access_key, Chbtc_secret_key)
	buyId, result, err := api.Buy("1100", "0.011")
	logger.Infoln(buyId, result, err)

	order, result, err := api.GetOrder(buyId)
	logger.Infoln(order, result, err)

	err = api.CancelOrder(order.Id)
	logger.Infoln(err)

	order, result, err = api.GetOrder(buyId)
	logger.Infoln(order, result, err)

	// sellId, result, err := api.Sell("2100", "0.012")
	// logger.Infoln(sellId, result, err)

	// order, result, err = api.GetOrder(sellId)
	// logger.Infoln(order, result, err)
}

func TaskTemplate(seconds time.Duration, f func()) {
	ticker := time.NewTicker(seconds * time.Second) // one second
	defer ticker.Stop()

	for _ = range ticker.C {
		f()
	}
}
