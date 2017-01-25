/*
  trader API Engine
*/

package chbtc

import (
	"common"
	. "config"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"logger"
	"strconv"
	"time"
	"util"

	. "github.com/bitly/go-simplejson"
)

/*
	https://vip.chbtc.com/u/api
*/
type ChbtcTrade struct {
	access_key string
	secret_key string
}

func NewChbtcTrade(access_key, secret_key string) *ChbtcTrade {
	w := new(ChbtcTrade)
	w.access_key = access_key
	w.secret_key = secret_key
	return w
}

func (w *ChbtcTrade) createSign(method, otherParams string) string {
	to_sign_str := "method=" + method + "&accesskey=" + w.access_key + otherParams

	secret_key := sha1.Sum([]byte(w.secret_key))
	secret_key_hexstr := hex.EncodeToString(secret_key[:])

	hmacMD5 := hmac.New(md5.New, []byte(secret_key_hexstr))
	hmacMD5.Write([]byte(to_sign_str))
	expectedMAC := hmacMD5.Sum(nil)

	sign := hex.EncodeToString(expectedMAC)

	now := time.Now().Unix() * 1000
	req_para := to_sign_str + "&sign=" + sign + "&reqTime=" + strconv.FormatInt(now, 10)

	return req_para
}

type ErrorMsg struct {
	Code    int
	Message string
}

func (w *ChbtcTrade) check_json_result(body string) (errorMsg ErrorMsg, ret bool) {
	//dummy function
	ret = true
	return
}

func (w *ChbtcTrade) GetAccount() (account common.Account, err error) {
	method := "getAccountInfo"

	req_para := w.createSign(method, "")

	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		return
	}

	js, err := NewJson([]byte(body))
	if err != nil {
		return
	}

	if js, ret := js.CheckGet("result"); ret {
		account.Available_cny = js.Get("balance").Get("CNY").Get("amount").MustFloat64()
		account.Available_btc = js.Get("balance").Get("BTC").Get("amount").MustFloat64()
		account.Available_ltc = js.Get("balance").Get("LTC").Get("amount").MustFloat64()
		account.Frozen_cny = js.Get("frozen").Get("CNY").Get("amount").MustFloat64()
		account.Frozen_btc = js.Get("frozen").Get("BTC").Get("amount").MustFloat64()
		account.Frozen_ltc = js.Get("frozen").Get("LTC").Get("amount").MustFloat64()
	}

	return
}

func (w *ChbtcTrade) _doTrade(tradeType, price, amount string) (string, string, error) {
	method := "order"
	otherParams := "&price=" + price + "&amount=" + amount + "&tradeType=" + tradeType + "&currency=btc"
	req_para := w.createSign(method, otherParams)
	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		return "", "", err
	}

	js, err := NewJson([]byte(body))
	if err != nil {
		return "", "", err
	}

	if jscode, ret := js.CheckGet("code"); ret {
		code := jscode.MustInt()

		if code == 1000 {
			return js.Get("id").MustString(), "", nil
		} else {
			result := string(body)
			return "", result, nil
		}
	}

	err = errors.New(string(body))
	return "", "", err
}

func (w *ChbtcTrade) doTrade(tradeType, price, amount string) (string, string, error) {
	var _tradeType string
	if tradeType == "buy" {
		_tradeType = "1"
	} else if tradeType == "sell" {
		_tradeType = "0"
	} else {
		panic(-1)
	}

	return w._doTrade(_tradeType, price, amount)
}

func (w *ChbtcTrade) Get_order(cointype string, id string) (order common.Order, result string, err error) {
	method := "getOrder"
	otherParams := "&id=" + id + "&currency=btc"

	req_para := w.createSign(method, otherParams)
	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		logger.Errorln(err)
		return
	}

	result = string(body)
	js, err := NewJson([]byte(body))
	if err != nil {
		logger.Errorln(err)
		return
	}

	if jsid, ret := js.CheckGet("id"); ret {
		order.Id = jsid.MustString()
		// order.Price = js.Get("price").MustFloat64()
		order.Amount = js.Get("total_amount").MustFloat64()
		order.Deal_amount = js.Get("trade_amount").MustFloat64()
		if order.Deal_amount < 0.000001 {
			order.Price = js.Get("price").MustFloat64()
		} else {
			trade_money := js.Get("trade_money").MustFloat64()
			order.Price = trade_money / order.Deal_amount
		}

		status := js.Get("status").MustInt()

		switch status {
		case 0, 3:
			order.Status = common.ORDER_STATE_PENDING
		case 2:
			order.Status = common.ORDER_STATE_SUCCESS
		case 1:
			order.Status = common.ORDER_STATE_CANCELED
		default:
			order.Status = common.ORDER_STATE_UNKNOWN
		}
		return
	} else if jscode, ret := js.CheckGet("code"); ret {
		jscode := jscode.MustInt()
		if jscode == 3001 {
			order.Status = common.ORDER_STATE_ERROR
			return
		}
	} else {
		err = errors.New("Get_order failed")
		return
	}

	return
}

func (w *ChbtcTrade) Cancel_order(cointype string, id string) (err error) {
	method := "cancelOrder"
	otherParams := "&id=" + id + "&currency=btc"

	req_para := w.createSign(method, otherParams)
	fmt.Println(req_para)
	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		return
	}

	js, err := NewJson([]byte(body))
	if err != nil {
		return
	}

	if jscode, ret := js.CheckGet("code"); ret {
		code := jscode.MustInt()

		if code == 1000 {
			return nil
		} else {
			logger.Errorln("code=", code)
			return errors.New("Cancel_order failed")
		}
	} else {
		return errors.New("Cancel_order failed")
	}
}

func (w *ChbtcTrade) withdraw(amount, fees, receiveAddr, safePwd string) (id string, err error) {
	method := "withdraw"
	otherParams := "&currency=btc_cny&fees=" + fees + "&amount=" + amount + "&receiveAddr=" + receiveAddr + "&safePwd=" + safePwd

	req_para := w.createSign(method, otherParams)
	fmt.Println(req_para)
	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		return
	}

	js, err := NewJson([]byte(body))
	if err != nil {
		return
	}
	logger.Infoln(js)
	logger.Errorln(err)

	if jscode, ret := js.CheckGet("code"); ret {
		code := jscode.MustInt()

		if code == 1000 {
			return js.Get("id").MustString(), nil
		} else {
			logger.Errorln("code=", code)
			return "", errors.New("withdraw failed")
		}
	} else {
		return "", errors.New("withdraw failed")
	}
}

func (w *ChbtcTrade) cancelWithdraw(downloadId, safePwd string) (err error) {
	method := "cancelWithdraw"
	otherParams := "&currency=btc_cny&downloadId=" + downloadId + "&safePwd=" + safePwd

	req_para := w.createSign(method, otherParams)
	fmt.Println(req_para)
	body, err := util.HttpPost(Config["ch_api_url"]+method, req_para)
	if err != nil {
		return
	}

	js, err := NewJson([]byte(body))
	if err != nil {
		return
	}
	logger.Infoln(js)
	logger.Errorln(err)
	if jscode, ret := js.CheckGet("code"); ret {
		code := jscode.MustInt()

		if code == 1000 {
			return nil
		} else {
			logger.Errorln("code=", code)
			return errors.New("cancelWithdraw failed")
		}
	} else {
		return errors.New("cancelWithdraw failed")
	}
}
