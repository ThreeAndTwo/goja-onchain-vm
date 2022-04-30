package goja_onchain_vm

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/deng00/ethutils"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/imroc/req"
	"math/big"
	"strings"
)

func (gvm *VMGlobal) Init() error {
	registry := require.NewRegistry()
	if !gvm.check() {
		return fmt.Errorf("gvm config error, please check your config")
	}

	vm := gvm.Runtime
	registry.Enable(vm)
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	err := vm.Set(string(Balance), gvm.GetBalance)
	if err != nil {
		return err
	}

	err = vm.Set(string(TokenBalance), gvm.GetTokenBalance)
	if err != nil {
		return err
	}

	err = vm.Set(string(CALL), gvm.Call)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetAddress), gvm.GetAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetPreAddress), gvm.GetPreAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetNextAddress), gvm.GetNextAddress)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetAddressByIndex), gvm.GetAddressByIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetAddressListByIndex), gvm.GetAddressListByIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(GetCurrentIndex), gvm.GetCurrentIndex)
	if err != nil {
		return err
	}

	err = vm.Set(string(PersonalSign), gvm.GetPersonalSign)
	if err != nil {
		return err
	}

	err = vm.Set(string(EncryptWithPubKey), gvm.EncryptWithPubKey)
	if err != nil {
		return err
	}

	// http get
	err = vm.Set(string(HttpGetRequest), gvm.HttpGet)
	if err != nil {
		return err
	}

	// http post
	return vm.Set(string(HttpPostRequest), gvm.HttpPost)
}

func (gvm *VMGlobal) check() bool {
	return nil != gvm && nil != gvm.Runtime
}

func (gvm *VMGlobal) GetBalance(account string) goja.Value {
	chain, err := ChainGetter(&gvm.ChainInfo)
	if err != nil {
		gvm.Runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	balance, err := chain.GetBalance(account)
	if err != nil {
		gvm.Runtime.Interrupt(`get chain balance error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	return gvm.Runtime.ToValue(balance)
}

func (gvm *VMGlobal) GetTokenBalance(_tokenType TokenType, _contractAddress, _account string, _tokenID *big.Int) goja.Value {
	chain, err := ChainGetter(&gvm.ChainInfo)
	if err != nil {
		gvm.Runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	balance, err := chain.GetTokenBalance(_tokenType, _contractAddress, _account, _tokenID)
	if err != nil {
		gvm.Runtime.Interrupt(`get chain token balance error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(balance)
}

func (gvm *VMGlobal) Call(to, data string) goja.Value {
	chain, err := ChainGetter(&gvm.ChainInfo)
	if err != nil {
		gvm.Runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	callData, err := chain.Call(to, data)
	if err != nil {
		gvm.Runtime.Interrupt(`call chain error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(callData)
}

func (gvm *VMGlobal) HttpGet(url, params, header string) goja.Value {
	reqHeader, reqParam, _, err := getReqParam(params, header)
	if err != nil {
		gvm.Runtime.Interrupt(err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	_req := NewGojaReq(url, reqHeader, reqParam, GET)
	data, err := _req.request()
	if err != nil {
		gvm.Runtime.Interrupt(`http request error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(data)
}

func (gvm *VMGlobal) HttpPost(url, params, header string) goja.Value {
	reqHeader, reqParam, isJson, err := getReqParam(params, header)
	if err != nil {
		gvm.Runtime.Interrupt(err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	_req := NewGojaReq(url, reqHeader, reqParam, POST)
	_req.isJson = isJson
	resp, err := _req.request()
	if err != nil {
		gvm.Runtime.Interrupt(`http request error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(resp)
}

func getReqParam(params, header string) (req.Header, req.Param, bool, error) {
	headerMap := make(map[string]string)
	paramsMap := make(map[string]string)

	if header != "" {
		err := json.Unmarshal([]byte(header), &headerMap)
		if err != nil {
			return nil, nil, false, fmt.Errorf("http params invalidate for header: %s", header)
		}
	}

	if params != "" {
		err := json.Unmarshal([]byte(params), &paramsMap)
		if err != nil {
			return nil, nil, false, fmt.Errorf("http params invalidate for params: %s", params)
		}
	}

	reqHeader, isJson := initHeader(headerMap)
	reqParam := initParam(paramsMap)
	return reqHeader, reqParam, isJson, nil
}

func (gvm *VMGlobal) GetAddress() goja.Value {
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetPreAddress() goja.Value {
	gvm.AccountInfo.Index--
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetNextAddress() goja.Value {
	gvm.AccountInfo.Index++
	return gvm.getAddress()
}

func (gvm *VMGlobal) GetAddressByIndex(index int) goja.Value {
	gvm.AccountInfo.Index = index
	return gvm.getAddress()
}

func (gvm *VMGlobal) getAddress() goja.Value {
	if gvm.checkAddress() {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	if gvm.AccountInfo.AccountType != "" {
		// remote
		header := `{"content-type": "application/json"}`
		params := fmt.Sprintf(`{"chain_id": %d, "account": "%s", "index": %d, "to": "%s"}`,
			gvm.ChainInfo.ChainId, gvm.AccountInfo.Key, gvm.AccountInfo.Index, gvm.AccountInfo.To)
		var encryptParam = gvm.EncryptWithPubKey(params)
		encryptMsg := `{"encryptMsg":"` + encryptParam.String() + `"}`
		return gvm.HttpPost(gvm.Url+"/address", encryptMsg, header)
	}
	account := NewAccount(gvm.AccountInfo.Key, gvm.AccountInfo.Index).GetAccount()
	if account == nil {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	address := account.Address.String()
	return gvm.Runtime.ToValue(address)
}

func (gvm *VMGlobal) GetAddressListByIndex(start, end int) goja.Value {
	gvm.AccountInfo.Index = start
	if gvm.checkAddress() {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	var arrAddr []string
	for k := start; k < end; k++ {
		gvm.AccountInfo.Index = k

		account := NewAccount(gvm.AccountInfo.Key, gvm.AccountInfo.Index).GetAccount()
		if account == nil {
			gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
			return gvm.Runtime.ToValue(`exception`)
		}

		addr := account.Address.String()
		arrAddr = append(arrAddr, addr)
	}
	addresses := strings.Join(arrAddr, ",")
	return gvm.Runtime.ToValue(addresses)
}

func (gvm *VMGlobal) checkAddress() bool {
	return gvm.AccountInfo.Key == "" || gvm.AccountInfo.Index < 0
}

func (gvm *VMGlobal) GetCurrentIndex() goja.Value {
	return gvm.Runtime.ToValue(gvm.AccountInfo.Index)
}

func (gvm *VMGlobal) GetPersonalSign(message string) goja.Value {
	if gvm.checkAddress() {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	account := NewAccount(gvm.AccountInfo.Key, gvm.AccountInfo.Index).GetAccount()
	if account == nil {
		gvm.Runtime.Interrupt(`params invalidate for address, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index))
		return gvm.Runtime.ToValue(`exception`)
	}

	signature, err := ethutils.Sign([]byte(message), account.PrivateKey)
	if err != nil {
		gvm.Runtime.Interrupt(`get sign error, index:` + fmt.Sprintf("%d", gvm.AccountInfo.Index) +
			`error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(hexutil.Encode(signature))
}

func (gvm *VMGlobal) EncryptWithPubKey(message string) goja.Value {
	if gvm.PublicKey == "" || message == "" {
		gvm.Runtime.Interrupt(`params invalidate for encryptWithPubKey`)
		return gvm.Runtime.ToValue(`exception`)
	}

	signerKey, err := hexutil.Decode("0x" + string(gvm.PublicKey))
	if err != nil {
		gvm.Runtime.Interrupt(`decode public key error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	pubKey, err := btcec.ParsePubKey(signerKey, btcec.S256())
	if err != nil {
		gvm.Runtime.Interrupt(`parse public key error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}

	encryptData, err := btcec.Encrypt(pubKey, []byte(message))
	if err != nil {
		gvm.Runtime.Interrupt(`encrypt data error:` + err.Error())
		return gvm.Runtime.ToValue(`exception`)
	}
	return gvm.Runtime.ToValue(hexutil.Encode(encryptData))
}
