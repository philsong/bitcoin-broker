package haobtc

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

type HaobtcTrade struct {
	name       string
	api_key    string
	secret_key string
	errno      int64
}

func NewHaobtcTrade(name, api_key, secret_key string) *HaobtcTrade {
	w := new(HaobtcTrade)
	w.name = name
	w.api_key = api_key
	w.secret_key = secret_key
	return w
}

func (w *HaobtcTrade) createSign(pParams map[string]string) string {
	ms := util.NewMapSorter(pParams)
	sort.Sort(ms)

	v := url.Values{}
	for _, item := range ms {
		v.Add(item.Key, item.Val)
	}

	to_sign_str := v.Encode()

	logger.Debugln(to_sign_str)

	h := md5.New()

	raw_str := v.Encode()

	raw_str += "&secret_key=" + w.secret_key

	logger.Debugln("raw_str", raw_str)
	io.WriteString(h, raw_str)
	sign := fmt.Sprintf("%X", h.Sum(nil))

	req_para := to_sign_str + "&sign=" + sign

	return req_para
}

type ErrorMsg struct {
	Code    string
	Message string
}

func (w *HaobtcTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	if !strings.Contains(body, "code") {
		ret = true
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err := doc.Decode(&errorMsg); err == io.EOF {
		logger.Errorln("HaobtcTrade errorMsg:", err, body)
		ret = false
		return
	} else if err != nil {
		logger.Errorln("HaobtcTrade errorMsg:", err, body)
		ret = false
		return
	}

	ret = false
	return
}

//////
type UserInfo struct {
	Exchange_cny        float64
	Exchange_btc        float64
	Wallet_cny          float64
	Wallet_btc          float64
	Exchange_frozen_cny float64
	Exchange_frozen_btc float64
}

func (w *HaobtcTrade) GetAccount() (userInfo UserInfo, err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config[w.name+"_api_userinfo"], req_para)
	if err != nil {
		return
	}

	errorMsg, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		logger.Errorln(ret, body, errorMsg)
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err = doc.Decode(&userInfo); err == io.EOF {
		logger.Errorln(err)
	} else if err != nil {
		logger.Errorln(err)
	}

	return
}

/////

func (w *HaobtcTrade) doTrade(symbol, method, price, amount string) (string, string, error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["type"] = method

	if method != "sell_market" && method != "buy_market" {
		pParams["price"] = price
	}

	pParams["amount"] = amount

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config[w.name+"_api_trade"], req_para)
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
		Order_id int64
	}

	var m Msg
	if err = doc.Decode(&m); err == io.EOF {
		logger.Errorln("HaobtcTrade errorMsg:", err, body)
	} else if err != nil {
		logger.Errorln("HaobtcTrade errorMsg:", err, body)
	}

	logger.Debugln(m)
	if m.Order_id == -1 {
		return "", "balance is not enough", err
	}

	return fmt.Sprintf("%d", m.Order_id), "", err
}

func (w *HaobtcTrade) BuyBTC(price, amount string) (string, string, error) {
	return w.doTrade("btc_cny", "buy", price, amount)
}

func (w *HaobtcTrade) SellBTC(price, amount string) (string, string, error) {
	return w.doTrade("btc_cny", "sell", price, amount)
}

func (w *HaobtcTrade) BuyMarketBTC(cny string) (string, string, error) {
	return w.doTrade("btc_cny", "buy_market", "", cny)
}

func (w *HaobtcTrade) SellMarketBTC(btc string) (string, string, error) {
	return w.doTrade("btc_cny", "sell_market", "", btc)
}

/////
type HaobtcOrder struct {
	Amount      float64
	Avg_price   float64
	Create_date string
	Deal_size   float64
	Order_id    int64
	Price       float64
	Status      string
	Side        string
	Type        string
}

func (w *HaobtcTrade) Get_order(symbol string, order_id string) (m HaobtcOrder, result string, err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["order_id"] = order_id

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config[w.name+"_api_order_info"], req_para)
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

func (w *HaobtcTrade) Get_BTCorder(order_id string) (m HaobtcOrder, result string, err error) {
	return w.Get_order("btc_cny", order_id)
}

func (w *HaobtcTrade) Cancel_order(symbol string, order_id string) (err error) {
	pParams := make(map[string]string)
	pParams["api_key"] = w.api_key
	pParams["order_id"] = order_id

	req_para := w.createSign(pParams)

	body, err := util.HttpPost(Config[w.name+"_api_cancelorder"], req_para)
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
		Order_id int64
	}

	var m Msg
	if err = doc.Decode(&m); err == io.EOF {
		logger.Errorln("cancel decode eof:", order_id, body)
	} else if err != nil {
		logger.Errorln("cancel decode err:", order_id, body)
	}

	logger.Debugln(m)

	return
}

func (w *HaobtcTrade) Cancel_BTCorder(order_id string) (err error) {
	return w.Cancel_order("btc_cny", order_id)
}
