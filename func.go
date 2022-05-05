package goja_onchain_vm

import (
	"fmt"
	"github.com/dop251/goja"
)

func NewFunc(runtime *goja.Runtime, chainInfo ChainInfo, accountInfo AccountInfo, remoteURL, publicKey string) (IFunc, error) {
	switch accountInfo.AccountType {
	case LocalTy:
		return NewLocal(runtime, accountInfo, publicKey)
	case RemoteTy:
		return NewRemote(runtime, chainInfo, accountInfo, remoteURL, publicKey)
	default:
		return nil, fmt.Errorf("unSupport account type")
	}
}
