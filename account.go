package goja_onchain_vm

import (
	"github.com/deng00/ethutils"
)

type Account struct {
	key   string
	index int
}

type AccountTy string

const (
	LocalTy  AccountTy = "Local"
	RemoteTy AccountTy = "Remote"
)

func NewAccount(key string, index int) *Account {
	return &Account{key: key, index: index}
}

func (a *Account) check() bool {
	return a.key == "" || a.index < 0
}

func (a *Account) GetAccount() *ethutils.Account {
	var _account *ethutils.Account
	if a.check() {
		return _account
	}

	if ethutils.IsMnemonic(a.key) {
		_account = ethutils.GetAccountFromMnemonic(a.key, a.index)
	} else {
		_account = ethutils.GetAccountFromPStr(a.key)
	}
	return _account
}
