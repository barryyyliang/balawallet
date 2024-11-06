package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	//"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

// //页面切换
// func doChangePage(page string, param string) {
// 	fmt.Println(resdir + page)
// 	content, err := os.ReadFile(resdir + page)
// 	if err != nil {
// 		fmt.Println("can not read html")
// 	}
// 	win.SetHtml(string(content))
// }

// 新创建一个账号
func doNewAccount(pass string, walletName string) string {
	re := getResult(Msg_success, "")
	//scryptN := keystore.StandardScryptN
	//scryptP := keystore.StandardScryptP
	//account, err := keystore.StoreKey(accountdir, pass, scryptN, scryptP)
	account, err := ks.NewAccount(pass)
	if err != nil {
		fmt.Println("failed to generate new account")
		re = getResult(Msg_createuser_fail, err.Error())
		return string(re)
	}

	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err == nil {
		defer db.Close()
		if isExist(db, walletName) == false {
			err = db.Put([]byte(account.Address.String()), []byte(walletName), nil)
			if err != nil {
				fmt.Println("put name err")
				re = getResult(Msg_writename_fail, err.Error())
			}
		} else {
			fmt.Println("name has existed")
			re = getResult(Msg_createuser_nameexists, "")
		}
	} else {
		re = getResult(Msg_writename_fail, err.Error())
	}

	return string(re)

}

// 日志接口
func doWebLog(para string) {
	fmt.Println(para)
}

// 获取账号列表
func doGetAccounts() string {
	type account struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	type accounts struct {
		Accounts []account `json:"accounts"`
	}

	//账户信息
	addressCache := make(map[string]string)
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err == nil {
		defer db.Close()
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			addressCache[string(it.Key())] = string(it.Value())
		}
	}

	acc := ks.Accounts()
	var acclist accounts
	for i := 0; i < len(acc); i++ {
		add := acc[i].Address.String()
		acclist.Accounts = append(acclist.Accounts, account{"eth", add, addressCache[add]})
	}

	result, err2 := json.Marshal(acclist)
	if err2 == nil {
		return string(result)
	} else {
		fmt.Println(err2.Error())
	}
	return ""
}

func doChangeAccountName(address string, name string) string {
	re := getResult(Msg_success, "")
	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err == nil {
		defer db.Close()
		if isExist(db, name) == false {
			err = db.Put([]byte(address), []byte(name), nil)
			if err != nil {
				re = getResult(Msg_writedb_fail, err.Error())
			}
		} else {
			re = getResult(Msg_changename_nameexists, "")
		}
	} else {
		re = getResult(Msg_writedb_fail, err.Error())
	}

	return string(re)
}

func doGetAccountName(address string) string {
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err == nil {
		defer db.Close()
		name, err := db.Get([]byte(address), nil)
		if err == nil {
			return string(name)
		} else {
			fmt.Println("get name err")
			return ""
		}
	} else {
		fmt.Println("open db error:" + err.Error())
		return ""
	}
}

// 通过私钥导入账号
func doImportAccount(pa string, name string, kk string) string {
	re := getResult(Msg_success, "")
	key := new(ecdsa.PrivateKey)
	key.PublicKey.Curve = crypto.S256()
	pp, err := hex.DecodeString(kk)
	if err != nil {
		re = getResult(Msg_privatekey_error, err.Error())
		return string(re)
	}
	key.D = new(big.Int)
	key.D.SetBytes(pp)
	key.PublicKey.X, key.PublicKey.Y = key.PublicKey.Curve.ScalarBaseMult(key.D.Bytes())

	add := crypto.PubkeyToAddress(key.PublicKey)
	if ks.HasAddress(add) {
		re = getResult(Msg_account_exist, "")
		return string(re)
	}

	ks.ImportECDSA(key, pa)
	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err == nil {
		defer db.Close()
		if isExist(db, name) == false {
			err = db.Put([]byte(add.String()), []byte(name), nil)
			if err != nil {
				fmt.Println("put name err")
				re = getResult(Msg_writename_fail, err.Error())
			}
		} else {
			fmt.Println("name has existed")
			re = getResult(Msg_createuser_nameexists, "")
		}
	} else {
		re = getResult(Msg_writename_fail, err.Error())
	}

	return string(re)
}

// 导出私钥
func doExportAccount(add string, pa string) string {
	re := getResult(Msg_success, "")
	selectAccountFile := ""
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		if acc[i].Address.String() == add {
			selectAccountFile = acc[i].URL.Path
		}
	}

	if selectAccountFile != "" {
		keyjson, err := os.ReadFile(selectAccountFile)
		if err != nil {
			re = getResult(Msg_readaccountfile_fail, err.Error())
			return string(re)
		}
		key, err := keystore.DecryptKey(keyjson, pa)
		if err != nil {
			re = getResult(Msg_decryptprivatekey_fail, err.Error())
			return string(re)
		}
		// 确保导出的是需要私钥
		if key.Address.String() != add {
			fmt.Errorf("key content mismatch: have account %x, want %s", key.Address, add)
			re = getResult(Msg_addressprivatekey_mismatch, key.Address.String()+"-"+add)
			return string(re)
		}

		file, err := os.OpenFile(exportdir+add, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("can not open file in export:" + err.Error())
			re = getResult(Msg_openexportfile_fail, err.Error())
			return string(re)
		}
		defer file.Close()

		value := hex.EncodeToString(key.PrivateKey.D.Bytes())
		writer := bufio.NewWriter(file)
		writer.WriteString(value)
		writer.Flush()
	}
	return string(re)
}

func doTransfer(from string, to string, value string, pass string) string {
	var re []byte
	connect := getConnection()
	if connect == nil {
		fmt.Println("doTransfer: connection is nil")
		re = getResult(Msg_connectnode_fail, "")
		return string(re)
	}

	ctx := context.Background()
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		if from == acc[i].Address.String() {
			bigvalue, err := etherTowei(value)
			if err != nil {
				re = getResult(Msg_weitoether_fail, err.Error())
				break
			}
			var gasLimit uint64 = 21000
			nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
			if err != nil {
				re = getResult(Msg_getchaininfo_fail, err.Error())
				break
			}
			//后续可以改为可设置消费的交易
			//tx := types.NewTransaction(nounce, common.HexToAddress(to), &bigvalue, gasLimit, gasPrice, make([]byte, 0))
			tx := newDynamicTx(nounce, common.HexToAddress(to), &bigvalue, gasLimit, gasPrice, tipsPrice, chainId)
			tx, err = ks.SignTxWithPassphrase(acc[i], pass, tx, chainId)
			if err == nil {
				err := connect.SendTransaction(ctx, tx)
				if err != nil {
					fmt.Println("transfer err:" + err.Error())
					re = getResult(Msg_sendtransaction_fail, err.Error())
				} else {
					re = getResult(Msg_success, tx.Hash().String())
				}
			} else {
				fmt.Println("sign tx err")
				re = getResult(Msg_signtransaction_fail, err.Error())
			}
			break
		}
	}
	return string(re)
}

// 获取账户余额
func doGetBalance(add string) string {
	var re []byte
	connect := getConnection()
	if connect == nil {
		re = getResult(Msg_connectnode_fail, "")
		return string(re)
	}

	now := time.Now()
	fmt.Println("start:")
	fmt.Println(now.UnixMilli())
	ctx, fuc := context.WithTimeout(context.Background(), G_Timeout)
	defer fuc()
	balance, err := connect.BalanceAt(ctx, common.HexToAddress(add), nil)
	now = time.Now()
	fmt.Println("end:")
	fmt.Println(now.UnixMilli())

	if err == nil {
		type bl struct {
			Address string
			Balance string
		}
		tmp, _ := json.Marshal(bl{add, weiToether(balance.String())})
		re, _ = json.Marshal(ErrorResult{0, string(tmp)})
		return string(re)
	} else {
		re = getResult(Msg_getbalance_fail, err.Error())
		return string(re)
	}
}

func doGetFee() string {
	connect := getConnection()
	if connect == nil {
		fmt.Println("doGetFee: connection is nil")
		return ""
	}
	ctx, fuc := context.WithTimeout(context.Background(), G_Timeout)
	defer fuc()
	price, err := connect.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("get price error:" + err.Error())
		return ""
	}
	//fmt.Println("doGetFee: Gas price ：" + weiTogwei(price.String()))
	//tips, err := connect.SuggestGasTipCap(ctx)
	//if err != nil {
	//	fmt.Println("get tips error:" + err.Error())
	//	return ""
	//}
	//fmt.Println("doGetFee: Tips  ：" + weiTogwei(tips.String()))
	gas := big.NewInt(21000)
	//新的 1559 交易升级后，gasprice 包含了 base fee  和  tips 费 不需要两者相加，如果不希望给 tips，可以减掉
	//fee := big.NewInt(0)
	//fee = fee.Add(price, tips)
	//fee = fee.Mul(fee, gas)
	fee := price.Mul(price, gas)
	str := weiToether(fee.String())
	return str
}

// func doGetComponent(com string) string {
// 	fmt.Println("enter get" + com)

// 	if com == "rightmenu" {
// 		type menujson struct {
// 			Menumcss string `json:"css"`
// 			Menuhtml string `json:"html"`
// 			Menujs   string `json:"js"`
// 		}
// 		menucss, err1 := os.ReadFile(dir + "/res/common/rightmenu.css")
// 		if err1 != nil {
// 			fmt.Println("doGetComponent get rightmenu.css error")
// 			return ""
// 		}
// 		menu, err2 := os.ReadFile(dir + "/res/common/rightmenu.html")
// 		if err2 != nil {
// 			fmt.Println("doGetComponent get rightmenu.html error")
// 			return ""
// 		}

// 		menujs, err3 := os.ReadFile(dir + "/res/common/rightmenu.js")
// 		if err3 != nil {
// 			fmt.Println("doGetComponent get rightmenu.js error")
// 			return ""
// 		}
// 		result := menujson{string(menucss), string(menu), string(menujs)}
// 		re, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println("doGetComponent json error")
// 			return ""
// 		}
// 		return string(re)
// 	} else {
// 		type pagejson struct {
// 			Pagecss  string `json:"css"`
// 			Pagehtml string `json:"html"`
// 			Pagejs   string `json:"js"`
// 		}
// 		Pagecss, err1 := os.ReadFile(dir + "/res/" + com + ".css")
// 		if err1 != nil {
// 			fmt.Println("doGetComponent get " + com + ".css error")
// 		}
// 		Pagejs, err3 := os.ReadFile(dir + "/res/" + com + ".js")
// 		if err3 != nil {
// 			fmt.Println("doGetComponent get " + com + ".js error")
// 		}
// 		Pagehtml, err2 := os.ReadFile(dir + "/res/" + com + ".html")
// 		if err2 != nil {
// 			fmt.Println("doGetComponent get " + com + ".html error")
// 			return ""
// 		}

// 		result := pagejson{string(Pagecss), string(Pagehtml), string(Pagejs)}
// 		re, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println("doGetComponent json error")
// 			return ""
// 		}
// 		return string(re)
// 	}
// }

func doGetComponent(com string) string {
	type pagejson struct {
		Css  string `json:"css"`
		Html string `json:"html"`
		Js   string `json:"js"`
		Res  string `json:"res"`
	}

	result := pagejson{"", "", "", ""}
	var cssfile, htmlfile, jsfile string

	if com == "rightmenu" {
		cssfile = "/common/rightmenu.css"
		htmlfile = "/common/rightmenu.html"
		jsfile = "/common/rightmenu.js"
	} else {
		cssfile = "/" + com + ".css"
		jsfile = "/" + com + ".js"
		htmlfile = "/" + com + ".html"
	}

	result.Css = string(G_Files[cssfile])
	result.Js = string(G_Files[jsfile])
	result.Html = string(G_Files[htmlfile])
	language := G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := G_Res[language][htmlfile]
	if len(res) != 0 {
		tmp, _ := json.Marshal(res)
		result.Res = string(tmp)
	}

	re, err := json.Marshal(result)
	if err != nil {
		fmt.Println("doGetComponent json error")
		return ""
	}
	return string(re)
}

func doGetRes(com string) string {
	var re string = ""
	htmlfile := "/" + com + ".html"
	language := G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := G_Res[language][htmlfile]
	if len(res) != 0 {
		tmp, _ := json.Marshal(res)
		re = string(tmp)
	}
	return re
}

func doConfig(configs string) string {
	re := getResult(Msg_success, "")
	var config = make(map[string]string)

	json.Unmarshal([]byte(configs), &config)

	for key, value := range config {
		G_Config[key] = value
	}

	db, _ := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()
	tmp, _ := json.Marshal(G_Config)
	db.Put([]byte("config"), tmp, nil)

	/*
		for key, value := range G_Config {
			err = db.Put([]byte(key), []byte(value), nil)
			if err != nil {
				fmt.Println("put config error")
				re = getResult(Msg_writedb_fail, err.Error())
			} else {
				fmt.Println("put config:" + key + ":" + value)
			}
		}
	*/
	// 配置更新后刷新链接
	G_Etherclient = nil
	return string(re)
}

func doGetConfig(configs string) string {
	var re string
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err != nil {
		fmt.Println("open db error:" + err.Error())
		return ""
	}

	defer db.Close()

	data, err := db.Get([]byte("config"), nil)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(data, &G_Config)
	if err != nil {
		return ""
	}

	if configs != "" {
		var configName []string
		json.Unmarshal([]byte(configs), &configName)
		var result = make(map[string]string)
		for _, value := range configName {
			result[value] = G_Config[value]
		}
		tmp, _ := json.Marshal(result)
		re = string(tmp)
	} else {
		tmp, _ := json.Marshal(G_Config)
		re = string(tmp)
	}
	return string(re)
}

func doDeployContract(from string, file string, pass string) string {
	var re []byte
	if from == "" || file == "" || pass == "" {
		fmt.Println("parameter error")
		re = getResult(Msg_param_error, "")
		return string(re)
	}
	filepath := G_Config["contractDir"]
	filepath = filepath + "/" + file
	connect := getConnection()

	if connect == nil {
		fmt.Println("can not connect to node:")
		re = getResult(Msg_connectnode_fail, "")
		return string(re)
	}
	ctx := context.Background()
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		fmt.Printf("get account: %s : \r\n", acc[i].Address.String())
		if from == acc[i].Address.String() {
			bigvalue, err := etherTowei("0")
			if err != nil {
				re = getResult(Msg_weitoether_fail, err.Error())
				break
			}

			nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
			if err != nil {
				re = getResult(Msg_getchaininfo_fail, err.Error())
				break
			}
			bin, err := os.ReadFile(filepath)
			if err != nil {
				re = getResult(Msg_openfile_fail, err.Error())
				break
			}
			bytecode, err := hex.DecodeString(string(bin))
			if err != nil {
				re = getResult(Msg_decodebinfile_fail, err.Error())
				break
			}
			msg := ethereum.CallMsg{
				From:      acc[i].Address, // the sender of the 'transaction'
				To:        nil,            // the destination contract (nil for contract creation)
				Gas:       0,              // if 0, the call executes with near-infinite gas
				GasPrice:  gasPrice,       // wei <-> gas exchange ratio
				GasFeeCap: tipsPrice,      // EIP-1559 fee cap per gas.
				GasTipCap: tipsPrice,      // EIP-1559 tip per gas.
				Value:     big.NewInt(0),  // amount of wei sent along with the call
				Data:      bytecode,       // input data, usually an ABI-encoded contract method invocation
			}

			gasLimit, err := connect.EstimateGas(ctx, msg)
			if err != nil {
				re = getResult(Msg_getgaslimit_fail, err.Error())
				break
			}

			//后续可以改为可设置消费的交易
			//tx := types.NewTransaction(nounce, common.HexToAddress(to), &bigvalue, gasLimit, gasPrice, make([]byte, 0))
			tx := newDeployTx(nounce, &bigvalue, gasLimit, gasPrice, tipsPrice, chainId, bytecode)
			tx, err = ks.SignTxWithPassphrase(acc[i], pass, tx, chainId)
			if err == nil {
				err := connect.SendTransaction(ctx, tx)
				if err != nil {
					re = getResult(Msg_sendtransaction_fail, err.Error())
				} else {
					re, _ = json.Marshal(ErrorResult{0, tx.Hash().Hex()})
					//connect.TransactionReceipt(ctx,tx.Hash())
				}
			} else {
				fmt.Println("sign tx err")
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
			}
			break
		}
	}
	return string(re)
}

func doGetContract() string {
	type contract struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}
	type contracts struct {
		Contracts []contract `json:"contracts"`
	}
	var res contracts
	//账户信息
	db, err := leveldb.OpenFile(datadir+"contract", nil)
	defer db.Close()
	if err == nil {
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			var aa contract
			json.Unmarshal(it.Value(), &aa)
			res.Contracts = append(res.Contracts, aa)
		}
	}

	result, err2 := json.Marshal(res)
	if err2 == nil {
		return string(result)
	} else {
		fmt.Println(err2.Error())
	}
	return ""
}

func doImportContract(cons string) string {
	re := getResult(Msg_success, "")
	type contract struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}

	var jcons contract
	err := json.Unmarshal([]byte(cons), &jcons)
	if err == nil {
		if jcons.Name == "" {
			ctx, fuc := context.WithTimeout(context.Background(), G_Timeout)
			defer fuc()
			connect := getConnection()
			if connect != nil {
				var conAdd common.Address = common.HexToAddress(jcons.Address)
				transferFnSignature := []byte("name()")
				hash := crypto.Keccak256(transferFnSignature)
				methodID := hash[:4]
				var data []byte
				data = append(data, methodID...)
				msg := ethereum.CallMsg{
					From:      common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d"), // the sender of the 'transaction'
					To:        &conAdd,                                                           // the destination contract (nil for contract creation)
					Gas:       0,                                                                 // if 0, the call executes with near-infinite gas
					GasPrice:  nil,                                                               // wei <-> gas exchange ratio
					GasFeeCap: nil,                                                               // EIP-1559 fee cap per gas.
					GasTipCap: nil,                                                               // EIP-1559 tip per gas.
					Value:     big.NewInt(0),                                                     // amount of wei sent along with the call
					Data:      data,                                                              // input data, usually an ABI-encoded contract method invocation
				}
				result, err5 := connect.CallContract(ctx, msg, nil)
				if err5 != nil {
					fmt.Println("get contract name error:" + err5.Error())
				}
				sresult := string(result)
				sresult = strings.ReplaceAll(sresult, string([]byte{0x00}), "")
				jcons.Name = sresult
			}
		}
		db, err2 := leveldb.OpenFile(datadir+"contract", nil)
		jvalue, _ := json.Marshal(jcons)
		if err2 == nil {
			defer db.Close()
			if db.Put([]byte(jcons.Address), jvalue, nil) != nil {
				re = getResult(Msg_writedb_fail, "")
			}
		} else {
			fmt.Println("open contract db error:" + err2.Error())
			re = getResult(Msg_opendb_fail, err2.Error())
		}

	} else {
		re = getResult(Msg_param_error, "")
	}

	return string(re)
}

func doCheckNode() string {
	var re []byte
	if G_Config["nodeAddress"] != "" {
		re = getResult(Msg_success, "")
	} else {
		re = getResult(Msg_nodeconfig_null, "")
	}
	return string(re)
}

func doLoadABI(file string) string {
	var re []byte
	filepath := G_Config["contractDir"]
	filepath = filepath + "/" + file
	if filepath == "" {
		re = getResult(Msg_param_error, "")
		return string(re)
	}

	abiContent, err := os.ReadFile(filepath)
	if err != nil {
		re = getResult(Msg_openfile_fail, err.Error())
		return string(re)
	}

	re, _ = json.Marshal(ErrorResult{0, string(abiContent)})
	return string(re)
}

func doReadContract(from string, to string, funcAbi string, paremters string) string {
	var toAdd common.Address = common.HexToAddress(to)
	contractAbi, err := abi.JSON(strings.NewReader(funcAbi))
	var parems []interface{}
	var re = getResult(Msg_success, "")
	var data []byte
	var outs []string
	if err == nil {
		for _, method := range contractAbi.Methods {
			var parStrings []string
			json.Unmarshal([]byte(paremters), &parStrings)
			for i, pType := range method.Inputs {
				parems = append(parems, stringToValue(pType.Type.GetType(), parStrings[i]))
			}
			data, err = contractAbi.Pack(method.RawName, parems...)
			if err != nil {
				re = getResult(Msg_packabiparam_error, err.Error())
				return string(re)
			}
			connect := getConnection()
			if connect != nil {
				ctx, fuc := context.WithTimeout(context.Background(), G_Timeout)
				defer fuc()
				msg := ethereum.CallMsg{
					From:      common.HexToAddress(from), // the sender of the 'transaction'
					To:        &toAdd,                    // the destination contract (nil for contract creation)
					Gas:       0,                         // if 0, the call executes with near-infinite gas
					GasPrice:  nil,                       // wei <-> gas exchange ratio
					GasFeeCap: nil,                       // EIP-1559 fee cap per gas.
					GasTipCap: nil,                       // EIP-1559 tip per gas.
					Value:     big.NewInt(0),             // amount of wei sent along with the call
					Data:      data,                      // input data, usually an ABI-encoded contract method invocation
				}
				result, err5 := connect.CallContract(ctx, msg, nil)
				if err5 == nil {
					ret, _ := method.Outputs.Unpack(result)
					for _, a := range ret {
						aa := valueToString(a)
						outs = append(outs, aa)
					}
					bouts, _ := json.Marshal(outs)
					re, _ = json.Marshal(ErrorResult{0, string(bouts)})
				} else {
					re = getResult(Msg_callcontract_fail, err5.Error())
				}
			}
			break
		}
	} else {
		re = getResult(Msg_parseabi_fail, err.Error())
	}
	return string(re)
}

func doInvokeContract(from string, conAddr string, funcAbi string, param string, pass string) string {
	var re []byte
	if from == "" || pass == "" || conAddr == "" {
		re = getResult(Msg_param_error, "")
		return string(re)
	}
	connect := getConnection()
	var toAdd common.Address = common.HexToAddress(conAddr)

	if connect == nil {
		re = getResult(Msg_connectnode_fail, "")
		return string(re)
	}
	ctx := context.Background()
	acc := ks.Accounts()
	var current accounts.Account
	var isExist bool = false
	for i := 0; i < len(acc); i++ {
		fmt.Printf("get account: %s : \r\n", acc[i].Address.String())
		if from == acc[i].Address.String() {
			current = acc[i]
			isExist = true
			break
		}
	}
	if isExist == false {
		re = getResult(Msg_findaccount_fail, from)
		return string(re)
	}

	bigvalue, _ := etherTowei("0")
	nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
	if err != nil {
		re = getResult(Msg_getchaininfo_fail, err.Error())
	}

	contractAbi, err := abi.JSON(strings.NewReader(funcAbi))
	var parems []interface{}
	var parStrings []string
	var data []byte
	if err == nil {
		for _, method := range contractAbi.Methods {
			json.Unmarshal([]byte(param), &parStrings)
			for i, pType := range method.Inputs {
				parems = append(parems, stringToValue(pType.Type.GetType(), parStrings[i]))
			}
			data, err = contractAbi.Pack(method.RawName, parems...)
			if err != nil {
				re = getResult(Msg_packabiparam_error, err.Error())
				return string(re)
			}
		}
	}

	msg := ethereum.CallMsg{
		From:      current.Address, // the sender of the 'transaction'
		To:        &toAdd,          // the destination contract (nil for contract creation)
		Gas:       0,               // if 0, the call executes with near-infinite gas
		GasPrice:  gasPrice,        // wei <-> gas exchange ratio
		GasFeeCap: tipsPrice,       // EIP-1559 fee cap per gas.
		GasTipCap: tipsPrice,       // EIP-1559 tip per gas.
		Value:     big.NewInt(0),   // amount of wei sent along with the call
		Data:      data,            // input data, usually an ABI-encoded contract method invocation
	}

	gasLimit, err := connect.EstimateGas(ctx, msg)
	if err != nil {
		re = getResult(Msg_getgaslimit_fail, err.Error())
	}

	tx := newCallContractTx(nounce, toAdd, &bigvalue, gasLimit, gasPrice, tipsPrice, chainId, data)
	tx, err = ks.SignTxWithPassphrase(current, pass, tx, chainId)
	if err == nil {
		err := connect.SendTransaction(ctx, tx)
		if err != nil {
			re, _ = json.Marshal(ErrorResult{1, err.Error()})
		} else {
			re, _ = json.Marshal(ErrorResult{0, tx.Hash().Hex()})
			//connect.TransactionReceipt(ctx,tx.Hash())
		}

	} else {
		re = getResult(Msg_signtransaction_fail, err.Error())
	}

	return string(re)
}

func doQueryTx(hs string) string {
	type txinfo struct {
		Transaction        *types.Transaction
		TransactionReceipt *types.Receipt
	}
	var re []byte
	var info txinfo
	txhash := common.HexToHash(hs)
	connect := getConnection()
	ctx, fuc := context.WithTimeout(context.Background(), G_Timeout)
	defer fuc()
	tx, _, err := connect.TransactionByHash(ctx, txhash)
	if err == nil {
		//re, _ = tx.MarshalJSON()
		info.Transaction = tx
		receipt, err2 := connect.TransactionReceipt(ctx, txhash)
		if err2 == nil {
			info.TransactionReceipt = receipt
		}
	} else {
		return err.Error()
	}
	re, _ = json.Marshal(info)
	return string(re)
}

func doSignature(info string) string {
	var re []byte
	var param = make(map[string]string)
	err := json.Unmarshal([]byte(info), &param)
	if err == nil {
		acc := ks.Accounts()
		var i = 0
		for i = 0; i < len(acc); i++ {
			if acc[i].Address.String() == param["Address"] {
				signature, err2 := ks.SignHashWithPassphrase(acc[i], param["Ass"], crypto.Keccak256([]byte(param["Info"])))
				if err2 == nil {
					re, _ = json.Marshal(ErrorResult{0, common.Bytes2Hex(signature)})
				} else {
					re = getResult(Msg_fail, err2.Error())
				}
				break
			}
		}
	} else {
		re = getResult(Msg_param_error, err.Error())
	}

	return string(re)
}

func doGetAbiFile() string {
	re := getResult(Msg_success, "")
	var conf = [1]string{"contractDir"}
	par, _ := json.Marshal(conf)
	doGetConfig(string(par))
	var dir = G_Config["contractDir"]
	if dir != "" {
		names, err := getFiles(dir, ".abi")
		if err == nil {
			tmp, _ := json.Marshal(names)
			re, _ = json.Marshal(ErrorResult{0, string(tmp)})
		} else {
			re = getResult(Msg_openfile_fail, err.Error())
		}
	} else {
		re = getResult(Msg_contractdirconfig_null, "")
	}
	return string(re)
}

func doGetBinFile() string {
	re := getResult(Msg_success, "")
	var conf = [1]string{"contractDir"}
	par, _ := json.Marshal(conf)
	doGetConfig(string(par))
	var dir = G_Config["contractDir"]
	if dir != "" {
		names, err := getFiles(dir, ".bin")
		if err == nil {
			tmp, _ := json.Marshal(names)
			re, _ = json.Marshal(ErrorResult{0, string(tmp)})
		} else {
			re = getResult(Msg_openfile_fail, err.Error())
		}
	} else {
		re = getResult(Msg_contractdirconfig_null, "")
	}
	return string(re)
}

func doGetMsg(code string) string {
	n, _ := strconv.Atoi(code)
	lang := G_Config["language"]
	if lang == "" {
		lang = "zh-CN"
	}
	msg := G_Msg[lang][uint(n)]
	if msg == "" {
		msg = "null"
	}

	return msg
}

func doSaveDataSpace(name string, content string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		value, err2 := db.Get([]byte(name), nil)
		if (err2 == errors.ErrNotFound) || (err2 == nil && len(value) == 0) {
			encryptData, err1 := keystore.EncryptDataV3([]byte(content), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
			if err1 == nil {
				data, _ := json.Marshal(encryptData)
				err3 := db.Put([]byte(name), data, nil)
				if err3 != nil {
					re = getResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = getResult(Msg_fail, err1.Error())
			}
		} else if err2 == nil {
			re = getResult(Msg_exist, string(value))
		} else {
			re = getResult(Msg_fail, err2.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}

	return string(re)
}

func doChangeDataSpace(name string, content string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		encryptData, err1 := keystore.EncryptDataV3([]byte(content), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
		if err1 == nil {
			data, _ := json.Marshal(encryptData)
			err2 := db.Put([]byte(name), data, nil)
			if err2 != nil {
				re = getResult(Msg_writedb_fail, err2.Error())
			}
		} else {
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doDeleteDataSpace(name string, auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		err2 := db.Delete([]byte(name), nil)
		if err2 != nil {
			re = getResult(Msg_fail, err2.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doGetDataSpace(auth string) string {
	re := getResult(Msg_success, "")
	check := checkDataSpaceAuth(auth)
	if check != nil {
		re = getResult(Msg_fail, check.Error())
		return string(re)
	}
	data := make(map[string]string)
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			if string(it.Key()) == "data_space_init" {
				continue
			}
			var tmp keystore.CryptoJSON
			if json.Unmarshal(it.Value(), &tmp) == nil {
				dedata, err1 := keystore.DecryptDataV3(tmp, auth)
				if err1 == nil {
					data[string(it.Key())] = string(dedata)
				} else {
					data[string(it.Key())] = err1.Error()
				}
			}
		}
		msg, _ := json.Marshal(data)
		re, _ = json.Marshal(ErrorResult{Msg_success, string(msg)})
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func doCheckDataSpace() string {
	re := getResult(Msg_success, "")
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		//have not init
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			re = getResult(Msg_space_notinited, "")
		} else if err1 == nil { //have inited
			re = getResult(Msg_space_inited, "")
		} else { //other error
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}

	return string(re)
}

func doInitDataSpace(auth string) string {
	re := getResult(Msg_success, "")
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		//have not init
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			encryptData, err2 := keystore.EncryptDataV3([]byte("data_space_init"), []byte(auth), keystore.StandardScryptN, keystore.StandardScryptP)
			if err2 == nil {
				data, _ := json.Marshal(encryptData)
				err3 := db.Put([]byte("data_space_init"), data, nil)
				if err3 != nil {
					re = getResult(Msg_writedb_fail, err3.Error())
				}
			} else {
				re = getResult(Msg_fail, err2.Error())
			}
		} else if err1 == nil { //have inited
			re = getResult(Msg_space_inited, "")
		} else { //other error
			re = getResult(Msg_fail, err1.Error())
		}
	} else {
		re = getResult(Msg_fail, err.Error())
	}
	return string(re)
}

func checkDataSpaceAuth(auth string) error {
	db, err := leveldb.OpenFile(datadir+"dataspace", nil)
	if db != nil {
		defer db.Close()
	}
	if err == nil {
		tag, err1 := db.Get([]byte("data_space_init"), nil)
		if (err1 == errors.ErrNotFound) || (err1 == nil && len(tag) == 0) {
			return errors.New(doGetMsg(strconv.Itoa(int(Msg_space_notinited))))
		} else if err1 == nil { //have inited
			var tmp keystore.CryptoJSON
			if json.Unmarshal(tag, &tmp) == nil {
				value, err2 := keystore.DecryptDataV3(tmp, auth)
				if err2 == nil && string(value) == "data_space_init" {
					return nil
				} else {
					return errors.New(doGetMsg(strconv.Itoa(int(Msg_password_error))))
				}
			}
		} else { //other error
			return err1
		}
	}
	return nil
}
