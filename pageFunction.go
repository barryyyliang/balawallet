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

//新创建一个账号
func doNewAccount(pass string, walletName string) string {
	re := getResult(0, "")
	//scryptN := keystore.StandardScryptN
	//scryptP := keystore.StandardScryptP
	//account, err := keystore.StoreKey(accountdir, pass, scryptN, scryptP)
	account, err := ks.NewAccount(pass)
	if err != nil {
		fmt.Println("failed to generate new account")
		re = getResult(1, err.Error())
		return string(re)
	}

	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	if isAddressExist(db, account.Address.String(), walletName) == false {
		err = db.Put([]byte(account.Address.String()), []byte(walletName), nil)
		if err != nil {
			fmt.Println("put name err")
			re = getResult(2, err.Error())
		}
	} else {
		fmt.Println("name has existed")
		re = getResult(3, "")
	}

	return string(re)

}

//日志接口
func doWebLog(para string) {
	fmt.Println(para)
}

//获取账号列表
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
	defer db.Close()
	if err == nil {
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
	re := getResult(0, "")
	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	if isAddressExist(db, address, name) == false {
		err = db.Put([]byte(address), []byte(name), nil)
		if err != nil {
			re = getResult(4, err.Error())
		} else {
			fmt.Println("put name:" + address + ":" + name)
		}
	} else {
		re = getResult(5, "")
	}
	return string(re)
}

func doGetAccountName(address string) string {
	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()
	if err == nil {
		fmt.Println(address)
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

//通过私钥导入账号
func doImportAccount(pa string, name string, kk string) string {
	re := getResult(0, "")
	fmt.Println("importAccount")
	key := new(ecdsa.PrivateKey)
	key.PublicKey.Curve = crypto.S256()
	pp, err := hex.DecodeString(kk)
	if err != nil {
		re = getResult(6, err.Error())
		return string(re)
	}
	key.D = new(big.Int)
	key.D.SetBytes(pp)
	key.PublicKey.X, key.PublicKey.Y = key.PublicKey.Curve.ScalarBaseMult(key.D.Bytes())
	ks.ImportECDSA(key, pa)
	return string(re)
}

//导出私钥
func doExportAccount(add string, pa string) string {
	re := getResult(0, "")
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
			re = getResult(7, err.Error())
			return string(re)
		}
		key, err := keystore.DecryptKey(keyjson, pa)
		if err != nil {
			re = getResult(8, err.Error())
			return string(re)
		}
		// 确保导出的是需要私钥
		if key.Address.String() != add {
			fmt.Errorf("key content mismatch: have account %x, want %s", key.Address, add)
			re = getResult(9, key.Address.String()+"-"+add)
			return string(re)
		}

		file, err := os.OpenFile(exportdir+add, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("can not open file in export:" + err.Error())
			re = getResult(10, err.Error())
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
		re, _ = json.Marshal(ErrorResult{1, "can not coonnect to node"})
		return string(re)
	}

	ctx := context.Background()
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		fmt.Printf("get account: %s : \r\n", acc[i].Address.String())
		if from == acc[i].Address.String() {
			bigvalue, err := etherTowei(value)
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
				break

			}
			var gasLimit uint64 = 21000
			nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
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
					re, _ = json.Marshal(ErrorResult{1, err.Error()})
				} else {
					fmt.Println("transfer success")
					re, _ = json.Marshal(ErrorResult{0, ""})
				}

				json, err2 := tx.MarshalJSON()
				if err2 == nil {
					fmt.Println(string(json))
				} else {
					fmt.Println("marshaljson tx err")
				}
			} else {
				fmt.Println("sign tx err")
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
			}
			break
		}
	}
	fmt.Println(string(re))
	return string(re)
}

//获取账户余额
func doGetBalance(add string) string {
	var re []byte
	fmt.Println("doGetBalance: address is " + add)
	if add == "" {
		re, _ = json.Marshal(ErrorResult{1, "address is null"})
		return string(re)
	}
	connect := getConnection()
	if connect == nil {
		re, _ = json.Marshal(ErrorResult{1, "can not get conneciton to node"})
		return string(re)
	}

	now := time.Now()
	fmt.Println("start:")
	fmt.Println(now.UnixMilli())
	ctx, fuc := context.WithTimeout(context.Background(), 3*time.Second)
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
		re, _ = json.Marshal(ErrorResult{1, "get balance error:" + err.Error()})
		return string(re)
	}
}

func doGetFee() string {
	connect := getConnection()
	if connect == nil {
		fmt.Println("doGetFee: connection is nil")
		return ""
	}
	ctx, fuc := context.WithTimeout(context.Background(), 3*time.Second)
	defer fuc()
	price, err := connect.SuggestGasPrice(ctx)
	if err != nil {
		fmt.Println("get price error:" + err.Error())
		return ""
	}
	fmt.Println("price:" + weiToether(price.String()))
	tips, err := connect.SuggestGasTipCap(ctx)
	if err != nil {
		fmt.Println("get tips error:" + err.Error())
		return ""
	}
	fmt.Println("tips:" + weiToether(tips.String()))
	gas := big.NewInt(21000)
	fee := big.NewInt(0)
	fee = fee.Add(price, tips)
	fee = fee.Mul(fee, gas)
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
	fmt.Println("enter get" + com)

	if com == "rightmenu" {
		type menujson struct {
			Menumcss string `json:"css"`
			Menuhtml string `json:"html"`
			Menujs   string `json:"js"`
		}
		menucss := G_Files["/common/rightmenu.css"]
		menu := G_Files["/common/rightmenu.html"]
		menujs := G_Files["/common/rightmenu.js"]

		result := menujson{string(menucss), string(menu), string(menujs)}
		re, err := json.Marshal(result)
		if err != nil {
			fmt.Println("doGetComponent json error")
			return ""
		}
		return string(re)
	} else {
		type pagejson struct {
			Pagecss  string `json:"css"`
			Pagehtml string `json:"html"`
			Pagejs   string `json:"js"`
		}
		Pagecss := G_Files["/"+com+".css"]
		Pagejs := G_Files["/"+com+".js"]
		Pagehtml := G_Files["/"+com+".html"]

		result := pagejson{string(Pagecss), string(Pagehtml), string(Pagejs)}
		re, err := json.Marshal(result)
		if err != nil {
			fmt.Println("doGetComponent json error")
			return ""
		}
		return string(re)
	}
}

func doConfig(configs string) string {
	fmt.Println("enter doConfig")
	re, _ := json.Marshal(ErrorResult{0, "success"})
	var config = make(map[string]string)

	json.Unmarshal([]byte(configs), &config)

	G_Config = config

	fmt.Println(config)

	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	for key, value := range G_Config {
		err = db.Put([]byte(key), []byte(value), nil)
		if err != nil {
			fmt.Println("put config error")
			re, _ = json.Marshal(ErrorResult{1, "write db error"})
		} else {
			fmt.Println("put config:" + key + ":" + value)
		}
	}

	// 配置更新后刷新链接
	G_Etherclient = nil

	return string(re)
}

func doGetConfig(configs string) string {
	fmt.Println(configs)
	db, err := leveldb.OpenFile(datadir+"db", nil)
	if err != nil {
		fmt.Println("open db error:" + err.Error())
		return ""
	}

	defer db.Close()

	var configName []string
	json.Unmarshal([]byte(configs), &configName)
	for _, key := range configName {
		data, err := db.Get([]byte(key), nil)
		if err == nil {
			G_Config[key] = string(data)
		} else {
			fmt.Println("get erro:" + err.Error())
		}
	}
	re, _ := json.Marshal(G_Config)
	return string(re)
}

func doDeployContract(from string, file string, pass string) string {
	var re []byte
	if from == "" || file == "" || pass == "" {
		fmt.Println("parameter error")
		re, _ = json.Marshal(ErrorResult{1, "parameter error"})
		return string(re)
	}
	filepath := G_Config["contractDir"]
	filepath = filepath + "/" + file
	connect := getConnection()

	if connect == nil {
		fmt.Println("can not connect to node:")
		re, _ = json.Marshal(ErrorResult{1, "can not connet to node"})
		return string(re)
	}
	ctx := context.Background()
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		fmt.Printf("get account: %s : \r\n", acc[i].Address.String())
		if from == acc[i].Address.String() {
			bigvalue, err := etherTowei("0")
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
				break
			}

			nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
				break
			}
			bin, err := os.ReadFile(filepath)
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
				break
			}
			bytecode, err := hex.DecodeString(string(bin))
			if err != nil {
				re, _ = json.Marshal(ErrorResult{1, "can not read bin file:" + filepath})
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
				fmt.Println("get gaslimit:" + strconv.FormatInt(int64(gasLimit), 10))
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
				break
			}

			//后续可以改为可设置消费的交易
			//tx := types.NewTransaction(nounce, common.HexToAddress(to), &bigvalue, gasLimit, gasPrice, make([]byte, 0))
			tx := newDeployTx(nounce, &bigvalue, gasLimit, gasPrice, tipsPrice, chainId, bytecode)
			tx, err = ks.SignTxWithPassphrase(acc[i], pass, tx, chainId)
			if err == nil {
				err := connect.SendTransaction(ctx, tx)
				if err != nil {
					re, _ = json.Marshal(ErrorResult{1, err.Error()})
				} else {
					fmt.Println("deploy success")
					re, _ = json.Marshal(ErrorResult{0, tx.Hash().Hex()})
					//connect.TransactionReceipt(ctx,tx.Hash())
				}

				json, err2 := tx.MarshalJSON()
				if err2 == nil {
					fmt.Println(string(json))
				} else {
					fmt.Println("marshaljson tx err")
				}
			} else {
				fmt.Println("sign tx err")
				re, _ = json.Marshal(ErrorResult{1, err.Error()})
			}
			break
		}
	}
	fmt.Println(string(re))
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
			fmt.Println("contract:" + string(it.Value()))
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
	re, _ := json.Marshal(ErrorResult{0, "success"})
	fmt.Println("doImportContract:" + cons)
	type contract struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
	}

	var jcons contract
	err := json.Unmarshal([]byte(cons), &jcons)
	if err == nil {
		if jcons.Name == "" {
			ctx, fuc := context.WithTimeout(context.Background(), 3*time.Second)
			defer fuc()
			connect := getConnection()
			if connect != nil {
				var conAdd common.Address = common.HexToAddress(jcons.Address)
				transferFnSignature := []byte("name()")
				hash := crypto.Keccak256(transferFnSignature)
				methodID := hash[:4]
				fmt.Println(string(methodID))
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
				fmt.Println("get contract name:" + string(sresult))
				fmt.Println(data)
				jcons.Name = sresult
			}
		}
		db, err2 := leveldb.OpenFile(datadir+"contract", nil)
		defer db.Close()
		jvalue, _ := json.Marshal(jcons)
		if err2 == nil {
			if db.Put([]byte(jcons.Address), jvalue, nil) != nil {
				fmt.Println("write db error:" + string(jvalue))
				re, _ = json.Marshal(ErrorResult{1, "write db error"})
			}
		} else {
			fmt.Println("open contract db error:" + err2.Error())
			re, _ = json.Marshal(ErrorResult{1, "open db error"})
		}

	} else {
		fmt.Println("get wrong parameter:" + cons)
		re, _ = json.Marshal(ErrorResult{1, "get wrong parameter"})
	}

	return string(re)
}

func doCheckNode() string {
	var re []byte
	if G_Config["nodeAddress"] != "" {
		re, _ = json.Marshal(ErrorResult{0, "have config blockchain node"})
	} else {
		re, _ = json.Marshal(ErrorResult{1, "have not config blockchain node"})
	}
	return string(re)
}

func doLoadABI(file string) string {
	var re []byte
	filepath := G_Config["contractDir"]
	filepath = filepath + "/" + file
	if filepath == "" {
		re, _ = json.Marshal(ErrorResult{1, "parameter error"})
		return string(re)
	}

	abiContent, err := os.ReadFile(filepath)
	if err != nil {
		re, _ = json.Marshal(ErrorResult{2, err.Error()})
		return string(re)
	}

	re, _ = json.Marshal(ErrorResult{0, string(abiContent)})
	return string(re)
}

func doReadContract(from string, to string, funcAbi string, paremters string) string {
	fmt.Println("doReadContract")
	fmt.Println(funcAbi)
	fmt.Println(paremters)
	var toAdd common.Address = common.HexToAddress(to)
	contractAbi, err := abi.JSON(strings.NewReader(funcAbi))
	var parems []interface{}
	var re ErrorResult = ErrorResult{0, ""}
	var data []byte
	var outs []string
	if err == nil {
		for _, method := range contractAbi.Methods {
			var parStrings []string
			json.Unmarshal([]byte(paremters), &parStrings)
			for i, pType := range method.Inputs {
				parems = append(parems, stringToValue(pType.Type.GetType(), parStrings[i]))
				fmt.Println(parStrings[i])
			}
			data, err = contractAbi.Pack(method.RawName, parems...)
			if err != nil {
				re.Code = 1
				re.Msg = "pack abi data err:" + err.Error()
				tmp, _ := json.Marshal(re)
				return string(tmp)
			}
			connect := getConnection()
			if connect != nil {
				ctx, fuc := context.WithTimeout(context.Background(), 3*time.Second)
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
						fmt.Println(aa)
					}
					bouts, _ := json.Marshal(outs)
					re.Msg = string(bouts)
				} else {
					re.Code = 1
					re.Msg = "call contract error:" + err5.Error()
				}
				fmt.Println(result)
			}
			break
		}
	} else {
		re.Code = 1
		re.Msg = "parse abi error"
	}
	fmt.Println(re.Msg)
	fmt.Println(re.Code)
	tmp, _ := json.Marshal(re)
	return string(tmp)
}

func doInvokeContract(from string, conAddr string, funcAbi string, param string, pass string) string {
	var re []byte
	if from == "" || pass == "" || conAddr == "" {
		fmt.Println("parameter error")
		re, _ = json.Marshal(ErrorResult{1, "parameter error"})
		return string(re)
	}
	connect := getConnection()
	var toAdd common.Address = common.HexToAddress(conAddr)

	if connect == nil {
		fmt.Println("can not connect to node:")
		re, _ = json.Marshal(ErrorResult{1, "can not connet to node"})
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
		re, _ = json.Marshal(ErrorResult{1, "can not find account"})
		return string(re)
	}

	bigvalue, err := etherTowei("0")
	nounce, gasPrice, tipsPrice, chainId, err := getChainInfo(connect, from)
	if err != nil {
		re, _ = json.Marshal(ErrorResult{1, err.Error()})
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
				fmt.Println(parStrings[i])
			}
			data, err = contractAbi.Pack(method.RawName, parems...)
			if err != nil {
				tmp, _ := json.Marshal(ErrorResult{1, "pack abi data err:" + err.Error()})
				return string(tmp)
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
		fmt.Println("get gaslimit:" + strconv.FormatInt(int64(gasLimit), 10))
		re, _ = json.Marshal(ErrorResult{1, err.Error()})
	}

	tx := newCallContractTx(nounce, toAdd, &bigvalue, gasLimit, gasPrice, tipsPrice, chainId, data)
	tx, err = ks.SignTxWithPassphrase(current, pass, tx, chainId)
	if err == nil {
		err := connect.SendTransaction(ctx, tx)
		if err != nil {
			re, _ = json.Marshal(ErrorResult{1, err.Error()})
		} else {
			fmt.Println("invoke success")
			re, _ = json.Marshal(ErrorResult{0, tx.Hash().Hex()})
			//connect.TransactionReceipt(ctx,tx.Hash())
		}

		json, err2 := tx.MarshalJSON()
		if err2 == nil {
			fmt.Println(string(json))
		} else {
			fmt.Println("marshaljson tx err")
		}
	} else {
		fmt.Println("sign tx err")
		re, _ = json.Marshal(ErrorResult{1, err.Error()})
	}

	fmt.Println(string(re))
	return string(re)
}

func doQueryTx(hs string) string {
	type txinfo struct {
		Transaction        *types.Transaction
		TransactionReceipt *types.Receipt
	}
	fmt.Println("doQueryTx:" + hs)
	var re []byte
	var info txinfo
	txhash := common.HexToHash(hs)
	connect := getConnection()
	ctx, fuc := context.WithTimeout(context.Background(), 3*time.Second)
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
	fmt.Println("server: doSignature :" + info)
	var re []byte
	var param = make(map[string]string)
	err := json.Unmarshal([]byte(info), &param)
	fmt.Println("server: doSignature 0:" + info)
	if err == nil {
		fmt.Println("server: doSignature :1")
		acc := ks.Accounts()
		var i = 0
		for i = 0; i < len(acc); i++ {
			fmt.Println("server: doSignature :2")
			if acc[i].Address.String() == param["Address"] {
				fmt.Println("server: doSignature :3.1")
				signature, err2 := ks.SignHashWithPassphrase(acc[i], param["Ass"], crypto.Keccak256([]byte(param["Info"])))
				if err2 == nil {
					fmt.Println("server: doSignature :3")
					re, _ = json.Marshal(ErrorResult{0, common.Bytes2Hex(signature)})
				} else {
					fmt.Println("server: doSignature :4" + err2.Error())
					re, _ = json.Marshal(ErrorResult{1, err2.Error()})
				}
				break
			}
		}
	} else {
		fmt.Println("server: doSignature :5" + err.Error())
		re, _ = json.Marshal(ErrorResult{1, err.Error()})
	}

	return string(re)
}

func doGetAbiFile() string {
	re := ErrorResult{0, ""}
	var conf = [1]string{"contractDir"}
	par, _ := json.Marshal(conf)
	doGetConfig(string(par))
	var dir = G_Config["contractDir"]
	if dir != "" {
		names, err := getFiles(dir, ".abi")
		if err == nil {
			tmp, _ := json.Marshal(names)
			re.Msg = string(tmp)
		} else {
			re.Code = 1
			re.Msg = err.Error()
		}
	} else {
		re.Code = 1
		re.Msg = "contract dir name not config"
	}

	bre, _ := json.Marshal(re)
	return string(bre)
}

func doGetBinFile() string {
	re := ErrorResult{0, ""}
	var conf = [1]string{"contractDir"}
	par, _ := json.Marshal(conf)
	doGetConfig(string(par))
	var dir = G_Config["contractDir"]
	if dir != "" {
		names, err := getFiles(dir, ".bin")
		if err == nil {
			tmp, _ := json.Marshal(names)
			re.Msg = string(tmp)
		} else {
			re.Code = 1
			re.Msg = err.Error()
		}
	} else {
		re.Code = 1
		re.Msg = "contract dir name not config"
	}
	bre, _ := json.Marshal(re)
	return string(bre)
}
