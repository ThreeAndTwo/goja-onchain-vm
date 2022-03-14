package goja_onchain_vm

import (
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	_ "github.com/dop251/goja_nodejs/require"
	"testing"
)

var vm = goja.New()

const (
	scriptCall = `
function run(){
    return contractCall("0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b");
}`
	scriptBalance = `
function run(){
	return balance("0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0")
}
`
	scriptTokenBalance = `
function run(){
	return tokenBalance(20, "0xdac17f958d2ee523a2206206994597c13d831ec7","0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0")
}
`
	scriptHttpGetRequest = `
function run(){
	return httpGetRequest()
}
`
	scriptHttpPostRequest = `
function run(){
	return httpPostRequest()
}
`
)

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
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: scriptCall,
			want:   true,
		},
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
			script: scriptBalance,
			want:   true,
		},
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
			script: scriptTokenBalance,
			want:   true,
		},
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
			script: scriptHttpGetRequest,
			want:   true,
		},
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
			script: scriptHttpPostRequest,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gvm.Init()
			if (err != nil) == tt.want {
				t.Logf("Init gvm error %s", err)
				return
			}

			_, err = vm.RunString(tt.script)
			runFunc, ok := goja.AssertFunction(vm.Get("run"))
			if !ok {
				t.Errorf("run not a function")
				return
			}
			value, err := runFunc(goja.Undefined())
			if err != nil {
				if _, ok := err.(*goja.InterruptedError); ok {
					t.Logf(`InterruptedError: %s`, err)
				} else {
					t.Errorf("unkonw error %s", err)
				}
				return
			}
			t.Log("value: ", value.String())
		})
	}
}
