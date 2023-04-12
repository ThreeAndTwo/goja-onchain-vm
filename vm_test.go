package goja_onchain_vm

import (
	"github.com/deng00/go-base/db/mongo"
	"github.com/dop251/goja"
	_ "github.com/dop251/goja_nodejs/console"
	_ "github.com/dop251/goja_nodejs/require"
	"reflect"

	mgo "go.mongodb.org/mongo-driver/mongo"
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

	jsHttpGetRequestData = `
function run() {
	return httpGetRequest("https://arweave.net/")['data']
}
`

	jsHttpGetRequestHeader = `
function run() {
	const data = httpGetRequest("https://arweave.net/")['header'];
	const dataObj = JSON.parse(data);
		return dataObj['Age'];
}
`

	jsHttpPostRequestData = `
function run() {
		var header = "{\"content-type\": \"application\/json\"}";
		var params = "{\"email\":\"eve.holt@reqres.in\",\"password\":\"pistol\"}";
		const data = httpPostRequest('https://reqres.in/api/register', params, header)['data'];
		const dataObj = JSON.parse(data);
		return dataObj['token'];
}
`

	jsHttpPostRequestHeader = `
function run() {
		var header = "{\"content-type\": \"application\/json\"}";
		var params = "{\"email\":\"eve.holt@reqres.in\",\"password\":\"pistol\"}";
		const _header = httpPostRequest('https://reqres.in/api/register', params, header)['header'];
		const headerObj = JSON.parse(_header);

		return headerObj['Access-Control-Allow-Origin'];
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
	return getAddressByIndex(2)
}
`
	jsGetAddressListByIndex = `
function run() {
	return getAddressListByIndex(1,3)
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
	jsRandomNumber = `
function run() {
	return randomNumber(20, 100);
}
`

	jsRandomNumberNav = `
function run() {
	return randomNumber(-10, -1);
}
`

	jsRandomBytes = `
function run() {
	return randomBytes();
}
`
	jsRandomBytes32 = `
function run() {
	return randomBytes(32);
}
`

	jsRandomBytes20 = `
function run() {
	return randomBytes(20);
}
`
	jsGetNonceOffset = `
function run(){
	const numArr = ["1", "2", "3", "4","5", "6", "7", "8", "9", "10"];
	return numArr[getNonceOffset()];
}
`
	jsGetNonce = `
function run(){
	const numArr = ["1", "2", "3", "4","5", "6", "7", "8", "9", "10"];
	return numArr[getNonce()];
}
`

	jsGetNonceOffset1 = `
function run(){
	const numArr = ["1", "2", "3", "4","5", "6", "7", "8", "9", "10"];
	return numArr[getNonceOffset()];
}
`

	jsGetNonceOffsetOutOfIndex = `
function run(){
	const numArr = ["1", "2", "3", "4","5", "6", "7", "8", "9", "10"];
	return numArr[getNonceOffset()];
}
`
	jsGetNonceOffsetOutOfIndex1 = `
function run(){
	const numArr = ["1", "2", "3", "4","5", "6", "7", "8", "9", "10"];
	return numArr[getNonceOffset()];
}
`

	jsDBGetViaCache = `
function run() {
	const value = db.get("key1");
	return value;
}
`

	jdDBGetViaDB = `
function run() {
	db.set("key14", "value14");
	return db.get("key14");
}
`

	jsDBGetErr = `
function run() {
	const value = db.get("");
	return value;
}
`

	jsDBSet = `
function run() {
	db.set("key4", "value44");
	return db.get("key4");
}
`
	jsDBSetErr = `
function run() {
	return db.set("key1", 123);
}
`
)

var kvMap map[string]string
var db *mgo.Database

func init() {
	kvMap = make(map[string]string)
	kvMap["key1"] = "value1"

	mongoConfig := &mongo.Config{
		URI:      "",
		Username: "",
		Password: "",
		Database: "",
	}
	_mongoClient, _ := mongo.New(mongoConfig)
	db = _mongoClient.Client.Database(mongoConfig.Database)
}

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
			script: jsHttpGetRequestData,
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
			script: jsHttpGetRequestHeader,
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
			script: jsHttpPostRequestData,
			want:   true,
		},
		{
			name: "http post req header",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					Key:   os.Getenv("MNEMONIC"),
					Index: 0,
				},
			},
			script: jsHttpPostRequestHeader,
			want:   true,
		},
		{
			name: "local for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       100,
				},
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 5,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "err: index not exists, remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "out of index, remote for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       25,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "index 1 for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 56,
				},
				AccountInfo: AccountInfo{
					AccountType: "Remote",
					Key:         "bnx",
					Index:       1,
					To:          "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
				},
				Url:       "https://aaa.bbb.io",
				PublicKey: os.Getenv("PUBKEY"),
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
					AccountType: LocalTy,
					Key:         "",
					Index:       1,
				},
			},
			script: jsGetAddress,
			want:   true,
		},
		{
			name: "local for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       1,
				},
			},
			script: jsGetPreAddress,
			want:   true,
		},
		{
			name: "remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       5,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getPreAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       30,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetPreAddress,
			want:   false,
		},
		{
			name: "index out for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
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
					AccountType: LocalTy,
					Key:         "",
					Index:       1,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
				},
			},
			script: jsGetNextAddress,
			want:   true,
		},
		{
			name: "remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       -1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "out of index: remote for getNextAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       24,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNextAddress,
			want:   false,
		},
		{
			name: "index out for getAddress func",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       -2,
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
					AccountType: LocalTy,
					Key:         "",
					Index:       0,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
				},
			},
			script: jsGetAddressByIndex,
			want:   true,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressByIndex,
			want:   false,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressByIndex,
			want:   false,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       24,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressByIndex,
			want:   false,
		},
		{
			name: "get address list by index success",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
				},
			},
			script: jsGetAddressListByIndex,
			want:   true,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressListByIndex,
			want:   false,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressListByIndex,
			want:   false,
		},
		{
			name: "remote for getAddressByIndex func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       24,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetAddressListByIndex,
			want:   false,
		},
		{
			name: "get address list by index failed",
			gvm: &VMGlobal{
				Runtime: vm,
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       1,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       2,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       3,
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
					AccountType: LocalTy,
					Key:         mnemonic,
					Index:       0,
				},
			},
			script: jsPersonalSign,
			want:   true,
		},
		{
			name: "remote for signature func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "remote for signature func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "remote for signature func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       24,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "remote for signature func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_PLAIN_PRIVATEKEY"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "remote for signature func",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_PLAIN_MNEMONIC"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "sign message via remote",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_SINGLE"),
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "sign message via remote",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_EV_MNEMONIC"),
					Index:       2,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "sign message via remote",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsPersonalSign,
			want:   false,
		},
		{
			name: "jsRandomNumber",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsRandomNumber,
			want:   false,
		},
		{
			name: "jsRandomNumberNav",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsRandomNumberNav,
			want:   false,
		},
		{
			name: "jsRandomBytes",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsRandomBytes,
			want:   false,
		},
		{
			name: "jsRandomBytes32",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsRandomBytes32,
			want:   false,
		},
		{
			name: "jsRandomBytes20",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     os.Getenv("RPC"),
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: RemoteTy,
					Key:         os.Getenv("TEST_MNEMONIC"),
					Index:       1,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsRandomBytes20,
			want:   false,
		},
		{
			name: "jsGetNonceOffset",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: 1,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNonceOffset,
			want:   false,
		},
		{
			name: "jsGetNonce",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: 0,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNonce,
			want:   false,
		},
		{
			name: "jsGetNonceOffset1",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -29,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNonceOffset1,
			want:   false,
		},
		{
			name: "jsGetNonceOffsetOutOfIndex",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -40,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNonceOffsetOutOfIndex,
			want:   false,
		},
		{
			name: "jsGetNonceOffsetOutOfIndex1",
			gvm: &VMGlobal{
				Runtime: vm,
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsGetNonceOffsetOutOfIndex1,
			want:   false,
		},
		{
			name: "error:db key duplicate",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBSet,
			want:   false,
		},
		{
			name: "error: db collection is null",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client: db,
					//Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBSet,
			want:   false,
		},
		{
			name: "normal: jsDBSet",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBSet,
			want:   false,
		},
		{
			name: "jsDBSetErr: client is null",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBSetErr,
			want:   false,
		},
		{
			name: "jsDBGetViaCache",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBGetViaCache,
			want:   false,
		},
		{
			name: "err:jsDBGetViaCache",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBGetErr,
			want:   false,
		},
		{
			name: "jdDBGetViaDB",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jdDBGetViaDB,
			want:   false,
		},
		{
			name: "jsDBSetErr",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: kvMap,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBSetErr,
			want:   false,
		},
		{
			name: "kvMap is null",
			gvm: &VMGlobal{
				Runtime: vm,
				DBConfig: GojaDB{
					Client:     db,
					Collection: "storage",
					Project: Project{
						Id:   "623d984e9e59cd594c3aac76",
						Name: "test",
					},
					KVMap: nil,
				},
				ChainInfo: ChainInfo{
					ChainId: 1,
					Rpc:     "https://mainnet.infura.io/v3/74312c6b77ac435fa2559c7e98277be5",
					Wss:     os.Getenv("WSS"),
				},
				AccountInfo: AccountInfo{
					AccountType: LocalTy,
					Key:         os.Getenv("TEST_PRIKEY"),
					NonceOffset: -25,
					Index:       0,
					To:          os.Getenv("TO"),
				},
				Url:       os.Getenv("URL"),
				PublicKey: os.Getenv("PUBLICKEY"),
			},
			script: jsDBGetViaCache,
			want:   false,
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
			_db := tt.gvm.DBConfig
			vm.Set("db", map[string]interface{}{
				"set": _db.Set,
				"get": _db.Get,
			})

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
			t.Logf("value: %s, value ty: %s", value.String(), reflect.TypeOf(value).String())
		})
	}
}
