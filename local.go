package goja_onchain_vm

import (
	"fmt"
	"github.com/deng00/ethutils"
	"github.com/dop251/goja"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Local struct {
	runtime     *goja.Runtime
	accountInfo AccountInfo
	publicKey   string
}

func NewLocal(runtime *goja.Runtime, account AccountInfo, publicKey string) (*Local, error) {
	return &Local{runtime: runtime, accountInfo: account, publicKey: publicKey}, nil
}

func (l *Local) GetAccountIndex() int {
	return l.accountInfo.Index
}

func (l *Local) SetAccountIndex(index int) {
	l.accountInfo.Index = index
}

func (l *Local) GetAddress() (string, error) {
	account := NewAccount(l.accountInfo.Key, l.accountInfo.Index).GetAccount()
	if account == nil {
		return "", fmt.Errorf("account invalidated")
	}
	return account.Address.String(), nil
}

func (l *Local) Signature(message []byte) (string, error) {
	account := NewAccount(l.accountInfo.Key, l.accountInfo.Index).GetAccount()
	if account == nil {
		return "", fmt.Errorf("account invalidated")
	}

	signature, err := ethutils.Sign(message, account.PrivateKey)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signature), err
}
