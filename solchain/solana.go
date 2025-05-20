package solchain

import (
	"balawallet/accounts/keystore"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	//"github.com/ethereum/go-ethereum/crypto/sha3"
	"balawallet/accounts/sol"
	balacommon "balawallet/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	solRPC "github.com/gagliardetto/solana-go/rpc"
	"github.com/mr-tron/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

// var ks *keystore.KeyStore
var ks *keystore.KeyStore
var exportdir string
var accountdir string

func Init(dir string) {

	var accountsRoot = dir + "/account/"
	accountdir = accountsRoot + "sol/"
	exportdir = dir + "/export/"

	//创建必要的目录
	_, err := os.Stat(accountsRoot)
	if err != nil {
		if !os.IsExist(err) {
			err = os.Mkdir(accountsRoot, os.ModePerm)
			if err != nil {
				fmt.Println("create accountsRoot dir error:" + err.Error())
				return
			}
		}
	}

	_, err = os.Stat(accountdir)
	if err != nil {
		if !os.IsExist(err) {
			err = os.Mkdir(accountdir, os.ModePerm)
			if err != nil {
				fmt.Println("create eth account dir error:" + err.Error())
				return
			}
		}
	}
	_, err = os.Stat(exportdir)
	if err != nil {
		if !os.IsExist(err) {
			err = os.Mkdir(exportdir, os.ModePerm)
			if err != nil {
				fmt.Println("create export dir error:" + err.Error())
			}
		}
	}

	//打开钱包
	solAdapter := &sol.SolChainAdapter{}
	ks = keystore.NewKeyStore(accountdir, keystore.StandardScryptN, keystore.StandardScryptP, solAdapter)

}

func GetSolChain() *SolChain {
	return &SolChain{}
}

type SolChain struct {
}

// 新创建一个账号
func (c *SolChain) NewAccount(pass string, walletName string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	//scryptN := keystore.StandardScryptN
	//scryptP := keystore.StandardScryptP
	//account, err := keystore.StoreKey(accountdir, pass, scryptN, scryptP)
	account, err := ks.NewAccount(pass)
	if err != nil {
		fmt.Println("failed to generate new account")
		re = balacommon.GetResult(balacommon.Msg_createuser_fail, err.Error())
		return string(re)
	}

	//账户信息存库
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"db", nil)
	if err == nil {
		defer db.Close()
		if balacommon.IsExist(db, walletName) == false {
			err = db.Put([]byte(account.GetAddress()), []byte(walletName), nil)
			if err != nil {
				fmt.Println("put name err")
				re = balacommon.GetResult(balacommon.Msg_writename_fail, err.Error())
			}
		} else {
			fmt.Println("name has existed")
			re = balacommon.GetResult(balacommon.Msg_createuser_nameexists, "")
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_writename_fail, err.Error())
	}

	return string(re)

}

// 获取账号列表
func (c *SolChain) GetAccounts() string {
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
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"db", nil)
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
		add := acc[i].GetAddress()
		acclist.Accounts = append(acclist.Accounts, account{"sol", add, addressCache[add]})
	}

	if len(acclist.Accounts) > 0 {
		result, err2 := json.Marshal(acclist)
		if err2 == nil {
			return string(result)
		} else {
			fmt.Println(err2.Error())
		}
	}

	return ""
}

func (c *SolChain) ChangeAccountName(address string, name string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	//账户信息存库
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"db", nil)
	if err == nil {
		defer db.Close()
		if balacommon.IsExist(db, name) == false {
			err = db.Put([]byte(address), []byte(name), nil)
			if err != nil {
				re = balacommon.GetResult(balacommon.Msg_writedb_fail, err.Error())
			}
		} else {
			re = balacommon.GetResult(balacommon.Msg_changename_nameexists, "")
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_writedb_fail, err.Error())
	}

	return string(re)
}

func (c *SolChain) GetAccountName(address string) string {
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"db", nil)
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
func (c *SolChain) ImportAccount(pa string, name string, kk string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	akey, err := sol.NewKeyFromBase58(kk)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	ks.ImportKey(akey, pa)
	//账户信息存库
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"db", nil)
	if err == nil {
		defer db.Close()
		if !balacommon.IsExist(db, name) {
			err = db.Put([]byte(akey.GetStringAddress()), []byte(name), nil)
			if err != nil {
				fmt.Println("put name err")
				re = balacommon.GetResult(balacommon.Msg_writename_fail, err.Error())
			}
		} else {
			fmt.Println("name has existed")
			re = balacommon.GetResult(balacommon.Msg_createuser_nameexists, "")
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_writename_fail, err.Error())
	}

	return string(re)
}

// 导出私钥
func (c *SolChain) ExportAccount(add string, pa string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	key, err := getKey(add, pa)
	if err == nil {
		defer clear(*key)
		file, err1 := os.OpenFile(exportdir+add, os.O_RDWR|os.O_CREATE, 0666)
		if err1 != nil {
			fmt.Println("can not open file in export:" + err1.Error())
			re = balacommon.GetResult(balacommon.Msg_openexportfile_fail, err1.Error())
			return string(re)
		}
		defer file.Close()

		value := base58.Encode(*key)
		writer := bufio.NewWriter(file)
		writer.WriteString(value)
		writer.Flush()
	} else {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
	}
	return string(re)
}

func (c *SolChain) Transfer(from string, to string, value string, pass string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")
	client := balacommon.GetSolConnection()
	if client == nil {
		re = balacommon.GetResult(balacommon.Msg_fail, "connect is nil")
		return string(re)
	}
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	recent, err := client.GetLatestBlockhash(
		ctx,
		solRPC.CommitmentFinalized,
	)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	amount, err := solTolam(value)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	pubkeyFrom, err := solana.PublicKeyFromBase58(from)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}
	pubkeyTo, err := solana.PublicKeyFromBase58(to)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				pubkeyFrom,
				pubkeyTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(pubkeyFrom),
	)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	_, err = signTx(from, pass, tx)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	// Or just send the transaction WITHOUT waiting for confirmation:
	ctx1, fuc1 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc1()
	sig, err := client.SendTransaction(
		ctx1,
		tx,
	)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}
	re = balacommon.GetResult(balacommon.Msg_success, sig.String())

	return string(re)
}

// 获取账户余额
func (c *SolChain) GetBalance(add string) string {
	fmt.Println("enter sol GetBalance")
	var re []byte

	client := balacommon.GetSolConnection()
	if client != nil {
		pub, err := solana.PublicKeyFromBase58(add)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_getbalance_fail, err.Error())
			return string(re)
		}
		ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
		defer fuc()
		res, err := client.GetBalance(ctx, pub, solRPC.CommitmentConfirmed)
		if err == nil {
			tmp, _ := json.Marshal(balacommon.ResBalance{add, balacommon.ToBig(9, strconv.FormatUint(res.Value, 10))})
			re, _ = json.Marshal(balacommon.ErrorResult{0, string(tmp)})
		} else {
			re = balacommon.GetResult(balacommon.Msg_getbalance_fail, err.Error())
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_connectnode_fail, "")
	}

	return string(re)
}

func (c *SolChain) GetFee() string {
	/* 	fmt.Println("enter sol GetFee")

	   	client := balacommon.GetSolConnection()
	   	if client != nil {
	   		ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	   		defer fuc()

	   		res, err := client.GetFeeForMessage(ctx,
	   			"AQABAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAA",
	   			solRPC.CommitmentConfirmed)
	   		fmt.Println("enter sol GetFee1")
	   		if err == nil {
	   			if res.Value != nil {
	   				return strconv.FormatUint(*res.Value, 10)
	   			} else {
	   				fmt.Println("can not get fee for message")
	   			}

	   		} else {
	   			fmt.Println("get fee:" + err.Error())
	   		}
	   	}
	   	fmt.Println("enter sol GetFee2") */
	return balacommon.ToBig(9, "5000")
}

func (c *SolChain) DeployContract(from, file, pass, funcAbi, param string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) GetContract() string {

	type contracts struct {
		Contracts []balacommon.ComContract `json:"contracts"`
	}
	var res contracts
	//账户信息
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"contract", nil)
	if err == nil {
		defer db.Close()
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			tmp := strings.Split(string(it.Key()), "_")
			t, err := strconv.Atoi(tmp[0])
			if err == nil {
				if t == balacommon.Nettype_solana {
					var aa balacommon.ComContract
					json.Unmarshal(it.Value(), &aa)
					res.Contracts = append(res.Contracts, aa)
				}
			}

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

func (c *SolChain) DeleteContract(address string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	db, err := leveldb.OpenFile(balacommon.G_DataDir+"contract", nil)
	if err == nil {
		defer db.Close()
		key := strconv.Itoa(balacommon.G_NetType) + "_" + address
		if db.Delete([]byte(key), nil) != nil {
			re = balacommon.GetResult(balacommon.Msg_writedb_fail, "")
		}
	} else {
		fmt.Println("open contract db error:" + err.Error())
		re = balacommon.GetResult(balacommon.Msg_opendb_fail, err.Error())
	}
	return string(re)
}

func (c *SolChain) ImportContract(cons string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	var jcons balacommon.InputImportContract
	err := json.Unmarshal([]byte(cons), &jcons)
	if err == nil {
		connection := balacommon.GetSolConnection()
		if connection == nil {
			re = balacommon.GetResult(balacommon.Msg_fail, "connect is nil")
			return string(re)
		}
		pubkey, err := solana.PublicKeyFromBase58(jcons.Address)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			return string(re)
		}
		ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
		defer fuc()
		out, err := connection.GetAccountInfo(ctx, pubkey)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			fmt.Println("get mint account:" + err.Error())
			return string(re)
		}
		min := new(token.Mint)
		data := out.Value.Data.GetBinary()
		decoder := bin.NewBinDecoder(data)
		err = decoder.Decode(min)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			return string(re)
		} else {
			db, err2 := leveldb.OpenFile(balacommon.G_DataDir+"contract", nil)
			co := balacommon.ComContract{Type: jcons.Type, Address: jcons.Address, Name: jcons.Name}
			jvalue, _ := json.Marshal(co)
			if err2 == nil {
				defer db.Close()
				key := strconv.Itoa(balacommon.Nettype_solana) + "_" + jcons.Address
				if db.Put([]byte(key), jvalue, nil) != nil {
					re = balacommon.GetResult(balacommon.Msg_writedb_fail, "")
				}
			} else {
				fmt.Println("open contract db error:" + err2.Error())
				re = balacommon.GetResult(balacommon.Msg_opendb_fail, err2.Error())
			}
		}

	}

	return string(re)
}

func (c *SolChain) CheckNode() string {
	var re []byte
	if balacommon.G_Config["solAddress"] != "" {
		re = balacommon.GetResult(balacommon.Msg_success, "")
	} else {
		re = balacommon.GetResult(balacommon.Msg_nodeconfig_null, "")
	}
	return string(re)
}

func (c *SolChain) LoadABI(file string) string {
	var re []byte
	if file == "SPL" {
		abiContent := balacommon.G_Abis[file]
		if abiContent != "" {
			re, _ = json.Marshal(balacommon.ErrorResult{0, string(abiContent)})
			return string(re)
		}
	}

	filepath := balacommon.G_Config["contractDir"]
	filepath = filepath + "/" + file
	if filepath == "" {
		re = balacommon.GetResult(balacommon.Msg_param_error, "")
		return string(re)
	}

	abiContent, err := os.ReadFile(filepath)
	if err != nil {
		re = balacommon.GetResult(balacommon.Msg_openfile_fail, err.Error())
		return string(re)
	}

	re, _ = json.Marshal(balacommon.ErrorResult{0, string(abiContent)})
	return string(re)
}

func (c *SolChain) ReadContract(from, conAdd, abiFile, funcName, paremters string) string {
	var re = balacommon.GetResult(balacommon.Msg_fail, "")
	sabi := c.LoadABI(abiFile)
	var abiResult = new(balacommon.ErrorResult)
	err := json.Unmarshal([]byte(sabi), abiResult)
	var contractAbi abi.ABI
	if err == nil && abiResult.Code == balacommon.Msg_success {
		contractAbi, err = abi.JSON(strings.NewReader(abiResult.Msg))
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_parseabi_fail, err.Error())
			return string(re)
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_parseabi_fail, "")
	}

	for _, method := range contractAbi.Methods {
		if method.RawName != funcName {
			continue
		}
		var parStrings []string
		json.Unmarshal([]byte(paremters), &parStrings)
		if len(parStrings) != len(method.Inputs) {
			re = balacommon.GetResult(balacommon.Msg_param_error, "")
		}
		tokenAdd, err := solana.PublicKeyFromBase58(conAdd)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			return string(re)
		}
		client := balacommon.GetSolConnection()
		if client == nil {
			re = balacommon.GetResult(balacommon.Msg_fail, "connect is nil")
			return string(re)
		}
		result, err := ReadTokenContract(tokenAdd, method, parStrings, client)
		if err == nil {
			re, _ = json.Marshal(balacommon.ErrorResult{balacommon.Msg_success, result})
		} else {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		}
	}

	return string(re)
}

func (c *SolChain) InvokeContract(from, conAdd, abiFile, funcName, param, pass string) string {
	var re = balacommon.GetResult(balacommon.Msg_fail, "")
	sabi := c.LoadABI(abiFile)
	var abiResult = new(balacommon.ErrorResult)
	err := json.Unmarshal([]byte(sabi), abiResult)
	var contractAbi abi.ABI
	if err == nil && abiResult.Code == balacommon.Msg_success {
		contractAbi, err = abi.JSON(strings.NewReader(abiResult.Msg))
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_parseabi_fail, err.Error())
			return string(re)
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_parseabi_fail, "")
	}

	for _, method := range contractAbi.Methods {
		if method.RawName != funcName {
			continue
		}
		var parStrings []string
		json.Unmarshal([]byte(param), &parStrings)
		if len(parStrings) != len(method.Inputs) {
			re = balacommon.GetResult(balacommon.Msg_param_error, "")
		}
		tokenAdd, err := solana.PublicKeyFromBase58(conAdd)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			return string(re)
		}
		payer, err := solana.PublicKeyFromBase58(from)
		if err != nil {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
			return string(re)
		}
		client := balacommon.GetSolConnection()
		if client == nil {
			re = balacommon.GetResult(balacommon.Msg_fail, "connect is nil")
			return string(re)
		}
		result, err := InvokeTokenContract(tokenAdd, method, parStrings, payer, pass, client)
		if err == nil {
			re, _ = json.Marshal(balacommon.ErrorResult{balacommon.Msg_success, result})
		} else {
			re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
		}
	}

	return string(re)
}

func (c *SolChain) InvokeEstimate(from, conAddr, funcAbi, param string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	return string(re)

}

func (c *SolChain) QueryTx(hs string) string {
	var re []byte
	client := balacommon.GetSolConnection()
	if client != nil {
		txSig := solana.MustSignatureFromBase58(hs)
		ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
		defer fuc()

		opts := new(solRPC.GetTransactionOpts)
		opts.MaxSupportedTransactionVersion = new(uint64)
		out, err := client.GetTransaction(
			ctx,
			txSig,
			opts,
		)
		if err == nil {
			re, err = json.Marshal(out)
			if err != nil {
				fmt.Println("so get transaction marshal result error :" + out.BlockTime.String())
			}
		} else {
			fmt.Println("get sol transaction err:" + err.Error())
		}
	}
	return string(re)
}

func (c *SolChain) Signature(info string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) SignatureVerify(info string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")
	return string(re)
}

func (c SolChain) GetAbiFile() string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	var conf = [1]string{"contractDir"}
	par, _ := json.Marshal(conf)
	balacommon.GetConfig(string(par))
	var dir = balacommon.G_Config["contractDir"]
	if dir != "" {
		names, err := balacommon.GetFilesName(dir, ".abi")
		if err == nil {
			tmp, _ := json.Marshal(names)
			re, _ = json.Marshal(balacommon.ErrorResult{0, string(tmp)})
		} else {
			re = balacommon.GetResult(balacommon.Msg_openfile_fail, err.Error())
		}
	} else {
		re = balacommon.GetResult(balacommon.Msg_contractdirconfig_null, "")
	}
	return string(re)
}

func (c *SolChain) GetBinFile() string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) ChangeAccountPass(atype, address, old, new string) string {
	re := balacommon.GetResult(balacommon.Msg_success, "")
	acc := ks.Accounts()
	var i = 0
	for i = 0; i < len(acc); i++ {
		if acc[i].GetAddress() == address {
			err2 := ks.Update(acc[i], old, new)
			if err2 != nil {
				re = balacommon.GetResult(balacommon.Msg_fail, err2.Error())
			}
			break
		}
	}
	if i == len(acc) {
		re = balacommon.GetResult(balacommon.Msg_account_noexist, "")
	}
	return string(re)
}

func (c *SolChain) GetContractBalance(assetType, conAddr, ownerAddr string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	abis := c.LoadABI(assetType)
	var abiResult = new(balacommon.ErrorResult)
	err := json.Unmarshal([]byte(abis), abiResult)
	if err == nil && abiResult.Code == balacommon.Msg_success {
		contractAbi, err := abi.JSON(strings.NewReader(abiResult.Msg))
		if err != nil {
			return string(re)
		}
		for _, value := range contractAbi.Methods {
			if value.Name == "balanceOf" &&
				len(value.Inputs) == 1 &&
				value.Inputs[0].Type.T == abi.AddressTy {
				if assetType == "SPL" {
					client := balacommon.GetSolConnection()
					if client == nil {
						re = balacommon.GetResult(balacommon.Msg_fail, "connect is nil")
						return string(re)
					}
					mint, err := solana.PublicKeyFromBase58(conAddr)
					if err != nil {
						re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
						return string(re)
					}

					result, err := ReadTokenContract(mint, value, []string{ownerAddr}, client)
					if err == nil {
						var arrRet []string
						err = json.Unmarshal([]byte(result), &arrRet)
						if err != nil {
							re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
							return string(re)
						}
						tmp, _ := json.Marshal(balacommon.ResBalance{ownerAddr, arrRet[0]})
						re, _ = json.Marshal(balacommon.ErrorResult{balacommon.Msg_success, string(tmp)})
						return string(re)
					} else {
						re = balacommon.GetResult(balacommon.Msg_fail, err.Error())
					}
				}
				break
			}
		}
	}

	return string(re)
}
func (c *SolChain) DeployEstimate(from, conFile, funcAbi, param string) string {
	re := balacommon.GetResult(balacommon.Msg_fail, "not support")

	return string(re)
}
