package goja_onchain_vm

import (
	"github.com/alethio/web3-multicall-go/multicall"
	"github.com/dop251/goja"
)

type VMGlobal struct {
	Runtime     *goja.Runtime
	ChainInfo   ChainInfo
	AccountInfo AccountInfo
}

type ChainInfo struct {
	ChainId int
	Rpc     string
	Wss     string
}

type AccountInfo struct {
	Key   string
	Index int
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

type ChainID int

const (
	ETH         ChainID = 1
	OPTIMISTIC  ChainID = 10
	CRONOS      ChainID = 25
	BSC         ChainID = 56
	OKEX        ChainID = 66
	HECO        ChainID = 128
	POLYGON     ChainID = 137
	FTM         ChainID = 250
	ARBITRUMONE ChainID = 42161
	AVALANCHE   ChainID = 43114
)

type VmFunc string

const (
	Balance         VmFunc = "balance"
	TokenBalance    VmFunc = "tokenBalance"
	CALL            VmFunc = "contractCall"
	GetAddress      VmFunc = "getAddress"
	GetPreAddress   VmFunc = "getPreAddress"
	GetNextAddress  VmFunc = "getNextAddress"
	HttpGetRequest  VmFunc = "httpGetRequest"
	HttpPostRequest VmFunc = "httpPostRequest"
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
