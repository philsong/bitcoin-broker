package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"trade_service"

	"github.com/apache/thrift/lib/go/thrift"
)

const (
	OKCoin_api_key    string = ""
	OKCoin_secret_key string = ""
	Huobi_access_key  string = ""
	Huobi_secret_key  string = ""
	Chbtc_access_key  string = ""
	Chbtc_secret_key  string = ""
)

func main() {
	startTime := currentTimeMillis()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "19090"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}

	useTransport := transportFactory.GetTransport(transport)
	client := trade_service.NewTradeServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:19090", " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	// for i := 0; i < 1000; i++ {
	// 	paramMap := make(map[string]string)
	// 	paramMap["name"] = "qinerg"
	// 	paramMap["passwd"] = "123456"
	// 	r1, e1 := client.FunCall(currentTimeMillis(), "login", paramMap)
	// 	fmt.Println(i, "Call->", r1, e1)
	// }
	if true {
		var configs []*trade_service.ExchangeConfig
		hb_config := new(trade_service.ExchangeConfig)

		hb_config.Exchange = "huobi"
		hb_config.AccessKey = Huobi_access_key
		hb_config.SecretKey = Huobi_secret_key
		configs = append(configs, hb_config)

		ok_config := new(trade_service.ExchangeConfig)
		ok_config.Exchange = "okcoin"
		ok_config.AccessKey = OKCoin_api_key
		ok_config.SecretKey = OKCoin_secret_key
		configs = append(configs, ok_config)

		ch_config := new(trade_service.ExchangeConfig)
		ch_config.Exchange = "chbtc"
		ch_config.AccessKey = Chbtc_access_key
		ch_config.SecretKey = Chbtc_secret_key
		configs = append(configs, ch_config)

		fmt.Println("configs->", configs)
		err = client.Config(configs)
		fmt.Println("Config<-", err)
	}

	accounts, err := client.GetAccount()
	fmt.Println("accounts->", err, accounts)

	ticker, err := client.GetTicker()
	fmt.Println("ticker->", err, ticker)

	for i := 0; i < 1; i++ {
		ticker, err := client.GetTicker()
		fmt.Println("ticker->", err, ticker)

		if true {
			var buyOrder trade_service.Trade
			buyOrder.ClientID = "10"
			buyOrder.Amount = 50
			order, err := client.Buy(&buyOrder)
			fmt.Println("buy->", err, order)
			//time.Sleep(1 * time.Second)

			order, err = client.GetOrder(order.ID)
			fmt.Println("buy order->", err, order)
		}

		if true {
			var sellOrder trade_service.Trade
			sellOrder.ClientID = "1"
			sellOrder.Amount = 0.015
			order, err := client.Sell(&sellOrder)
			fmt.Println("Sell->", err, order)
			//time.Sleep(1 * time.Second)

			order, err = client.GetOrder(order.ID)
			fmt.Println("Sell order->", err, order)
		}
	}

	accounts, err = client.GetAccount()
	fmt.Println("accounts->", err, accounts)

	alertOrders, err := client.GetAlertOrders()
	fmt.Println("errorOrders->", err, len(alertOrders))

	// for i := 0; i < 10; i++ {
	// 	order, err := client.Buy("1000000")
	// 	fmt.Println("Buy->", err, order)
	// 	time.Sleep(3 * time.Second)
	// }

	endTime := currentTimeMillis()
	fmt.Println("Program exit. time->", endTime, startTime, (endTime - startTime))
}

// 转换成毫秒
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
