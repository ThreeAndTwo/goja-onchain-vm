package goja_onchain_vm

import (
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	_ "github.com/dop251/goja_nodejs/require"
	"testing"
)

var vm = goja.New()

func TestEVMChain_GetProvider(t *testing.T) {
	const script = `
				var console = require('console')
				function provider(){
					return newProvider();
				}
				`

	tests := []struct {
		name string
		gvm  *VMGlobal
		want bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				Runtime: nil,
				ChainInfo: ChainInfo{
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
				Runtime:   vm,
				ChainInfo: ChainInfo{},
			},
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			want: true,
		},
		{
			name: "rpc node is inValidate",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e7be5",
					"",
				},
			},
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

			_, err = vm.RunString(script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}

			provider, ok := goja.AssertFunction(vm.Get("provider"))
			if !ok {
				t.Errorf("provider not a function")
				return
			}

			value, err := provider(goja.Undefined())
			if (err != nil) == true {
				t.Logf("new provider error %s", err)
				return
			}

			t.Log("value: ", value)
		})
	}
}

func TestEVMChain_GetBalance(t *testing.T) {
	const script = `
					var console = require('console')
					function balance(address) {
						return getBalance(address)
					}
				`

	tests := []struct {
		name    string
		gvm     *VMGlobal
		account string
		want    bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			account: "0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
			want:    true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				Runtime: nil,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			account: "0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
			want:    true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				Runtime:   vm,
				ChainInfo: ChainInfo{},
			},
			account: "0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
			want:    true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			account: "0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
			want:    true,
		},
		{
			name: "account is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			account: "",
			want:    true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			account: "0xa8C731e9259CE796B417A02ad0Cdcdd2057a0",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.init()
			if (err != nil) == tt.want {
				t.Logf("init gvm error %s", err)
				return
			}

			_, err = vm.RunString(script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}

			balance, ok := goja.AssertFunction(vm.Get("balance"))
			if !ok {
				t.Errorf("balance not a function")
				return
			}

			value, err := balance(goja.Undefined(), vm.ToValue(tt.account))
			if err != nil {
				t.Logf("new balance error %s", err)
				return
			}
			t.Log("value: ", value)
		})
	}
}

func TestEVMChain_GetTokenBalance(t *testing.T) {
	const script = `
				function tokenBalance(tokenType, contractAddr, account, tokenId) {
					var tokenId = string2BigInt(tokenId)
					return getTokenBalance(tokenType, contractAddr, account, tokenId)
				}
				`

	type params struct {
		tokenType    int
		contractAddr string
		account      string
		tokenId      string
	}

	tests := []struct {
		name  string
		gvm   *VMGlobal
		filed params
		want  bool
	}{
		{
			name: "normal",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				Runtime: nil,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				Runtime:   vm,
				ChainInfo: ChainInfo{},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "account is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"",
				"0",
			},
			want: true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "contract address is invalidation",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f93a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"",
			},
			want: true,
		},
		{
			name: "normal wax",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				1,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "token type unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				0,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
			want: true,
		},
		{
			name: "token type unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			filed: params{
				10,
				"0xdac17f958d2ee523a2206206994597c13d831ec7",
				"0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0",
				"0",
			},
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

			_, err = vm.RunString(script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}

			tokenBalance, ok := goja.AssertFunction(vm.Get("tokenBalance"))
			if !ok {
				t.Errorf("tokenBalance not a function")
				return
			}

			value, err := tokenBalance(goja.Undefined(), vm.ToValue(tt.filed.tokenType), vm.ToValue(tt.filed.contractAddr),
				vm.ToValue(tt.filed.account), vm.ToValue(tt.filed.tokenId))
			if err != nil {
				t.Logf("get token balance error %s", err)
				return
			}
			t.Log("value: ", value)
		})
	}
}

func TestEVMChain_Call(t *testing.T) {
	const script = `
					function callChain(contractAddr, data) {
						return call(contractAddr, data)
					}
				`

	type params struct {
		contract string
		data     string
	}

	tests := []struct {
		name   string
		gvm    *VMGlobal
		field params
		want   bool
	}{
		{
			name: "call normal",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "vm is null",
			gvm: &VMGlobal{
				Runtime: nil,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "chain not exist",
			gvm: &VMGlobal{
				Runtime:   vm,
				ChainInfo: ChainInfo{},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "chainId unSupport",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					-1,
					"",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "inValidate account",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4Dc610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "contract address is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "contract address is invalidation",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3396c610ea", "0x8da5cb5b"},
			want: true,
		},
		{
			name: "data is null",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", ""},
			want: true,
		},
		{
			name: "data is invalidation",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			field: params{"0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8dab5b"},
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

			_, err = vm.RunString(script)
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}

			call, ok := goja.AssertFunction(vm.Get("callChain"))
			if !ok {
				t.Errorf("callChain not a function")
				return
			}

			value, err := call(goja.Undefined(), vm.ToValue(tt.field.contract), vm.ToValue(tt.field.data))
			if err != nil {
				t.Logf("call error %s", err)
				return
			}
			t.Log("value: ", value)
		})
	}
}
