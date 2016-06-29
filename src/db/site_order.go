package db

import (
	"github.com/jinzhu/gorm"
	"logger"
	"time"
	"trade_service"
)

type SiteOrder struct {
	ID            int64                     `json:"id"`
	ClientID      string                    `json:"client_id"`
	TradeType     trade_service.TradeType   `json:"trade_type"`
	OrderStatus   trade_service.OrderStatus `json:"order_status"`
	Amount        float64                   `json:"amount"`
	Price         float64                   `json:"price"`
	EstimatePrice float64                   `json:"estimate_price"`
	EstimateCny   float64                   `json:"estimate_cny"`
	EstimateBtc   float64                   `json:"estimate_btc"`
	Created       time.Time                 `json:"created"`
}

type DBSiteOrder struct {
	ID            int64     `json:"id"`
	ClientID      string    `json:"client_id"`
	TradeType     string    `json:"trade_type"`
	OrderStatus   string    `json:"order_status"`
	Amount        float64   `json:"amount"`
	Price         float64   `json:"price"`
	EstimatePrice float64   `json:"estimate_price"`
	EstimateCny   float64   `json:"estimate_cny"`
	EstimateBtc   float64   `json:"estimate_btc"`
	Created       time.Time `json:"created"`
}

func (DBSiteOrder) TableName() string {
	return "site_order"
}

func InsertOrder(tx *gorm.DB, order *SiteOrder) (order_id int64, err error) {
	dbSiteOrder := _convert_to_db_site_order(order)
	if err = tx.Create(&dbSiteOrder).Error; err != nil {
		logger.Errorln("InsertOrder err:", err)
		return
	}

	logger.Infoln("InsertOrder ok", dbSiteOrder)

	order.ID = dbSiteOrder.ID

	return order.ID, nil
}

func GetOrderByClientID(tx *gorm.DB, client_id string) (orders []SiteOrder, err error) {
	var dbSiteOrders []DBSiteOrder
	if err = tx.Where("client_id = ?", client_id).Find(&dbSiteOrders).Error; err != nil {
		logger.Errorln("GetOrderByClientID err:", err, client_id)
		return
	}

	return convert_to_site_orders(dbSiteOrders)
}

func GetOrdersByStatus(tx *gorm.DB, orderStatus trade_service.OrderStatus) (orders []SiteOrder, err error) {
	var dbSiteOrders []DBSiteOrder
	if err = tx.Where("order_status = ?", orderStatus.String()).Find(&orders).Error; err != nil {
		logger.Errorln("GetOrdersByStatus err:", err, orderStatus.String())
		return
	}

	return convert_to_site_orders(dbSiteOrders)
}

//remove ID and Created
func _convert_to_db_site_order(order *SiteOrder) (dbSiteOrder DBSiteOrder) {
	dbSiteOrder.ClientID = order.ClientID
	dbSiteOrder.TradeType = order.TradeType.String()
	dbSiteOrder.OrderStatus = order.OrderStatus.String()
	dbSiteOrder.Amount = order.Amount
	dbSiteOrder.Price = order.Price
	dbSiteOrder.EstimatePrice = order.EstimatePrice
	dbSiteOrder.EstimateCny = order.EstimateCny
	dbSiteOrder.EstimateBtc = order.EstimateBtc
	dbSiteOrder.Created = time.Now()

	return
}

func convert_to_site_order(dbSiteOrder DBSiteOrder) (order SiteOrder, err error) {
	order.ID = dbSiteOrder.ID
	order.ClientID = dbSiteOrder.ClientID
	order.TradeType, err = trade_service.TradeTypeFromString(dbSiteOrder.TradeType)
	if err != nil {
		logger.Errorln("TradeType panic", err)
		return
	}
	order.OrderStatus, err = trade_service.OrderStatusFromString(dbSiteOrder.OrderStatus)
	if err != nil {
		logger.Errorln("OrderStatus panic", err)
		return
	}
	order.Amount = dbSiteOrder.Amount
	order.Price = dbSiteOrder.Price
	order.EstimatePrice = dbSiteOrder.EstimatePrice
	order.EstimateCny = dbSiteOrder.EstimateCny
	order.EstimateBtc = dbSiteOrder.EstimateBtc
	order.Created = dbSiteOrder.Created

	return
}

func convert_to_site_orders(dbSiteOrders []DBSiteOrder) (orders []SiteOrder, err error) {
	for _, dbSiteOrder := range dbSiteOrders {
		var order SiteOrder
		order, err = convert_to_site_order(dbSiteOrder)
		if err != nil {
			logger.Errorln("convert_to_site_order err", err)
			return
		}

		orders = append(orders, order)
	}

	return
}
