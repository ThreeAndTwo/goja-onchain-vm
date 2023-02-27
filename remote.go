package goja_onchain_vm

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/dop251/goja"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
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
	params := fmt.Sprintf(`{"chain_id": %d, "account": "%s", "index": %d, "to": "%s"}`,
		r.chainInfo.ChainId, r.accountInfo.Key, r.accountInfo.Index, r.accountInfo.To)
	encryptParam, err := r.encryptWithPubKey(params)
	if err != nil {
		return "", err
	}

	encryptMsg := `{"encryptMsg":"` + encryptParam + `"}`
	data := &RemoteData{}
	res, err := r.post(r.url+"/v1/address", encryptMsg, header)
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
	params := fmt.Sprintf(`{"chain_id": %d, "account": "%s", "index": %d, "to": "%s", "message": "%s"}`,
		r.chainInfo.ChainId, r.accountInfo.Key, r.accountInfo.Index, r.accountInfo.To, message)

	encryptParam, err := r.encryptWithPubKey(params)
	if err != nil {
		return "", err
	}

	encryptMsg := `{"encryptMsg":"` + encryptParam + `"}`
	data := &RemoteData{}
	res, err := r.post(r.url+"/signature", encryptMsg, header)
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

	signerKey, err := hexutil.Decode("0x" + string(r.publicKey))
	if err != nil {
		return "", fmt.Errorf(`decode public key error:` + err.Error())
	}

	pubKey, err := btcec.ParsePubKey(signerKey, btcec.S256())
	if err != nil {
		return "", fmt.Errorf(`parse public key error:` + err.Error())
	}

	encryptData, err := btcec.Encrypt(pubKey, []byte(message))
	if err != nil {
		return "", fmt.Errorf(`encrypt data error:` + err.Error())
	}
	return hexutil.Encode(encryptData), nil
}
