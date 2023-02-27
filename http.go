package goja_onchain_vm

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"strings"
)

type ReqType string

const (
	POST ReqType = "post"
	GET  ReqType = "get"
)

type gojaReq struct {
	url     string
	header  req.Header
	param   req.Param
	reqType ReqType
	isJson  bool
}

func NewGojaReq(url string, header req.Header, param req.Param, reqType ReqType) *gojaReq {
	return &gojaReq{url: url, header: header, param: param, reqType: reqType}
}

func initHeader(header map[string]string) (req.Header, bool) {
	authHeader := req.Header{}
	hasJson := false

	for k, v := range header {
		authHeader[k] = v
		if hasJsonInHeader(k, v) {
			hasJson = true
		}
	}
	return authHeader, hasJson
}

func hasJsonInHeader(key, value string) bool {
	return strings.ToLower(key) == "content-type" && strings.Contains(strings.ToLower(value), "json")
}

func initParam(params map[string]string) req.Param {
	reqParams := req.Param{}
	for k, v := range params {
		reqParams[k] = v
	}
	return reqParams
}

func (gReq *gojaReq) check() bool {
	return gReq.url == ""
}

func (gReq *gojaReq) request() (map[string]string, error) {
	if gReq.check() {
		return nil, fmt.Errorf("url should be null")
	}

	switch gReq.reqType {
	case POST:
		return gReq.post()
	case GET:
		return gReq.get()
	default:
		return gReq.get()
	}
}

func (gReq *gojaReq) post() (map[string]string, error) {
	var reqResp = &req.Resp{}
	var err error
	if gReq.isJson {
		jsonParam, _ := json.Marshal(gReq.param)
		reqResp, err = req.Post(gReq.url, jsonParam, gReq.header)
	} else {
		reqResp, err = req.Post(gReq.url, gReq.param, gReq.header)
	}

	if reqResp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("%d status code != 200, error: %s", reqResp.Response().StatusCode, reqResp.String())
	}

	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	res["data"] = reqResp.String()
	bytesHeader, _ := json.Marshal(reqResp.Response().Header)
	res["header"] = string(bytesHeader)
	return res, err
}

func (gReq *gojaReq) get() (map[string]string, error) {
	resp, err := req.Get(gReq.url, gReq.header)
	if err != nil {
		return nil, err
	}

	if resp.Response().StatusCode != 200 {
		return nil, fmt.Errorf("%d status code != 200, error: %s", resp.Response().StatusCode, resp.String())
	}

	res := make(map[string]string)
	res["data"] = resp.String()
	bytesHeader, _ := json.Marshal(resp.Response().Header)
	res["header"] = string(bytesHeader)
	return res, nil
}
