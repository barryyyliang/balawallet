package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/webview/webview"
)

var win webview.WebView
var ks *keystore.KeyStore
var dir string
var datadir string
var accountdir string
var resdir string
var exportdir string
var nodeaddress string = "http://localhost:8545"

func main() {
	win = webview.New(false)
	defer win.Destroy()
	win.SetTitle("我的钱包")
	win.SetSize(600, 480, webview.HintNone)

	//注册全局函数
	win.Bind("changePage", changePage)
	win.Bind("newaccount", newaccount)
	win.Bind("weblog", weblog)
	win.Bind("getAccounts", getAccounts)
	win.Bind("changeAccountName", changeAccountName)
	win.Bind("importAccount", importAccount)
	win.Bind("exportAccount", exportAccount)
	win.Bind("getAccountName", getAccountName)
	win.Bind("transfer", transfer)
	win.Bind("doGetBalance", doGetBalance)
	win.Bind("doGetFee", doGetFee)

	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(path, "\\", "/", -1)
	datadir = dir + "/data/"
	resdir = dir + "/res/"
	accountdir = dir + "/account/"
	exportdir = dir + "/export/"

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
