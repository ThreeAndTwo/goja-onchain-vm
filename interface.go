package goja_onchain_vm

import (
	"github.com/deng00/ethutils"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type IChain interface {
	GetProvider() (*ethclient.Client, error)

	GetBalance(account string) (*big.Int, error)
	GetNonce(address string, isPending bool) (uint64, error)
	GetTokenBalance(tokenType TokenType, contractAddress, account string, tokenId *big.Int) (*big.Int, error)
	Call(to, data string) (string, error)
}

type IAccount interface {
	GetAccount() *ethutils.Account
}

type IFunc interface {
	GetAccountIndex() int
	SetAccountIndex(index int)
	GetAddress() (string, error)
	Signature(message []byte) (string, error)
}
