package common

const DEPTH = 150

// trade interface type and method
type Account struct {
	Available_cny float64
	Available_btc float64
	Available_ltc float64
	Frozen_cny    float64
	Frozen_btc    float64
	Frozen_ltc    float64
}

type Ticker struct {
	Date   string
	Ticker TickerPrice
}

type TickerPrice struct {
	Buy  float64
	High float64
	Last float64
	Low  float64
	Sell float64
	Vol  float64
}

type MarketOrder struct {
	Price  float64 // 价格
	Amount float64 // 委单量
}

// price from high to low: asks[0] > .....>asks[DEPTH] > bids[0] > ......> bids[DEPTH]
type OrderBook struct {
	Asks [DEPTH]MarketOrder // sell
	Bids [DEPTH]MarketOrder // buy
}

const (
	ORDER_STATE_PENDING  string = "PENDING"
	ORDER_STATE_SUCCESS  string = "SUCCESS"
	ORDER_STATE_CANCELED string = "CANCELED"
	ORDER_STATE_ERROR    string = "ERROR"
	ORDER_STATE_UNKNOWN  string = "UNKNOWN"
)

type Order struct {
	Id          string
	Price       float64
	Amount      float64 //因为火币在buy order中返回的与ok,chbtc不同，我们common统一处理为比特币金额，注意：坑！ok的买单，amount返回为空
	Deal_amount float64 //因为火币在buy order中返回的与ok,chbtc不同，我们common统一处理为比特币成交金额
	Status      string
}

type TradeAPI interface {
	GetTicker() (Ticker, error)
	GetAccount() (Account, error)
	GetDepth() (OrderBook, error)
	Buy(price, amount string) (order_id string, result string, err error)
	Sell(price, amount string) (order_id string, result string, err error)

	BuyMarket(amount string) (order_id string, result string, err error)
	SellMarket(amount string) (order_id string, result string, err error)
	GetOrder(order_id string) (order Order, result string, err error)
	CancelOrder(order_id string) error
}

//////////////////
// below for test
//////////////////
var g_MockObjs map[string]TradeAPI = make(map[string]TradeAPI)

func GetMockTradeAPI(exchange string) TradeAPI {
	return g_MockObjs[exchange]
}

func SetMockTradeAPI(exchange string, mockObj TradeAPI) {
	g_MockObjs[exchange] = mockObj
}
