package goja_onchain_vm

import (
	"encoding/json"
	"fmt"
	goEthutils "github.com/CoinSummer/go-ethutils"
	"github.com/dop251/goja"
	"strings"
)

type remotePath string

const (
	RemoteAddress     remotePath = "/v1/address"
	RemoteSignMessage remotePath = "/v1/sign/message"
)

type Remote struct {
	runtime     *goja.Runtime
	chainInfo   ChainInfo
	accountInfo AccountInfo
	url         string
	publicKey   string
}

func NewRemote(runtime *goja.Runtime, chainInfo ChainInfo, accountInfo AccountInfo, url string, publicKey string) (*Remote, error) {
	return &Remote{runtime: runtime, chainInfo: chainInfo, accountInfo: accountInfo, url: url, publicKey: publicKey}, nil
}

func (r *Remote) GetAccountIndex() int {
	return r.accountInfo.Index
}

func (r *Remote) SetAccountIndex(index int) {
	r.accountInfo.Index = index
}

func (r *Remote) GetAddress() (string, error) {
	header := `{"content-type": "application/json"}`
	params := fmt.Sprintf(`{"chain_id": %d, "index": %d}`,
		r.chainInfo.ChainId, r.accountInfo.Index)
	encryptParam, err := r.encryptWithPubKey(params)
	if err != nil {
		return "", err
	}

	encryptMsg := `{"encryptMsg":"` + encryptParam + `"}`
	data := &RemoteData{}
	res, err := r.post(r.url+string(RemoteAddress), encryptMsg, header)
	if err != nil {
		return "", err
	}
	resArr := strings.Split(res, "||")

	if len(resArr) == 0 {
		return "", fmt.Errorf(`get remote address error, index:` + fmt.Sprintf("%d", r.accountInfo.Index))
	}

	err = json.Unmarshal([]byte(resArr[0]), data)
	if err != nil {
		return "", fmt.Errorf(`get remote address error, index:` + fmt.Sprintf("%d", r.accountInfo.Index))
	}

	if data.Data == "" {
		return "", fmt.Errorf("address not exist for remote")
	}

	return data.Data, nil
}

func (r *Remote) Signature(message []byte) (string, error) {
	header := `{"content-type": "application/json"}`
	params := fmt.Sprintf(`{"chain_id": %d, "index": %d, "to": "%s", "message": "%s"}`,
		r.chainInfo.ChainId, r.accountInfo.Index, r.accountInfo.To, message)

	encryptParam, err := r.encryptWithPubKey(params)
	if err != nil {
		return "", err
	}

	encryptMsg := `{"encryptMsg":"` + encryptParam + `"}`
	data := &RemoteData{}
	res, err := r.post(r.url+string(RemoteSignMessage), encryptMsg, header)
	if err != nil {
		return "", err
	}
	resArr := strings.Split(res, "||")

	if len(resArr) == 0 {
		return "", fmt.Errorf(`get remote signature error, index:` + fmt.Sprintf("%d", r.accountInfo.Index))
	}

	err = json.Unmarshal([]byte(resArr[0]), data)
	if err != nil {
		return "", fmt.Errorf(`get remote signature error, index:` + fmt.Sprintf("%d", r.accountInfo.Index))
	}

	if data.Data == "" {
		return "", fmt.Errorf("can not signature the message for remote")
	}

	return data.Data, nil
}

func (r *Remote) post(url, params, header string) (string, error) {
	reqHeader, reqParam, isJson, err := getReqParam(params, header)
	if err != nil {
		return "", err
	}

	_req := NewGojaReq(url, reqHeader, reqParam, POST)
	_req.isJson = isJson
	resp, err := _req.request()
	if err != nil {
		return "", err
	}
	return resp["data"], nil
}

func (r *Remote) get(url, params, header string) (string, error) {
	reqHeader, reqParam, _, err := getReqParam(params, header)
	if err != nil {
		return "", err
	}

	_req := NewGojaReq(url, reqHeader, reqParam, GET)
	data, err := _req.request()
	if err != nil {
		return "", err
	}
	return data["data"], err
}

func (r *Remote) encryptWithPubKey(message string) (string, error) {
	if r.publicKey == "" || message == "" {
		return "", fmt.Errorf("params invalidate for encryptWithPubKey")
	}

	key, err := goEthutils.EncryptByPubKey("0x"+string(r.publicKey), message)
	if err != nil {
		return "", err
	}

	ed, err := key.Stringify()
	if err != nil {
		return "", err
	}
	return ed, nil
}
