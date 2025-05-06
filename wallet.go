package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"balawallet/common"
	"balawallet/ethchain"
	"balawallet/pagefunc"
	"balawallet/serverapi"
	"balawallet/solchain"

	webview "github.com/webview/webview_go"
)

var win webview.WebView

var dir string
var resdir string
var currentAccount string

func main() {
	common.InitHtmlFile() //generate by tool
	common.InitAbiFile()  //generate by tool
	common.InitMsg()
	common.InitRes()
	var local bool = false
	if local {
		//for local
		path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		dir = strings.Replace(path, "\\", "/", -1)

	} else {
		home, _ := user.Current()
		dir = home.HomeDir + "/.balawallet"
	}

	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}

	common.InitGlobal(dir)
	ethchain.Init(dir)
	solchain.Init(dir)
	tmp := ethchain.GetEthChain()
	server := serverapi.GetServer(tmp)

	win = webview.New(false)
	defer win.Destroy()
	win.SetTitle("BalaWallet")
	win.SetSize(1280, 800, webview.HintNone)

	//注册全局函数

	win.Bind("doWebLog", pagefunc.DoWebLog)
	win.Bind("doConfig", pagefunc.DoConfig)
	win.Bind("doGetConfig", pagefunc.DoGetConfig)
	win.Bind("doGetRes", pagefunc.DoGetRes)
	win.Bind("doGetMsg", pagefunc.DoGetMsg)
	win.Bind("doGetComponent", pagefunc.DoGetComponent)
	win.Bind("doChangeNet", pagefunc.DoChangeNet)
	win.Bind("doSaveDataSpace", common.DoSaveDataSpace)
	win.Bind("doChangeDataSpace", common.DoChangeDataSpace)
	win.Bind("doDeleteDataSpace", common.DoDeleteDataSpace)
	win.Bind("doGetDataSpace", common.DoGetDataSpace)
	win.Bind("doCheckDataSpace", common.DoCheckDataSpace)
	win.Bind("doInitDataSpace", common.DoInitDataSpace)
	win.Bind("doGetDataItem", common.DoGetDataItem)
	win.Bind("doChangeSpacePass", common.DoChangeSpacePass)

	win.Bind("doNewAccount", server.DoNewAccount)
	win.Bind("doGetAccounts", server.DoGetAccounts)
	win.Bind("doChangeAccountName", server.DoChangeAccountName)
	win.Bind("doImportAccount", server.DoImportAccount)
	win.Bind("doExportAccount", server.DoExportAccount)
	win.Bind("doGetAccountName", server.DoGetAccountName)
	win.Bind("doTransfer", server.DoTransfer)
	win.Bind("doGetBalance", server.DoGetBalance)
	win.Bind("doGetFee", server.DoGetFee)
	win.Bind("doGetContract", server.DoGetContract)
	win.Bind("doImportContract", server.DoImportContract)
	win.Bind("doDeployContract", server.DoDeployContract)
	win.Bind("doCheckNode", server.DoCheckNode)
	win.Bind("doLoadABI", server.DoLoadABI)
	win.Bind("doReadContract", server.DoReadContract)
	win.Bind("doQueryTx", server.DoQueryTx)
	win.Bind("doSignature", server.DoSignature)
	win.Bind("doGetAbiFile", server.DoGetAbiFile)
	win.Bind("doGetBinFile", server.DoGetBinFile)
	win.Bind("doInvokeContract", server.DoInvokeContract)
	win.Bind("doDeleteContract", server.DoDeleteContract)
	win.Bind("doSignatureVerify", server.DoSignatureVerify)
	win.Bind("doChangeAccountPass", server.DoChangeAccountPass)
	win.Bind("doInvokeEstimate", server.DoInvokeEstimate)
	win.Bind("doGetContractBalance", server.DoGetContractBalance)
	win.Bind("doDeployEstimate", server.DoDeployEstimate)

	content := common.G_Files["/main.html"]
	if len(content) == 0 {
		fmt.Println("can not read html")
	}
	win.SetHtml(string(content))
	win.Run()
}
