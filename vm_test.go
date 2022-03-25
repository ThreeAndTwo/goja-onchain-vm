package goja_onchain_vm

import (
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	_ "github.com/dop251/goja_nodejs/require"
	"testing"
	"time"
)

var vm = goja.New()

const (
	mnemonic = `abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon cactus`

	jsCall = `
function run(){
    return contractCall("0xCc13Fc627EFfd6E35D2D2706Ea3C4D7396c610ea", "0x8da5cb5b");
}`
	jsBalance = `
function run(){
	return balance("0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0");
}
`
	jsTokenBalance = `
function run(){
	return tokenBalance(20, "0xdac17f958d2ee523a2206206994597c13d831ec7","0xa8C731e9259CE796B417A02aE7cd0Cdcdd2057a0");
}`

	jsGetAddress = `
function run(){
	return getAddress();
}
`
	jsGetPreAddress = `
function run(){
	return getPreAddress();
}
`

	jsGetNextAddress = `
function run(){
	return getNextAddress();
}
`

	jsHttpGetRequest = `
function run(){
	return httpGetRequest("");
}
`
	jsHttpPostRequest = `
function run(){
	return httpPostRequest("");
}
`
	jsEndlessLoop = `
function run(){
	var i = 0;
    for (;;) {
        i++;
    }
}
`
	jsGetAddressByIndex = `
function run() {
	return getAddressByIndex(100)
}
`
	jsGetAddressListByIndex = `
function run() {
	return getAddressListByIndex(1,100)
}
`
	jsGetAddressListByIndexError = `
function run() {
	return getAddressListByIndex(-1, 1)
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
			script: jsCall,
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
			script: jsBalance,
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
			script: jsTokenBalance,
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
			script: jsHttpGetRequest,
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
			script: jsHttpPostRequest,
			want:   true,
		},
		{
			name: "normal for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 100,
				},
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "index out for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: -1,
				},
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "mnemonic is null",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   "",
					Index: 1,
				},
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "normal for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 1,
				},
			},
			script: jsGetPreAddress,
			want:   true,
		},
		{
			name: "index out for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetPreAddress,
			want:   true,
		},
		{
			name: "mnemonic is null",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   "",
					Index: 1,
				},
			},
			script: jsGetPreAddress,
			want:   true,
		},
		{
			name: "normal for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetNextAddress,
			want:   true,
		},
		{
			name: "index out for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: -2,
				},
			},
			script: jsGetNextAddress,
			want:   true,
		},
		{
			name: "mnemonic is null",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   "",
					Index: 0,
				},
			},
			script: jsGetNextAddress,
			want:   true,
		},
		{
			name: "endlessLoop",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   "",
					Index: 0,
				},
			},
			script: jsEndlessLoop,
			want:   true,
		},
		{
			name: "get address by index",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetAddressByIndex,
			want:   true,
		},
		{
			name: "get address list by index success",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetAddressListByIndex,
			want:   true,
		},
		{
			name: "get address list by index failed",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetAddressListByIndexError,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now().UnixMilli()
			err := tt.gvm.Init()
			t.Logf("duration time %d", time.Now().UnixMilli()-startTime)
			if err != nil {
				t.Logf("Init gvm error %s", err)
				return
			}
			time.AfterFunc(5*time.Second, func() {
				vm.Interrupt("halt")
			})

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
