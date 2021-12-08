package goja_onchain_vm

import (
	"github.com/dop251/goja"
	"testing"
)

var vm = goja.New()

func TestEVMChain_GetProvider(t *testing.T) {
	tests := []struct{
		name   string
		gvm    *VMGlobal
		script string
		want   bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `var provider = newProvider('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('provider: ', provider)`,
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				runtime: nil,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `var provider = newProvider('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('provider: ', provider)`,
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				runtime:   vm,
				chainInfo: ChainInfo{},
			},
			script: `var provider = newProvider('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('provider: ', provider)`,
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			script: `var provider = newProvider('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('provider: ', provider)`,
			want: true,
		},
		{
			name: "rpc node is inValidate",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e7be5",
					"",
				},
			},
			script: `var provider = newProvider('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('provider: ', provider)`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.init()
			if (err != nil) == tt.want {
				t.Logf("init gvm error %s", err)
				return
			}

			rs, err := vm.RunString(tt.script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}
			t.Log("rs: ", rs)
		})
	}
}

func TestEVMChain_GetBalance(t *testing.T) {
	tests := []struct {
		name   string
		gvm    *VMGlobal
		script string
		want   bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `var amount = getBalance('0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0')
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				runtime: nil,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				runtime:   vm,
				chainInfo: ChainInfo{},
			},
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			want: true,
		},
		{
			name: "account is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `var amount = getBalance('')
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `var amount = getBalance('0xa8C731e9259CE796B417A02aE7cd0Cdcda0')
				logTx('amount: ', amount)`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.init()
			if (err != nil) == tt.want {
				t.Logf("init gvm error %s", err)
				return
			}

			rs, err := vm.RunString(tt.script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}
			t.Log("rs: ", rs)
		})
	}
}

func TestEVMChain_GetTokenBalance(t *testing.T) {
	tests := []struct {
		name   string
		gvm    *VMGlobal
		script string
		want   bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				runtime: nil,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				runtime:   vm,
				chainInfo: ChainInfo{},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "account is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d831ec7','0xa8C731e9259CE796B417A02aE7cd0Cdcd057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'','0xa8C731e9259CE796B417A02aE7cd0Cdcd057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "contract address is invalidation",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0xdac17f958d2ee523a2206206994597c13d8','0xa8C731e9259CE796B417A02aE7cd0Cdcd057a0', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "normal wax",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(1,'0x39bb259f66e1c59d5abef88375979b4d20d98022','0x7be8076f4ea4a4ad08075c2508e481d6c946d12b', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "token type unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(0,'0x39bb259f66e1c59d5abef88375979b4d20d98022','0x7be8076f4ea4a4ad08075c2508e481d6c946d12b', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
		{
			name: "token type unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var tokenId = string2BigInt('0')
				var amount = getTokenBalance(10,'0x39bb259f66e1c59d5abef88375979b4d20d98022','0x7be8076f4ea4a4ad08075c2508e481d6c946d12b', tokenId)
				logTx('amount: ', amount)`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.init()
			if (err != nil) == tt.want {
				t.Logf("init gvm error %s", err)
				return
			}

			rs, err := vm.RunString(tt.script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}
			t.Log("rs: ", rs)
		})
	}
}

func TestEVMChain_Call(t *testing.T) {
	tests := []struct {
		name   string
		gvm    *VMGlobal
		script string
		want   bool
	}{
		{
			name: "call normal",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				runtime: nil,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea','0x8da5cb5b')
				logTx('res: ', res)`,
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				runtime:   vm,
				chainInfo: ChainInfo{},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4Dc610ea', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "contract address is invalidation",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3396c610ea', '0x8da5cb5b')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "data is null",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea', '')
				logTx('res ', res)`,
			want: true,
		},
		{
			name: "data is invalidation",
			gvm: &VMGlobal{
				runtime: vm,
				chainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: `
				var res = call('0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea', '0x8dab5b')
				logTx('res ', res)`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.init()
			if (err != nil) == tt.want {
				t.Logf("init gvm error %s", err)
				return
			}

			rs, err := vm.RunString(tt.script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}
			t.Log("rs: ", rs)
		})
	}
}
