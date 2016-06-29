package huobi

import (
	"config"
	"logger"
	"testing"
	"time"
)

const (
	Huobi_access_key string = ""
	Huobi_secret_key string = ""
)

func setup() {
	config.Env = config.Test
	config.Root = "/Users/phil/dev/work/haobtc/trader"
	config.LoadConfig()

	//strategy.Task()
}

func Test_API(t *testing.T) {
	setup()

	api := NewExchange(Huobi_access_key, Huobi_secret_key)

	if false {

		buyId, result, err := api.Buy("1100", "0.001")
		logger.Infoln(buyId, result, err)

		order, result, err := api.GetOrder(buyId)
		logger.Infoln(order, result, err)

		sellId, result, err := api.Sell("2100", "0.002")
		logger.Infoln(sellId, result, err)

		order, result, err = api.GetOrder(sellId)
		logger.Infoln(order, result, err)
	}

	if true {

		buyId, result, err := api.BuyMarket("1")
		logger.Infoln(buyId, result, err)

		order, result, err := api.GetOrder(buyId)
		logger.Infoln(order, result, err)

		sellId, result, err := api.SellMarket("0.0022")
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
