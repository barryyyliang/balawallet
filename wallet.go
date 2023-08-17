package main

import (
	"encoding/json"
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

//var nodeaddress string = "http://localhost:8545"
var G_Config = make(map[string]string)

func main() {
	win = webview.New(true)
	defer win.Destroy()
	win.SetTitle("我的钱包")
	win.SetSize(600, 480, webview.HintNone)

	//注册全局函数
	win.Bind("doChangePage", doChangePage)
	win.Bind("doNewAccount", doNewAccount)
	win.Bind("doWebLog", doWebLog)
	win.Bind("doGetAccounts", doGetAccounts)
	win.Bind("doChangeAccountName", doChangeAccountName)
	win.Bind("doImportAccount", doImportAccount)
	win.Bind("doExportAccount", doExportAccount)
	win.Bind("doGetAccountName", doGetAccountName)
	win.Bind("doTransfer", doTransfer)
	win.Bind("doGetBalance", doGetBalance)
	win.Bind("doGetFee", doGetFee)
	win.Bind("doGetComponent", doGetComponent)
	win.Bind("doConfig", doConfig)
	win.Bind("doGetConfig", doGetConfig)

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
	//初始化配置
	var conf = [1]string{"nodeAddress"}
	par, err := json.Marshal(conf)

	doGetConfig(string(par))

	content, err := os.ReadFile(dir + "/res/qianbao.html")
	if err != nil {
		fmt.Println("can not read html")
	}
	win.SetHtml(string(content))

	win.Run()
}
