/*
  trader API Engine
*/

package util

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
)

func http_req(req *http.Request) (body string, err error) {

	logger.Debugln("http_req req:", req)

	c := NewTimeoutClient()

	resp, err := c.Do(req)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer resp.Body.Close()

	logger.Debugln("http_req resp:", resp)
	if resp.StatusCode%200 != 0 {
		logger.Errorln("http_req resp:", resp)
		return
	}

	if resp.StatusCode != 200 {
		logger.Infoln("http_req resp:", resp.StatusCode)
	}

	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		body = DumpGZIP(resp.Body)
	default:
		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorln("read the http stream failed")
		} else {
			body = string(bodyByte)
		}
	}

	logger.Debugln("http_req body:", body)
	return
}

func HttpPost(api_url, req_para string) (body string, err error) {
	req_para_reader := strings.NewReader(req_para)

	req, err := http.NewRequest("POST", api_url, req_para_reader)
	if err != nil {
		logger.Errorln(err)
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return http_req(req)
}

func HttpGet(url string) (body string, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Errorln(err)
		return
	}

	return http_req(req)
}

func DumpGZIP(r io.Reader) string {
	var body string
	reader, _ := gzip.NewReader(r)
	for {
		buf := make([]byte, 1024)
		n, err := reader.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}
		body += string(buf)
	}
	return body
}
