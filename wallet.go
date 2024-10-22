package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	webview "github.com/webview/webview_go"
)

var win webview.WebView
var ks *keystore.KeyStore
var dir string
var datadir string
var accountdir string
var resdir string
var exportdir string
var currentAccount string

var G_Config = make(map[string]string)
var G_Etherclient *ethclient.Client = nil
var G_BlockchainConfig bool = true

func main() {
	InitHtmlFile()
	initMsg()
	initRes()
	win = webview.New(true)
	defer win.Destroy()
	win.SetTitle("Web3Wallet")
	win.SetSize(1024, 640, webview.HintNone)

	//注册全局函数
	//win.Bind("doChangePage", doChangePage)
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
	win.Bind("doGetContract", doGetContract)
	win.Bind("doImportContract", doImportContract)
	win.Bind("doDeployContract", doDeployContract)
	win.Bind("doCheckNode", doCheckNode)
	win.Bind("doLoadABI", doLoadABI)
	win.Bind("doReadContract", doReadContract)
	win.Bind("doQueryTx", doQueryTx)
	win.Bind("doSignature", doSignature)
	win.Bind("doGetAbiFile", doGetAbiFile)
	win.Bind("doGetBinFile", doGetBinFile)
	win.Bind("doInvokeContract", doInvokeContract)
	win.Bind("doGetRes", doGetRes)
	win.Bind("doGetMsg", doGetMsg)

	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.Replace(path, "\\", "/", -1)
	datadir = dir + "/data/"
	//resdir = dir + "/res/"
	accountdir = dir + "/account/"
	exportdir = dir + "/export/"

	var err error
	//创建必要的目录
	_, err = os.Stat(datadir)
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(datadir, os.ModePerm)
		}
	}
	_, err = os.Stat(accountdir)
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(accountdir, os.ModePerm)
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
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP
	ks = keystore.NewKeyStore(dir+"/account/", scryptN, scryptP)
	//初始化配置
	var conf = [1]string{"nodeAddress"}
	par, _ := json.Marshal(conf)

	doGetConfig(string(par))
	if G_Config["nodeAddress"] != "" {
		G_Etherclient, err = ethclient.Dial(G_Config["nodeAddress"])
		if err != nil {
			G_Etherclient = nil
			fmt.Println("can not connect to node")
		}
	} else {
		fmt.Println("have not config blockchain node")
		G_BlockchainConfig = false
	}

	content := G_Files["/main.html"]
	if len(content) == 0 {
		fmt.Println("can not read html")
	}
	win.SetHtml(string(content))
	//same website can not be opened by this way
	//win.Navigate("https://github.com/webview/webview")

	win.Run()
}
