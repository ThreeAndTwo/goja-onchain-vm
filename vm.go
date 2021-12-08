package goja_onchain_vm

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"math/big"
)

func (gvm *VMGlobal) init() error {
	registry := require.NewRegistry()
	if !gvm.check() {
		return fmt.Errorf("gvm config error, please check your config")
	}

	vm := gvm.runtime
	registry.Enable(vm)
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	err := vm.Set(string(NEWPROVIDER), gvm.NewProvider)
	if err != nil {
		return err
	}

	err = vm.Set(string(GETBALANCE), gvm.GetBalance)
	if err != nil {
		return err
	}

	err = vm.Set(string(GETTOKENBALANCE), gvm.GetTokenBalance)
	if err != nil {
		return err
	}

	err = vm.Set(string(CALL), gvm.Call)
	if err != nil {
		return err
	}

	return vm.Set(string(STRING2BIGINT), gvm.String2BigInt)
}

func (gvm *VMGlobal) check() bool {
	return nil != gvm && nil != gvm.runtime
}

func (gvm *VMGlobal) NewProvider() goja.Value {
	chain, err := ChainGetter(&gvm.chainInfo)
	if err != nil {
		gvm.runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.runtime.ToValue(`should be catch exception`)
	}

	provider, err := chain.GetProvider()
	if err != nil {
		gvm.runtime.Interrupt(`get chain provider error:` + err.Error())
		return gvm.runtime.ToValue(`should be catch exception`)
	}

	return gvm.runtime.ToValue(provider)
}

func (gvm *VMGlobal) GetBalance(account string) goja.Value {
	chain, err := ChainGetter(&gvm.chainInfo)
	if err != nil {
		gvm.runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}

	balance, err := chain.GetBalance(account)
	if err != nil {
		gvm.runtime.Interrupt(`get chain balance error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}

	return gvm.runtime.ToValue(balance)
}

func (gvm *VMGlobal) GetTokenBalance(_tokenType TokenType, _contractAddress, _account string, _tokenID *big.Int) goja.Value {
	chain, err := ChainGetter(&gvm.chainInfo)
	if err != nil {
		gvm.runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}

	balance, err := chain.GetTokenBalance(_tokenType, _contractAddress, _account, _tokenID)
	if err != nil {
		gvm.runtime.Interrupt(`get chain token balance error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}
	return gvm.runtime.ToValue(balance)
}

func (gvm *VMGlobal) Call(to, data string) goja.Value {
	chain, err := ChainGetter(&gvm.chainInfo)
	if err != nil {
		gvm.runtime.Interrupt(`new chain error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}

	callData, err := chain.Call(to, data)
	if err != nil {
		gvm.runtime.Interrupt(`call chain error:` + err.Error())
		return gvm.runtime.ToValue(`exception`)
	}
	return  gvm.runtime.ToValue(callData)
}

func (gvm *VMGlobal) String2BigInt(number string) goja.Value {
	return gvm.runtime.ToValue(String2BigInt(number))
}
