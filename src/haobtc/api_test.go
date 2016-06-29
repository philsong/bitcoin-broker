package haobtc

import (
	"config"
	"logger"
	"testing"
	"time"
)

const (
	OKCoin_api_key    string = ""
	OKCoin_secret_key string = ""
)

func setup() {
	config.Env = config.Test
	config.Root = "/Users/phil/dev/work/haobtc/trader"
	config.LoadConfig()

	//strategy.Task()
}

func Test_API(t *testing.T) {
	setup()

	api := NewExchange(OKCoin_api_key, OKCoin_secret_key)

	if false {

		buyId, result, err := api.Buy("1100", "0.011")
		logger.Infoln(buyId, result, err)

		order, result, err := api.GetOrder(buyId)
		logger.Infoln(order, result, err)

		sellId, result, err := api.Sell("2100", "0.012")
		logger.Infoln(sellId, result, err)

		order, result, err = api.GetOrder(sellId)
		logger.Infoln(order, result, err)
	}

	if true {
		buyId, result, err := api.BuyMarket("30")
		logger.Infoln(buyId, result, err)

		order, result, err := api.GetOrder(buyId)
		logger.Infoln(order, result, err)

		sellId, result, err := api.SellMarket("0.013")
		logger.Infoln(sellId, result, err)

		order, result, err = api.GetOrder(sellId)
		logger.Infoln(order, result, err)
	}

}

func TaskTemplate(seconds time.Duration, f func()) {
	ticker := time.NewTicker(seconds * time.Second) // one second
	defer ticker.Stop()

	for _ = range ticker.C {
		f()
	}
}
