package goja_onchain_vm

import "math/big"

func String2BigInt(number string) *big.Int {
	n := new(big.Int)
	n, _ = n.SetString(number, 10)
	return n
}


