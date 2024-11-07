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

func doDeleteContract(address string) string {
	re := getResult(Msg_success, "")
	db, err := leveldb.OpenFile(datadir+"contract", nil)
	if err == nil {
		defer db.Close()
		if db.Delete([]byte(address), nil) != nil {
			re = getResult(Msg_writedb_fail, "")
		}
	} else {
		fmt.Println("open contract db error:" + err.Error())
		re = getResult(Msg_opendb_fail, err.Error())
	}
	return string(re)
}

func doImportContract(cons string) string {
	re := getResult(Msg_success, "")
	type contract struct {
		Type    string `json:"type"`
		Address string `json:"address"`
		Name    string `json:"name"`
		From    string `json:"from"`
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
				contractAbi, _ := abi.JSON(strings.NewReader(G_Abis[jcons.Type]))
				transferFnSignature := []byte("name()")
				hash := crypto.Keccak256(transferFnSignature)
				methodID := hash[:4]
				method, _ := contractAbi.MethodById(methodID)

				var data []byte
				data = append(data, methodID...)
				msg := ethereum.CallMsg{
					From:      common.HexToAddress(jcons.From), // the sender of the 'transaction'
					To:        &conAdd,                         // the destination contract (nil for contract creation)
					Gas:       0,                               // if 0, the call executes with near-infinite gas
					GasPrice:  nil,                             // wei <-> gas exchange ratio
					GasFeeCap: nil,                             // EIP-1559 fee cap per gas.
					GasTipCap: nil,                             // EIP-1559 tip per gas.
					Value:     big.NewInt(0),                   // amount of wei sent along with the call
					Data:      data,                            // input data, usually an ABI-encoded contract method invocation
				}
				result, err5 := connect.CallContract(ctx, msg, nil)
				fmt.Println("get contract name:" + string(result))
				if err5 == nil {
					ret, _ := method.Outputs.Unpack(result)
					for _, a := range ret {
						aa := valueToString(a)
						fmt.Println("get contract name2:" + aa)
						jcons.Name = jcons.Name + aa
					}
				} else {
					fmt.Println("get contract name error:" + err5.Error())
				}
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
