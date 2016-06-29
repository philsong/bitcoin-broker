package mocks

import "common"
import "github.com/stretchr/testify/mock"

type TradeAPI struct {
	mock.Mock
}

func (m *TradeAPI) GetTicker() (common.Ticker, error) {
	ret := m.Called()

	r0 := ret.Get(0).(common.Ticker)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *TradeAPI) GetAccount() (common.Account, error) {
	ret := m.Called()

	r0 := ret.Get(0).(common.Account)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *TradeAPI) GetDepth() (common.OrderBook, error) {
	ret := m.Called()

	r0 := ret.Get(0).(common.OrderBook)
	r1 := ret.Error(1)

	return r0, r1
}
func (m *TradeAPI) Buy(price string, amount string) (string, string, error) {
	ret := m.Called(price, amount)

	r0 := ret.Get(0).(string)
	r1 := ret.Get(1).(string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *TradeAPI) Sell(price string, amount string) (string, string, error) {
	ret := m.Called(price, amount)

	r0 := ret.Get(0).(string)
	r1 := ret.Get(1).(string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *TradeAPI) BuyMarket(amount string) (string, string, error) {
	ret := m.Called(amount)

	r0 := ret.Get(0).(string)
	r1 := ret.Get(1).(string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *TradeAPI) SellMarket(amount string) (string, string, error) {
	ret := m.Called(amount)

	r0 := ret.Get(0).(string)
	r1 := ret.Get(1).(string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *TradeAPI) GetOrder(order_id string) (common.Order, string, error) {
	ret := m.Called(order_id)

	r0 := ret.Get(0).(common.Order)
	r1 := ret.Get(1).(string)
	r2 := ret.Error(2)

	return r0, r1, r2
}
func (m *TradeAPI) CancelOrder(order_id string) error {
	ret := m.Called(order_id)

	r0 := ret.Error(0)

	return r0
}
