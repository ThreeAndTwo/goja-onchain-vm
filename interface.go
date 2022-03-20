package goja_onchain_vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type IChain interface {
	GetProvider() (*ethclient.Client, error)

	GetBalance(account string) (*big.Int, error)
	GetTokenBalance(tokenType TokenType, contractAddress, account string, tokenId *big.Int) (*big.Int, error)
	Call(to, data string) (string, error)
}

type IAccount interface {
	GetAddress() common.Address
}
