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
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	//"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
)

type ErrorResult struct {
	Code int
	Err  string
}

//页面切换
func doChangePage(page string, param string) {
	fmt.Println(resdir + page)
	content, err := os.ReadFile(resdir + page)
	if err != nil {
		fmt.Println("can not read html")
	}
	win.SetHtml(string(content))
}

func isAddressExist(db *leveldb.DB, adress string, name string) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		if string(it.Value()) == name {
			return true
		}
	}
	return false
}

//新创建一个账号
func doNewAccount(pass string, walletName string) string {
	var re []byte
	re, _ = json.Marshal(ErrorResult{0, "success"})
	fmt.Println("password:" + pass)
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	account, err := keystore.StoreKey(accountdir, pass, scryptN, scryptP)
	if err != nil {
		fmt.Println("failed to generate new account")
		re, _ = json.Marshal(ErrorResult{1, err.Error()})
		return string(re)
	}

	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	if isAddressExist(db, account.Address.String(), walletName) == false {
		err = db.Put([]byte(account.Address.String()), []byte(walletName), nil)
		if err != nil {
			fmt.Println("put name err")
			re, _ = json.Marshal(ErrorResult{1, "create success, but write name error"})
		}
	} else {
		fmt.Println("name has existed")
		re, _ = json.Marshal(ErrorResult{1, "create success, but name is exist"})
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

func doChangeAccountName(address string, name string) {
	//账户信息存库
	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	if isAddressExist(db, address, name) == false {
		err = db.Put([]byte(address), []byte(name), nil)
		if err != nil {
			fmt.Println("put name err")
		} else {
			fmt.Println("put name:" + address + ":" + name)
		}
	} else {
		fmt.Println("name has existed")
	}
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
func doImportAccount(pa string, name string, kk string) {
	fmt.Println("importAccount")
	key := new(ecdsa.PrivateKey)
	key.PublicKey.Curve = crypto.S256()
	pp, err := hex.DecodeString(kk)
	if err != nil {
		fmt.Println("bad private key string")
	}
	key.D = new(big.Int)
	key.D.SetBytes(pp)
	key.PublicKey.X, key.PublicKey.Y = key.PublicKey.Curve.ScalarBaseMult(key.D.Bytes())
	ks.ImportECDSA(key, pa)
}

//导出私钥
func doExportAccount(add string, pa string) {
	//fmt.Println("exportAccount")
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
			return
		}
		key, err := keystore.DecryptKey(keyjson, pa)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// Make sure we're really operating on the requested key (no swap attacks)
		if key.Address.String() != add {
			fmt.Errorf("key content mismatch: have account %x, want %s", key.Address, add)
			return
		}

		file, err := os.OpenFile(exportdir+add, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("can not open file in export:" + err.Error())
			return
		}
		defer file.Close()

		value := hex.EncodeToString(key.PrivateKey.D.Bytes())
		writer := bufio.NewWriter(file)
		writer.WriteString(value)
		writer.Flush()
	}
}

func newDynamicTx(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tipsPrice *big.Int, chain *big.Int) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: tipsPrice,
		ChainID:   chain,
	})
}

func newDeployTx(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tipsPrice *big.Int, chain *big.Int, data []byte) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Value:     amount,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: tipsPrice,
		ChainID:   chain,
		Data:      data,
	})
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
	fmt.Println("doGetBalance: address is " + add)
	if add == "" {
		fmt.Println("doGetBalance: address is null")
		return ""
	}
	connect := getConnection()
	if connect == nil {
		fmt.Println("doGetBalance: connection is nil")
		return ""
	}

	now := time.Now()
	fmt.Println("start:")
	fmt.Println(now.UnixMilli())
	balance, err := connect.BalanceAt(context.Background(), common.HexToAddress(add), nil)
	now = time.Now()
	fmt.Println("end:")
	fmt.Println(now.UnixMilli())

	if err == nil {
		return weiToether(balance.String())
	} else {
		fmt.Println("get balance error:" + err.Error())
		return ""
	}
}

func weiTox(unittag int, wei string) string {
	ballength := len(wei)
	var balance string
	fmt.Println("wei:" + wei)
	if ballength > unittag {
		var first []byte = make([]byte, ballength-unittag)
		copy(first, []byte(wei)[0:ballength-unittag])
		//fmt.Println("first:" + string(first))
		var second []byte = make([]byte, unittag)
		copy(second, []byte(wei)[ballength-unittag:])
		//fmt.Println("second:" + string(second))
		secondStr := strings.TrimRight(string(second), "0")
		if secondStr != "" {
			balance = string(first) + "." + secondStr
		} else {
			balance = string(first)
		}
	} else {
		wei = strings.TrimRight(wei, "0")
		if wei != "" {
			gap := unittag - ballength
			balance = "0."

			for i := 0; i < gap; i++ {
				balance = balance + "0"
			}
			balance = balance + wei
		} else {
			balance = "0"
		}
	}
	return balance
}

func weiToether(wei string) string {
	return weiTox(18, wei)
}

func weiTogwei(wei string) string {
	return weiTox(9, wei)
}

func etherTowei(ether string) (big.Int, error) {
	// 要调用isnumber判断是否为正确数字格式
	if isNumber(ether) {
		strs := strings.Split(ether, ".")
		var balstr string
		if len(strs) == 1 {
			if balstr != "0" {
				balstr = strs[0] + "000000000000000000"
			}
		} else {
			lens := len(strs[1])

			for i := 0; i < 18-lens; i++ {
				strs[1] = strs[1] + "0"
			}

			if "0" == strs[0] {
				balstr = strings.TrimLeft(strs[1], "0")
			} else {
				balstr = strs[0] + strs[1]
			}
		}

		var wei big.Int = *big.NewInt(0)
		wei.SetString(balstr, 10)
		return wei, nil
	}
	return *big.NewInt(0), fmt.Errorf("bad number")

}

//判断字符串是否为10进制数字，，每个字符必须为数字【1--9】，最多有一个小数点”.“
func isNumber(value string) bool {
	//todo
	result, err := strconv.ParseFloat(value, 64)
	fmt.Println(result)
	if err != nil && err.(*strconv.NumError).Err == strconv.ErrSyntax {
		fmt.Println("bad number:" + err.Error())
		return false
	}
	return true
}

func doGetFee() string {
	connect := getConnection()
	if connect == nil {
		fmt.Println("doGetFee: connection is nil")
		return ""
	}
	ctx := context.Background()

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

func doGetComponent(com string) string {
	fmt.Println("enter get" + com)

	if com == "menu" {
		type menujson struct {
			Menumcss string `json:"css"`
			Menuhtml string `json:"html"`
		}
		menucss, err1 := os.ReadFile(dir + "/res/common/menu.css")
		if err1 != nil {
			fmt.Println("doGetComponent get menu.css error")
			return ""
		}
		menu, err2 := os.ReadFile(dir + "/res/common/menu.html")
		if err2 != nil {
			fmt.Println("doGetComponent get menu.html error")
			return ""
		}
		result := menujson{string(menucss), string(menu)}
		re, err := json.Marshal(result)
		if err != nil {
			fmt.Println("doGetComponent json error")
			return ""
		}
		return string(re)
	} else if com == "rightmenu" {
		type menujson struct {
			Menumcss string `json:"css"`
			Menuhtml string `json:"html"`
			Menujs   string `json:"js"`
		}
		menucss, err1 := os.ReadFile(dir + "/res/common/rightmenu.css")
		if err1 != nil {
			fmt.Println("doGetComponent get rightmenu.css error")
			return ""
		}
		menu, err2 := os.ReadFile(dir + "/res/common/rightmenu.html")
		if err2 != nil {
			fmt.Println("doGetComponent get rightmenu.html error")
			return ""
		}

		menujs, err3 := os.ReadFile(dir + "/res/common/rightmenu.js")
		if err3 != nil {
			fmt.Println("doGetComponent get rightmenu.js error")
			return ""
		}
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
		Pagecss, err1 := os.ReadFile(dir + "/res/" + com + ".css")
		if err1 != nil {
			fmt.Println("doGetComponent get " + com + ".css error")
		}
		Pagejs, err3 := os.ReadFile(dir + "/res/" + com + ".js")
		if err3 != nil {
			fmt.Println("doGetComponent get " + com + ".js error")
		}
		Pagehtml, err2 := os.ReadFile(dir + "/res/" + com + ".html")
		if err2 != nil {
			fmt.Println("doGetComponent get " + com + ".html error")
			return ""
		}

		result := pagejson{string(Pagecss), string(Pagehtml), string(Pagejs)}
		re, err := json.Marshal(result)
		if err != nil {
			fmt.Println("doGetComponent json error")
			return ""
		}
		return string(re)
	}
}

func doConfig(configs string) {
	fmt.Println("enter doConfig")

	var config = make(map[string]string)

	json.Unmarshal([]byte(configs), &config)

	G_Config = config

	fmt.Println(config)

	db, err := leveldb.OpenFile(datadir+"db", nil)
	defer db.Close()

	for key, value := range G_Config {
		err = db.Put([]byte(key), []byte(value), nil)
		if err != nil {
			fmt.Println("put configi err")
		} else {
			fmt.Println("put config:" + key + ":" + value)
		}
	}

	// 配置更新后刷新链接
	G_Etherclient = nil
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

func getChainInfo(connect *ethclient.Client, from string) (uint64, *big.Int, *big.Int, *big.Int, error) {
	ctx := context.Background()
	nounce, err := connect.NonceAt(ctx, common.HexToAddress(from), nil)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	gasPrice, err := connect.SuggestGasPrice(ctx) ///big.NewInt(20000000000)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	tipsPrice, err := connect.SuggestGasTipCap(ctx)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	chainId, err := connect.ChainID(ctx)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	return nounce, gasPrice, tipsPrice, chainId, nil

}

func doDeployContract(from string, filepath string, pass string) string {
	var re []byte
	if from == "" || filepath == "" || pass == "" {
		fmt.Println("parameter error")
		re, _ = json.Marshal(ErrorResult{1, "parameter error"})
		return string(re)
	}
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
					re, _ = json.Marshal(ErrorResult{0, ""})
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

func doImportContract(cons string) {
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
			ctx := context.Background()
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
				jcons.Name = sresult
			}
		}
		db, err2 := leveldb.OpenFile(datadir+"contract", nil)
		defer db.Close()
		jvalue, _ := json.Marshal(jcons)
		if err2 == nil {
			if db.Put([]byte(jcons.Address), jvalue, nil) != nil {
				fmt.Println("write db error:" + string(jvalue))
			}
		} else {
			fmt.Println("open contract db error:" + err2.Error())
		}

	} else {
		fmt.Println("get wrong parameter:" + cons)
	}
}

func getConnection() *ethclient.Client {
	if G_Etherclient == nil {
		connect, err := ethclient.Dial(G_Config["nodeAddress"])
		if err != nil {
			G_Etherclient = nil

		} else {
			G_Etherclient = connect
		}
	}
	return G_Etherclient
}
