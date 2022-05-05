package goja_onchain_vm

import (
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/dop251/goja"
)

type VMGlobal struct {
	Runtime     *goja.Runtime
	ChainInfo   ChainInfo
	AccountInfo AccountInfo
	Url         string
	PublicKey   string
}

type ChainInfo struct {
	ChainId int64
	Rpc     string
	Wss     string
}

type RemoteInfo struct {
	Url    string
	Params string
}

type AccountInfo struct {
	AccountType AccountTy
	Key         string
	Index       int
	To          string
}

type ChainMCallGetter struct {
	MChainInfo *ChainInfo
	McClient   *multicall.Multicall
	McCounter  int
}

type TokenInfo struct {
	Address     string                      `json:"address"`
	Name        string                      `json:"name"`
	Symbol      string                      `json:"symbol"`
	Decimals    uint8                       `json:"decimals"`
	TotalSupply *multicall.BigIntJSONString `json:"total_supply"`
}

type VmFunc string

const (
	Balance               VmFunc = "balance"
	TokenBalance          VmFunc = "tokenBalance"
	CALL                  VmFunc = "contractCall"
	GetAddress            VmFunc = "getAddress"
	GetPreAddress         VmFunc = "getPreAddress"
	GetNextAddress        VmFunc = "getNextAddress"
	GetAddressByIndex     VmFunc = "getAddressByIndex"
	GetAddressListByIndex VmFunc = "getAddressListByIndex"
	GetCurrentIndex       VmFunc = "getCurrentIndex"
	PersonalSign          VmFunc = "personalSign"
	HttpGetRequest        VmFunc = "httpGetRequest"
	HttpPostRequest       VmFunc = "httpPostRequest"
	EncryptWithPubKey     VmFunc = "encryptWithPubKey"
)

type TokenType int64

const (
	NoType  TokenType = iota
	ERC20             = 20
	ERC721            = 721
	ERC1155           = 1155
)

type Call struct {
	To   string `json:"to"`
	Data string `json:"data"`
}

type RemoteData struct {
	Data string `json:"data"`
}
