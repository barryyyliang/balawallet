package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webview/webview"
)

var win webview.WebView
var dir string
var ks *keystore.KeyStore

//页面切换
func changePage(page string, param string) {
	fmt.Println(dir + "/res/" + page)
	content, err := os.ReadFile(dir + "/res/" + page)
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
func newaccount(pass string, walletName string) {
	fmt.Println("password:" + pass)
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	account, err := keystore.StoreKey(dir+"/account/", pass, scryptN, scryptP)
	if err != nil {
		fmt.Println("failed to generate new account")
	}
	fmt.Println(account.Address)
	fmt.Println(account.URL)
	fmt.Println(account.URL.Path)
	fmt.Println(account.URL.Scheme)

	//账户信息存库
	db, err := leveldb.OpenFile(dir+"/data/db", nil)
	defer db.Close()

	if isAddressExist(db, account.Address.String(), walletName) == false {
		err = db.Put([]byte(account.Address.String()), []byte(walletName), nil)
		if err != nil {
			fmt.Println("put name err")
		}
	} else {
		fmt.Println("name has existed")
	}

}

//日志接口
func weblog(para string) {
	fmt.Println(para)
}

//获取账号列表
func getAccounts() string {
	//账户信息
	addressCache := make(map[string]string)
	db, err := leveldb.OpenFile(dir+"/data/db", nil)
	defer db.Close()
	if err == nil {
		it := db.NewIterator(nil, nil)
		defer it.Release()
		for it.Next() {
			//fmt.Println("get address:" + string(it.Key()) + ":" + string(it.Value()))
			fmt.Printf("get address: %s :\r\n", string(it.Key()))
			addressCache[string(it.Key())] = string(it.Value())
			fmt.Println("map:" + string(it.Key()) + ":" + addressCache[string(it.Key())])
		}
	}

	acc := ks.Accounts()
	var acclist string
	for i := 0; i < len(acc); i++ {
		if 0 == i {
			//fmt.Println("get name:" + acc[i].Address.String())
			fmt.Printf("get name: %s : \r\n", acc[i].Address.String())
			add := acc[i].Address.String()
			acclist = add + "-" + addressCache[add]
		} else {
			//fmt.Println("get name:" + acc[i].Address.String())
			fmt.Printf("get name: %s : \r\n", acc[i].Address.String())
			add := acc[i].Address.String()
			acclist = acclist + "," + add + "-" + addressCache[add]
		}
	}
	return acclist
}

func changeWalletName(address string, name string) {
	//账户信息存库
	db, err := leveldb.OpenFile(dir+"/data/db", nil)
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

//通过私钥导入账号
func importAccount(pa string, name string, kk string) {
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
func exportAccount(add string, pa string) {
	fmt.Println("exportAccount")
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

		file, err := os.OpenFile(dir+"/export/"+add, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("can not open file in export")
			return
		}
		defer file.Close()

		value := hex.EncodeToString(key.PrivateKey.D.Bytes())
		writer := bufio.NewWriter(file)
		writer.WriteString(value)
		writer.Flush()
	}
}

func main() {
	win = webview.New(true)
	defer win.Destroy()
	win.SetTitle("Basic Example")
	win.SetSize(600, 480, webview.HintNone)

	//注册全局函数
	win.Bind("changePage", changePage)
	win.Bind("newaccount", newaccount)
	win.Bind("weblog", weblog)
	win.Bind("getAccounts", getAccounts)
	win.Bind("changeWalletName", changeWalletName)
	win.Bind("importAccount", importAccount)
	win.Bind("exportAccount", exportAccount)

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(path, "\\", "/", -1)

	//打开钱包
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks = keystore.NewKeyStore(dir+"/account/", scryptN, scryptP)

	content, err := os.ReadFile(dir + "/res/qianbao.html")
	if err != nil {
		fmt.Println("can not read html")
	}
	win.SetHtml(string(content))

	win.Run()
}
