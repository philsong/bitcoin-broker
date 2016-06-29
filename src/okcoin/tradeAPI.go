/*
	SEE DOC:
	TRADE API
	https://www.okcoin.cn/about/rest_api.do
*/

package okcoin

import (
	. "config"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"logger"
	"net/url"
	"sort"
	"strings"
	"util"
)

type OkcoinTrade struct {
	api_key    string
	secret_key string
	errno      int64
}

func NewOkcoinTrade(api_key, secret_key string) *OkcoinTrade {
	w := new(OkcoinTrade)
	w.api_key = api_key
	w.secret_key = secret_key
	return w
}

func (w *OkcoinTrade) createSign(pParams map[string]string) string {
	ms := util.NewMapSorter(pParams)
	sort.Sort(ms)

	v := url.Values{}
	for _, item := range ms {
		v.Add(item.Key, item.Val)
	}

	to_sign_str := v.Encode()

	//v.Add("secret_key", w.secret_key)

	h := md5.New()

	raw_str := v.Encode()

	raw_str += "&secret_key=" + w.secret_key

	// logger.Infoln("raw_str",raw_str)
	io.WriteString(h, raw_str)
	sign := fmt.Sprintf("%X", h.Sum(nil))

	req_para := to_sign_str + "&sign=" + sign

	return req_para
}

type ErrorMsg struct {
	Result     bool
	Error_code int
}

func (w *OkcoinTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	if !strings.Contains(body, "result") {
		ret = false
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err := doc.Decode(&errorMsg); err == io.EOF {
		logger.Errorln("OkcoinTrade errorMsg:", err, body)
		ret = false
		return
	} else if err != nil {
		logger.Errorln("OkcoinTrade errorMsg:", err, body)
		ret = false
		return
	}

	if errorMsg.Result != true {
		logger.Errorln("OkcoinTrade errorMsg:", errorMsg)
		ret = false
		return
	}
	ret = true
	return
}

//////
type Asset struct {
	Net   string
	Total string
}
type UnionFund struct {
	BTC string
	LTC string
}

type Money struct {
	BTC string
	CNY string
	LTC string
}

type Funds struct {
	Asset     Asset
	Borrow    Money
	Free      Money
	Freezed   Money
	UnionFund UnionFund
}

type Info struct {
	Funds Funds
}

type UserInfo struct {
	Result bool
	Info   Info
}

func (w *OkcoinTrade) GetAccount() (userInfo UserInfo, err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config["ok_api_userinfo"], req_para)
	if err != nil {
		return
	}

	errorMsg, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		logger.Infoln(ret, errorMsg)
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err = doc.Decode(&userInfo); err == io.EOF {
		logger.Debugln(err)
	} else if err != nil {
		logger.Errorln(err)
	}

	return
}

/////

func (w *OkcoinTrade) doTrade(symbol, method, price, amount string) (string, string, error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["symbol"] = symbol
	pParams["type"] = method

	if method != "sell_market" {
		pParams["price"] = price
	}

	if method != "buy_market" {
		pParams["amount"] = amount
	}

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config["ok_api_trade"], req_para)
	if err != nil {
		return "", "", err
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		result := string(body)
		return "", result, nil
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result   bool
		Order_id int64
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Errorln("OkcoinTrade errorMsg:", err, body)
	} else if err != nil {
		logger.Errorln("OkcoinTrade errorMsg:", err, body)
	}

	if m.Result == true {
		return fmt.Sprintf("%d", m.Order_id), "", nil
	} else {
		err = errors.New(string(body))
		return "", "", err
	}
}

func (w *OkcoinTrade) BuyBTC(price, amount string) (string, string, error) {
	return w.doTrade("btc_cny", "buy", price, amount)
}

func (w *OkcoinTrade) SellBTC(price, amount string) (string, string, error) {
	return w.doTrade("btc_cny", "sell", price, amount)
}

func (w *OkcoinTrade) BuyLTC(price, amount string) (string, string, error) {
	return w.doTrade("ltc_cny", "buy", price, amount)
}

func (w *OkcoinTrade) SellLTC(price, amount string) (string, string, error) {
	return w.doTrade("ltc_cny", "sell", price, amount)
}

func (w *OkcoinTrade) BuyMarketBTC(cny string) (string, string, error) {
	return w.doTrade("btc_cny", "buy_market", cny, "")
}

func (w *OkcoinTrade) SellMarketBTC(btc string) (string, string, error) {
	return w.doTrade("btc_cny", "sell_market", "", btc)
}

func (w *OkcoinTrade) BuyMarketLTC(cny string) (string, string, error) {
	return w.doTrade("ltc_cny", "buy_market", cny, "")
}

func (w *OkcoinTrade) SellMarketLTC(ltc string) (string, string, error) {
	return w.doTrade("ltc_cny", "sell_market", "", ltc)
}

/////
type OKOrder struct {
	Amount      float64
	Avg_price   float64
	Create_date int
	Deal_amount float64
	Order_id    int64
	Orders_id   int64
	Price       float64
	Status      int
	Symbol      string
	Type        string
}

type OKOrderTable struct {
	Result bool
	Orders []OKOrder
}

func (w *OkcoinTrade) Get_order(symbol string, order_id string) (m OKOrderTable, result string, err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["symbol"] = symbol
	pParams["order_id"] = order_id

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config["ok_api_order_info"], req_para)
	if err != nil {
		return
	}

	result = string(body)
	_, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		logger.Errorln("Get_order check_json_result:", order_id, body)

		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err = doc.Decode(&m); err == io.EOF {
		logger.Errorln(err)
	} else if err != nil {
		logger.Errorln(err)
		logger.Errorln(body)
		logger.Errorln(m)
	}

	return
}

func (w *OkcoinTrade) Get_BTCorder(order_id string) (m OKOrderTable, result string, err error) {
	return w.Get_order("btc_cny", order_id)
}

func (w *OkcoinTrade) Get_LTCorder(order_id string) (m OKOrderTable, result string, err error) {
	return w.Get_order("ltc_cny", order_id)
}

func (w *OkcoinTrade) Cancel_order(symbol string, order_id string) (err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["symbol"] = symbol
	pParams["order_id"] = order_id

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config["ok_api_cancelorder"], req_para)
	if err != nil {
		return
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		logger.Errorln("cancel check_json_result:", order_id, body)
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result   bool
		Order_id int64
	}

	var m Msg
	if err := doc.Decode(&m); err == io.EOF {
		logger.Errorln("cancel decode eof:", order_id, body)
	} else if err != nil {
		logger.Errorln("cancel decode err:", order_id, body)
	}

	logger.Debugln(m)

	if m.Result == true {
		logger.Infoln(m)
		return nil
	} else {
		logger.Infoln(m)
		err = errors.New(string(body))

		return
	}
}

func (w *OkcoinTrade) Cancel_BTCorder(order_id string) (err error) {
	return w.Cancel_order("btc_cny", order_id)
}

func (w *OkcoinTrade) Cancel_LTCorder(order_id string) (err error) {
	return w.Cancel_order("ltc_cny", order_id)
}
