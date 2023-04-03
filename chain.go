package goja_onchain_vm

import (
	"context"
	"fmt"
	"github.com/ThreeAndTwo/goja-onchain-vm/contract"
	"github.com/deng00/ethutils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type EVMChain struct {
	chainInfo    *ChainInfo
	chainClient  *ethclient.Client
	mcCallGetter *ChainMCallGetter
}

func ChainGetter(_chainInfo *ChainInfo) (IChain, error) {
	return NewEVMChain(_chainInfo)
}

func NewEVMChain(_chainInfo *ChainInfo) (*EVMChain, error) {
	client, err := ethclient.Dial(_chainInfo.Rpc)
	if err != nil {
		return nil, err
	}
	return &EVMChain{chainInfo: _chainInfo, chainClient: client}, nil
}

func (E *EVMChain) GetProvider() (*ethclient.Client, error) {
	return ethclient.Dial(E.chainInfo.Rpc)
}

func (E *EVMChain) GetBalance(account string) (*big.Int, error) {
	if !isValidateAddress(account) {
		return nil, fmt.Errorf("account is invalidation")
	}
	return E.chainClient.PendingBalanceAt(context.Background(), common.HexToAddress(account))
}

func (E *EVMChain) GetTokenBalance(tokenType TokenType, contractAddress, account string, tokenId *big.Int) (*big.Int, error) {
	if !isValidateAddress(contractAddress) || !isValidateAddress(account) {
		return nil, fmt.Errorf("address is invalidation")
	}

	tokenAddr := common.HexToAddress(contractAddress)
	address := common.HexToAddress(account)

	var balance *big.Int
	var err error
	switch tokenType {
	case ERC20:
		balance, err = E.erc20TokenBalance(tokenAddr, address)
	case ERC721:
		balance, err = E.erc721TokenBalance(tokenAddr, address)
	case ERC1155:
		balance, err = E.erc1155TokenBalance(tokenAddr, address, tokenId)
	default:
		return nil, fmt.Errorf("not support: %d", tokenType)
	}
	return balance, err
}

func (E *EVMChain) GetNonce(address string, isPending bool) (uint64, error) {
	_address := common.HexToAddress(address)
	if !isPending {
		return E.chainClient.NonceAt(context.TODO(), _address, nil)
	}
	return E.chainClient.PendingNonceAt(context.TODO(), _address)
}

func (E *EVMChain) Call(to, data string) (string, error) {
	if !isValidateAddress(to) {
		return "", fmt.Errorf("address is invalidation")
	}

	receipt := common.HexToAddress(to)
	callReq := ethereum.CallMsg{
		To:   &receipt,
		Data: common.FromHex(data),
	}

	res, err := E.chainClient.CallContract(context.TODO(), callReq, nil)
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(res), nil
}

func (E *EVMChain) erc20TokenBalance(contractAddress, account common.Address) (*big.Int, error) {
	_contract, err := contract.NewErc20(contractAddress, E.chainClient)
	if err != nil {
		return nil, fmt.Errorf("NewErc20 error: %s", err)
	}
	return _contract.BalanceOf(ethutils.MakeCallOpts(account), account)
}

func (E *EVMChain) erc721TokenBalance(contractAddress, account common.Address) (*big.Int, error) {
	_contract, err := contract.NewErc721(contractAddress, E.chainClient)
	if err != nil {
		return nil, fmt.Errorf("NewErc721 error: %s", err)
	}
	return _contract.BalanceOf(ethutils.MakeCallOpts(account), account)
}

func (E *EVMChain) erc1155TokenBalance(contractAddress, account common.Address, tokenId *big.Int) (*big.Int, error) {
	_contract, err := contract.NewERC1155(contractAddress, E.chainClient)
	if err != nil {
		return nil, fmt.Errorf("NewErc1155 error: %s", err)
	}
	return _contract.BalanceOf(ethutils.MakeCallOpts(account), account, tokenId)
}

func isValidateAddress(account string) bool {
	return 42 == len(account)
}
