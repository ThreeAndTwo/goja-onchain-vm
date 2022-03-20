package goja_onchain_vm

import (
	"github.com/deng00/ethutils"
	"github.com/ethereum/go-ethereum/common"
)

type Account struct {
	key   string
	index int
}

func NewAccount(key string, index int) *Account {
	return &Account{key: key, index: index}
}

func (a *Account) GetAddress() common.Address {
	var _account *ethutils.Account
	if ethutils.IsMnemonic(a.key) {
		_account = ethutils.GetAccountFromMnemonic(a.key, a.index)
	} else {
		_account = ethutils.GetAccountFromPStr(a.key)
	}
	if _account == nil {
		return [20]byte{}
	}
	return _account.Address
}
