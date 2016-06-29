package db

import (
	"github.com/jinzhu/gorm"
	"logger"
	"time"
	"trade_service"
)

type DBTradeOrder struct {
	ID            int64     `thrift:"id,1" json:"id"`
	SiteOrderID   int64     `thrift:"site_order_id,2" json:"site_order_id"`
	Exchange      string    `thrift:"exchange,3" json:"exchange"`
	Price         float64   `thrift:"price,4" json:"price"`
	TradeType     string    `thrift:"trade_type,5" json:"trade_type"`
	OrderStatus   string    `thrift:"order_status,6" json:"order_status"`
	EstimateCny   float64   `thrift:"estimate_cny,7" json:"estimate_cny"`
	EstimateBtc   float64   `thrift:"estimate_btc,8" json:"estimate_btc"`
	EstimatePrice float64   `thrift:"estimate_price,9" json:"estimate_price"`
	DealCny       float64   `thrift:"deal_cny,10" json:"deal_cny"`
	DealBtc       float64   `thrift:"deal_btc,11" json:"deal_btc"`
	DealPrice     float64   `thrift:"deal_price,12" json:"deal_price"`
	PriceMargin   float64   `thrift:"price_margin,13" json:"price_margin"`
	OrderID       string    `thrift:"order_id,14" json:"order_id"`
	Created       time.Time `thrift:"created,15" json:"created"`
	UpdateAt      time.Time `thrift:"update_at,16" json:"update_at"`
	TryTimes      int64     `thrift:"try_times,17" json:"try_times"`
	Info          string    `thrift:"info,18" json:"info"`
	Memo          string    `thrift:"memo,19" json:"memo"`
	MatchID       int64     `thrift:"match_id,20" json:"match_id"`
}

func (DBTradeOrder) TableName() string {
	return "trade_order"
}

func InsertTradeOrder(tx *gorm.DB, tradeOrder *trade_service.TradeOrder) (id int64, err error) {
	dbTradeOrder := convert_to_db_trade_order(true, tradeOrder)
	dbTradeOrder.ID = 0

	if err = tx.Create(&dbTradeOrder).Error; err != nil {
		logger.Errorln("InsertOrder err:", err)
		return
	}

	*tradeOrder, err = convert_to_trade_order(dbTradeOrder)
	if err != nil {
		return
	}

	return tradeOrder.ID, nil
}

func UpdateTradeOrder(tx *gorm.DB, tradeOrder trade_service.TradeOrder) (err error) {

	var dbTradeOrder DBTradeOrder
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", tradeOrder.ID).First(&dbTradeOrder).Error; err != nil {
		logger.Errorln("UpdateTradeOrder not found:", err)
		return err
	}

	newDbTradeOrder := convert_to_db_trade_order(false, &tradeOrder)

	if err = tx.Save(&newDbTradeOrder).Error; err != nil {
		logger.Errorln("UpdateTradeOrder Save err:", err, newDbTradeOrder)
		return
	}

	return
}

func GetAlertOrders() (tradeOrders []trade_service.TradeOrder, err error) {
	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	return GetAllTradeOrders(db, trade_service.OrderStatus_ERROR)
}

func GetTotalReadyNow() (buy_total, sell_total float64) {
	buy_total = 0
	sell_total = 0

	db, err := getORMDB()
	if err != nil {
		logger.Errorln(err)
		return
	}

	type Result struct {
		TradeType string
		SellTotal float64
		BuyTotal  float64
	}

	var results []Result
	db.Table("trade_order").Select("trade_type, sum(estimate_btc) as sell_total, sum(estimate_cny) as buy_total").Group("trade_type").Where("order_status in (?)", []string{"READY", "TIME_WEIGHTED"}).Scan(&results)
	for _, result := range results {
		if result.TradeType == "BUY" {
			buy_total = result.BuyTotal
		} else if result.TradeType == "SELL" {
			sell_total = result.SellTotal
		}
	}

	return
}

func GetAllTradeOrders(tx *gorm.DB, order_status trade_service.OrderStatus) (tradeOrders []trade_service.TradeOrder, err error) {
	var dbTradeOrders []DBTradeOrder
	if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("order_status = ?", order_status.String()).Find(&dbTradeOrders).Error; err != nil {
		logger.Errorln("GetAllTradeOrders err:", order_status.String(), err)
		return
	}

	return convert_to_trade_orders(dbTradeOrders)
}

func GetTradeOrders(tx *gorm.DB, exchange string, order_status trade_service.OrderStatus) (tradeOrders []trade_service.TradeOrder, err error) {
	var dbTradeOrders []DBTradeOrder
	if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("order_status = ? AND exchange = ?", order_status.String(), exchange).Find(&dbTradeOrders).Error; err != nil {
		logger.Errorln("GetTradeOrders err:", order_status.String(), exchange, err)
		return
	}

	return convert_to_trade_orders(dbTradeOrders)
}

func GetFirstTradeOrder(tx *gorm.DB, exchange string, order_status trade_service.OrderStatus) (tradeOrders []trade_service.TradeOrder, err error) {
	var dbTradeOrders []DBTradeOrder
	if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("order_status = ? AND exchange = ?", order_status.String(), exchange).Order("try_times desc").Limit(1).Find(&dbTradeOrders).Error; err != nil {
		logger.Errorln("GetFirstTradeOrder err:", order_status.String(), exchange, err)
		return
	}

	return convert_to_trade_orders(dbTradeOrders)
}

func GetTimeoutOrder(tx *gorm.DB, exchange string) (tradeOrders []trade_service.TradeOrder, err error) {
	var dbTradeOrders []DBTradeOrder
	if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("order_status = 'ORDERED' AND exchange = ?", exchange).Order("update_at desc").Limit(10).Find(&dbTradeOrders).Error; err != nil {
		logger.Errorln("GetTimeoutOrder err:", exchange, err)
		return
	}

	var final_dbTradeOrders []DBTradeOrder
	for _, dbTradeOrder := range dbTradeOrders {
		now_time := time.Now()
		diff := now_time.Sub(dbTradeOrder.UpdateAt).Seconds()

		logger.Debugln("time:diff", diff, dbTradeOrder.UpdateAt.Format(time.RFC3339), now_time.Format(time.RFC3339))
		if diff < 15 {
			break
		}

		final_dbTradeOrders = append(final_dbTradeOrders, dbTradeOrder)
	}

	logger.Debugln("final_dbTradeOrders:", final_dbTradeOrders)

	return convert_to_trade_orders(final_dbTradeOrders)
}

//remove ID and Created
func convert_to_db_trade_order(isInsert bool, tradeOrder *trade_service.TradeOrder) (dbTradeOrder DBTradeOrder) {
	dbTradeOrder.ID = tradeOrder.ID
	dbTradeOrder.SiteOrderID = tradeOrder.SiteOrderID
	dbTradeOrder.Exchange = tradeOrder.Exchange
	dbTradeOrder.Price = tradeOrder.Price
	dbTradeOrder.TradeType = tradeOrder.TradeType.String()
	dbTradeOrder.OrderStatus = tradeOrder.OrderStatus.String()
	dbTradeOrder.EstimateCny = tradeOrder.EstimateCny
	dbTradeOrder.EstimateBtc = tradeOrder.EstimateBtc
	dbTradeOrder.EstimatePrice = tradeOrder.EstimatePrice
	dbTradeOrder.DealCny = tradeOrder.DealCny
	dbTradeOrder.DealBtc = tradeOrder.DealBtc
	dbTradeOrder.DealPrice = tradeOrder.DealPrice
	dbTradeOrder.PriceMargin = tradeOrder.PriceMargin
	dbTradeOrder.OrderID = tradeOrder.OrderID

	if isInsert {
		dbTradeOrder.Created = time.Now()
		dbTradeOrder.UpdateAt = time.Now()
	} else {
		created_at, err := time.Parse(
			time.RFC3339,
			tradeOrder.Created)
		if err == nil {
			dbTradeOrder.Created = created_at
		} else {
			logger.Errorln("Parse Created time.RFC3339 error:", tradeOrder.Created)
		}

		update_at, err := time.Parse(
			time.RFC3339,
			tradeOrder.UpdateAt)
		if err == nil {
			dbTradeOrder.UpdateAt = update_at
		} else {
			logger.Errorln("Parse UpdateAt time.RFC3339 error:", tradeOrder.UpdateAt)
		}
	}

	dbTradeOrder.TryTimes = tradeOrder.TryTimes
	dbTradeOrder.Info = tradeOrder.Info
	dbTradeOrder.Memo = tradeOrder.Memo
	dbTradeOrder.MatchID = tradeOrder.MatchID

	return
}

func convert_to_trade_order(dbTradeOrder DBTradeOrder) (tradeOrder trade_service.TradeOrder, err error) {
	tradeOrder.ID = dbTradeOrder.ID
	tradeOrder.SiteOrderID = dbTradeOrder.SiteOrderID
	tradeOrder.Exchange = dbTradeOrder.Exchange
	tradeOrder.Price = dbTradeOrder.Price
	tradeOrder.TradeType, err = trade_service.TradeTypeFromString(dbTradeOrder.TradeType)
	if err != nil {
		logger.Errorln("TradeType panic", err)
		return
	}
	tradeOrder.OrderStatus, err = trade_service.OrderStatusFromString(dbTradeOrder.OrderStatus)
	if err != nil {
		logger.Errorln("OrderStatus panic", err)
		return
	}
	tradeOrder.EstimateCny = dbTradeOrder.EstimateCny
	tradeOrder.EstimateBtc = dbTradeOrder.EstimateBtc
	tradeOrder.EstimatePrice = dbTradeOrder.EstimatePrice
	tradeOrder.DealCny = dbTradeOrder.DealCny
	tradeOrder.DealBtc = dbTradeOrder.DealBtc
	tradeOrder.DealPrice = dbTradeOrder.DealPrice
	tradeOrder.PriceMargin = dbTradeOrder.PriceMargin
	tradeOrder.OrderID = dbTradeOrder.OrderID
	tradeOrder.Created = dbTradeOrder.Created.Format(time.RFC3339)
	tradeOrder.UpdateAt = dbTradeOrder.UpdateAt.Format(time.RFC3339)
	tradeOrder.TryTimes = dbTradeOrder.TryTimes
	tradeOrder.Info = dbTradeOrder.Info
	tradeOrder.Memo = dbTradeOrder.Memo
	tradeOrder.MatchID = dbTradeOrder.MatchID
	return
}

func convert_to_trade_orders(dbTradeOrders []DBTradeOrder) (tradeOrders []trade_service.TradeOrder, err error) {
	for _, dbTradeOrder := range dbTradeOrders {
		var tradeOrder trade_service.TradeOrder
		tradeOrder, err = convert_to_trade_order(dbTradeOrder)
		if err != nil {
			logger.Errorln("convert_to_trade_order err", err)
			return
		}

		tradeOrders = append(tradeOrders, tradeOrder)
	}

	return
}
