package solchain

import (
	. "balawallet/common"
	"fmt"
	"os"
	//"github.com/ethereum/go-ethereum/crypto/sha3"
)

//var ks *keystore.KeyStore

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
	//ks = keystore.NewKeyStore(accountdir, keystore.StandardScryptN, keystore.StandardScryptP)

}

func GetSolChain() *SolChain {
	return &SolChain{}
}

type SolChain struct {
}

// 新创建一个账号
func (c *SolChain) NewAccount(pass string, walletName string) string {
	re := GetResult(Msg_success, "")

	return string(re)

}

// 获取账号列表
func (c *SolChain) GetAccounts() string {

	return ""
}

func (c *SolChain) ChangeAccountName(address string, name string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) GetAccountName(address string) string {

	return ""

}

// 通过私钥导入账号
func (c *SolChain) ImportAccount(pa string, name string, kk string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}

// 导出私钥
func (c *SolChain) ExportAccount(add string, pa string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) Transfer(from string, to string, value string, pass string) string {
	var re []byte

	return string(re)
}

// 获取账户余额
func (c *SolChain) GetBalance(add string) string {
	var re []byte
	return string(re)
}

func (c *SolChain) GetFee() string {
	return ""
}

func (c *SolChain) DeployContract(from string, file string, pass string) string {
	var re []byte

	return string(re)
}

func (c *SolChain) GetContract() string {

	return ""
}

func (c *SolChain) DeleteContract(address string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) ImportContract(cons string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) CheckNode() string {
	var re []byte

	return string(re)
}

func (c *SolChain) LoadABI(file string) string {
	var re []byte

	return string(re)
}

func (c *SolChain) ReadContract(from, to, funcAbi, paremters string) string {

	return ""
}

func (c *SolChain) InvokeContract(from, conAddr, funcAbi, param, pass string) string {
	var re []byte

	return string(re)
}

func (c *SolChain) InvokeEstimate(from, conAddr, funcAbi, param string) string {
	var re []byte

	return string(re)

}

func (c *SolChain) QueryTx(hs string) string {
	var re []byte

	return string(re)
}

func (c *SolChain) Signature(info string) string {
	var re []byte

	return string(re)
}

func (c *SolChain) SignatureVerify(info string) string {
	var re []byte
	return string(re)
}

func (c SolChain) GetAbiFile() string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) GetBinFile() string {
	re := GetResult(Msg_success, "")

	return string(re)
}

func (c *SolChain) ChangeAccountPass(atype, address, old, new string) string {
	re := GetResult(Msg_success, "")

	return string(re)
}
