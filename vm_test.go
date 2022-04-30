package goja_onchain_vm

import (
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	_ "github.com/dop251/goja_nodejs/require"
	"os"
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
	var console = require('console')
	var data = httpGetRequest("https://demo.example.com/api/v1/connect?walletAddress=0x00")
	var dataArr = data.split('||')
	
	var reqData = JSON.parse(dataArr[0])
	console.log("cookie:", dataArr[1])
	return reqData.data
}
`
	jsHttpPostRequest = `
function run(){
	return httpPostRequest("");
}
`

	jsHttpPostRequestContainsHeader = `
function run(){
	var console = require('console');
	var header = "{\"content-type\": \"application\/json\"}";
	var params = "{\"signature\":\"0x00\",\"walletAddress\":\"0x00\"}";
	
	var reqData = httpPostRequest('https://demo.example.com/api/v1/sign-in', params, header);
	var dataArr = reqData.split('||');

	var data = JSON.parse(dataArr[0]);
	console.log("data:", dataArr[0]);
	console.log("cookie:", dataArr[1]);
	return data.data;
}
`
	jsHttpMix = `
		var console = require('console');
		var data = httpGetRequest("https://demo.example.com/api/v1/connect?walletAddress=0x00");
		var dataArr = data.split('||');

		var message = JSON.parse(dataArr[0]);
		var signMessage = personalSign(message.data);

		var header = "{\"content-type\": \"application\/json\"}";
		var params = "{\"signature\":\""+ signMessage +"\",\"walletAddress\":\"0x00\"}";

		var signInData = httpPostRequest('https://demo.example.com/api/v1/sign-in', params, header);
		var signInArr = signInData.split('||');
		var signCookie = JSON.parse(signInArr[1]);
		
		var myRe = /authorization=(\S*); Path=/;
		var myArray = myRe.exec(signCookie['Set-Cookie']);
		var cookieArr = myArray[0].split(';');

		var signCookie = "NEXT_LOCALE=en; " + cookieArr[0] + "; tr=" + Date.parse(new Date());
		var signHeader = "{\"cookie\": \""+ signCookie +"\"}";	
		var signData = httpGetRequest("https://demo.example.com/api/v1/firework-mint-data?walletAddress=0x00&amount=3", signHeader);
		
		var signDataArr = signData.split('||');
		var inputData = JSON.parse(signDataArr[0]);
		console.log("inputData:", inputData.data);
		return inputData.data
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
	jsGetCurrentIndex = `
function run() {
	var tokenList = [100,200,300,400]
	return tokenList[getCurrentIndex()]
}
`
	jsGetCurrentIndexViaString = `
function run() {
	var tokenList = ["100","200","300","400"]
	return tokenList[getCurrentIndex()]
}
`
	jsGetCurrentIndexViaMix = `
function run() {
	var tokenList = [1, "200", 30, "4"]
	return tokenList[getCurrentIndex()]
}
`
	jsGetCurrentIndexViaErr = `
function run() {
	var tokenList = [100,"200","30,"40"]
	return tokenList[getCurrentIndex()]
}
`
	jsPersonalSign = `
function run() {
	var message = "Welcome to TRLab!\nWallet address:\n0x00\nNonce\n653888"
	var signMessage = personalSign(message)
	return signMessage
}
`
	jsEncryptPubKey = `
function run() {
	var message = '{"chain_id":1,"account":"test","to":"0x287e21B9201E98ef3E2E0e8759Ee36Ca8257a6d2","message":"aaaaaa"}'
	return encryptWithPubKey(message)
}
`
	jsEncryptPubKeyGetAccount = `
function run() {
	var message = '{"chain_id":1,"account":"test","to":"0x287e21B9201E98ef3E2E0e8759Ee36Ca8257a6d2"}'
	return encryptWithPubKey(message)
}
`
)

func TestEVMChain(t *testing.T) {
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
			name: "normal http",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					1,
					"https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					"",
				},
			},
			script: jsHttpPostRequestContainsHeader,
			want:   true,
		},
		{
			name: "normal js mix",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   os.Getenv("MNEMONIC"),
					Index: 0,
				},
			},
			script: jsHttpMix,
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
					AccountType: "remote",
					Key:         "aaa",
					Index:       0,
					To:          "0x00",
				},
				Url:       "http://127.0.0.1/ping",
				PublicKey: os.Getenv("PUBKEY"),
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
		{
			name: "get current index",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 0,
				},
			},
			script: jsGetCurrentIndex,
			want:   true,
		},
		{
			name: "get current index via string",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 1,
				},
			},
			script: jsGetCurrentIndexViaString,
			want:   true,
		},
		{
			name: "get current index via mix",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 2,
				},
			},
			script: jsGetCurrentIndexViaMix,
			want:   true,
		},
		{
			name: "get current index via error",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   mnemonic,
					Index: 3,
				},
			},
			script: jsGetCurrentIndexViaErr,
			want:   true,
		},
		{
			name: "sign message",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     "",
				},
				AccountInfo: AccountInfo{
					Key:   os.Getenv("MNEMONIC"),
					Index: 0,
				},
			},
			script: jsPersonalSign,
			want:   true,
		},
		{
			name: "encrypt with pubKey",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "",
					Wss:     "",
				},
				AccountInfo: AccountInfo{},
				PublicKey:   os.Getenv("PUBKEY"),
			},
			script: jsEncryptPubKey,
			want:   true,
		},
		{
			name: "encrypt with pubKey2",
			gvm: &VMGlobal{
				Runtime:   vm,
				PublicKey: os.Getenv("PUBKEY"),
			},
			script: jsEncryptPubKeyGetAccount,
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
			time.AfterFunc(100*time.Second, func() {
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
				if _, ok = err.(*goja.InterruptedError); ok {
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
