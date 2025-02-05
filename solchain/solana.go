package solchain

import (
	"balawallet/accounts/keystore"
	"balawallet/common"
	"encoding/json"
	"fmt"
	"os"

	//"github.com/ethereum/go-ethereum/crypto/sha3"
	"balawallet/accounts/sol"

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
	re := common.GetResult(common.Msg_success, "")
	//scryptN := keystore.StandardScryptN
	//scryptP := keystore.StandardScryptP
	//account, err := keystore.StoreKey(accountdir, pass, scryptN, scryptP)
	account, err := ks.NewAccount(pass)
	if err != nil {
		fmt.Println("failed to generate new account")
		re = common.GetResult(common.Msg_createuser_fail, err.Error())
		return string(re)
	}

	//账户信息存库
	db, err := leveldb.OpenFile(common.G_DataDir+"db", nil)
	if err == nil {
		defer db.Close()
		if common.IsExist(db, walletName) == false {
			err = db.Put([]byte(account.GetAddress()), []byte(walletName), nil)
			if err != nil {
				fmt.Println("put name err")
				re = common.GetResult(common.Msg_writename_fail, err.Error())
			}
		} else {
			fmt.Println("name has existed")
			re = common.GetResult(common.Msg_createuser_nameexists, "")
		}
	} else {
		re = common.GetResult(common.Msg_writename_fail, err.Error())
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
	db, err := leveldb.OpenFile(common.G_DataDir+"db", nil)
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
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) GetAccountName(address string) string {

	return ""

}

// 通过私钥导入账号
func (c *SolChain) ImportAccount(pa string, name string, kk string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

// 导出私钥
func (c *SolChain) ExportAccount(add string, pa string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) Transfer(from string, to string, value string, pass string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

// 获取账户余额
func (c *SolChain) GetBalance(add string) string {
	re := common.GetResult(common.Msg_fail, "not support")
	return string(re)
}

func (c *SolChain) GetFee() string {
	return ""
}

func (c *SolChain) DeployContract(from string, file string, pass string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) GetContract() string {

	return ""
}

func (c *SolChain) DeleteContract(address string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) ImportContract(cons string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) CheckNode() string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) LoadABI(file string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) ReadContract(from, to, funcAbi, paremters string) string {
	re := common.GetResult(common.Msg_fail, "not support")
	return string(re)
}

func (c *SolChain) InvokeContract(from, conAddr, funcAbi, param, pass string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) InvokeEstimate(from, conAddr, funcAbi, param string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)

}

func (c *SolChain) QueryTx(hs string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) Signature(info string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) SignatureVerify(info string) string {
	re := common.GetResult(common.Msg_fail, "not support")
	return string(re)
}

func (c SolChain) GetAbiFile() string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) GetBinFile() string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}

func (c *SolChain) ChangeAccountPass(atype, address, old, new string) string {
	re := common.GetResult(common.Msg_fail, "not support")

	return string(re)
}
