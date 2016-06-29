/*
  trader API Engine
*/

package huobi

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
	"strconv"
	"strings"
	"time"
	"util"
)

/*
	https://www.huobi.com/help/index.php?a=api_help_v3
*/
type HuobiTrade struct {
	access_key string
	secret_key string
}

func NewHuobiTrade(access_key, secret_key string) *HuobiTrade {
	w := new(HuobiTrade)
	w.access_key = access_key
	w.secret_key = secret_key
	return w
}

func (w *HuobiTrade) createSign(pParams map[string]string) string {
	pParams["secret_key"] = w.secret_key

	ms := util.NewMapSorter(pParams)
	sort.Sort(ms)

	v := url.Values{}
	for _, item := range ms {
		v.Add(item.Key, item.Val)
	}

	to_sign_str := v.Encode()

	h := md5.New()

	io.WriteString(h, v.Encode())

	sign := fmt.Sprintf("%x", h.Sum(nil))

	req_para := to_sign_str + "&sign=" + sign

	return req_para
}

type ErrorMsg struct {
	Code    string
	Message string
}

func (w *HuobiTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	if !strings.Contains(body, "code") {
		ret = true
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err := doc.Decode(&errorMsg); err == io.EOF {
		logger.Errorln("HuobiTrade errorMsg:", err, body)
		ret = false
		return
	} else if err != nil {
		logger.Errorln("HuobiTrade errorMsg:", err, body)
		ret = false
		return
	}

	if errorMsg.Code != "0" {
		logger.Errorln("HuobiTrade errorMsg:", errorMsg)
		ret = false
		return
	}

	ret = true
	return
}

type Account_info struct {
	Total                 string
	Net_asset             string
	Available_cny_display string
	Available_btc_display string
	Available_ltc_display string
	Frozen_cny_display    string
	Frozen_btc_display    string
	Frozen_ltc_display    string
	Loan_cny_display      string
	Loan_btc_display      string
	Loan_ltc_display      string
}

func (w *HuobiTrade) GetAccount() (account_info Account_info, err error) {
	pParams := make(map[string]string)
	pParams["method"] = "get_account_info"
	pParams["access_key"] = w.access_key
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)

	req_para := w.createSign(pParams)
	body, err := util.HttpPost(Config["hb_api_url"], req_para)
	if err != nil {
		return
	}

	_, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	if err := doc.Decode(&account_info); err == io.EOF {
		logger.Errorln(err)
	} else if err != nil {
		logger.Errorln(err)
	}

	return
}

func (w *HuobiTrade) _doTrade(method, cointype, price, amount string) (string, string, error) {
	pParams := make(map[string]string)
	pParams["method"] = method
	pParams["access_key"] = w.access_key
	pParams["coin_type"] = cointype
	if method == "buy" || method == "sell" {
		pParams["price"] = price
	}

	pParams["amount"] = amount
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)
	req_para := w.createSign(pParams)
	body, err := util.HttpPost(Config["hb_api_url"], req_para)
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
		Result string
		Id     int64
	}

	var m Msg
	if err = doc.Decode(&m); err == io.EOF {
		logger.Errorln("HuobiTrade errorMsg:", err, body)
	} else if err != nil {
		logger.Errorln("HuobiTrade errorMsg:", err, body)
	}

	if m.Result == "success" {
		return fmt.Sprintf("%d", m.Id), "", nil
	} else {
		err = errors.New(string(body))
		return "", "", err
	}
}

func (w *HuobiTrade) doTrade(method, symbol, price, amount string) (string, string, error) {
	var cointype string
	if symbol == "btc_cny" {
		cointype = "1"
	} else if symbol == "ltc_cny" {
		cointype = "0"
	} else {
		panic(-1)
	}

	return w._doTrade(method, cointype, price, amount)
}

type HBOrder struct {
	Id               int64
	Type             int
	Order_price      string
	Order_amount     string
	Processed_price  string
	Processed_amount string
	Vot              string
	Fee              string
	Total            string
	Status           int
}

func (w *HuobiTrade) Get_order(cointype string, order_id string) (m HBOrder, result string, err error) {
	pParams := make(map[string]string)
	pParams["method"] = "order_info"
	pParams["access_key"] = w.access_key
	pParams["coin_type"] = cointype
	pParams["id"] = order_id
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)

	req_para := w.createSign(pParams)
	body, err := util.HttpPost(Config["hb_api_url"], req_para)
	if err != nil {
		return
	}

	result = string(body)
	_, ret := w.check_json_result(body)
	if ret == false {
		logger.Infoln(body)
		err = errors.New(string(body))
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))
	if err = doc.Decode(&m); err == io.EOF {
		logger.Infoln(err)
	} else if err != nil {
		logger.Infoln(err)
	}

	return
}

func (w *HuobiTrade) Cancel_order(cointype string, order_id string) (err error) {
	pParams := make(map[string]string)
	pParams["method"] = "cancel_order"
	pParams["access_key"] = w.access_key
	pParams["coin_type"] = cointype
	pParams["id"] = order_id
	now := time.Now().Unix()
	pParams["created"] = strconv.FormatInt(now, 10)

	req_para := w.createSign(pParams)
	body, err := util.HttpPost(Config["hb_api_url"], req_para)
	if err != nil {
		return
	}
	_, ret := w.check_json_result(body)
	if ret == false {
		err = errors.New(string(body))
		return
	}

	doc := json.NewDecoder(strings.NewReader(body))

	type Msg struct {
		Result string
	}

	var m Msg
	if err = doc.Decode(&m); err == io.EOF {
		logger.Errorln(err)
	} else if err != nil {
		logger.Errorln(err)
	}

	logger.Debugln(m)

	if m.Result == "success" {
		return nil
	} else {
		err = errors.New(string(body))
		return
	}
}
