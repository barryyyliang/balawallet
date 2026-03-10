package pagefunc

import (
	"balawallet/common"
	"balawallet/ethchain"
	"balawallet/serverapi"
	"balawallet/solchain"
	"encoding/json"
	"fmt"

	//"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/syndtr/goleveldb/leveldb"
)

// //页面切换
// func doChangePage(page string, param string) {
// 	fmt.Println(resdir + page)
// 	content, err := os.ReadFile(resdir + page)
// 	if err != nil {
// 		fmt.Println("can not read html")
// 	}
// 	win.SetHtml(string(content))
// }

// 日志接口
func DoWebLog(para string) {
	fmt.Println(para)
}

// func doGetComponent(com string) string {
// 	fmt.Println("enter get" + com)

// 	if com == "rightmenu" {
// 		type menujson struct {
// 			Menumcss string `json:"css"`
// 			Menuhtml string `json:"html"`
// 			Menujs   string `json:"js"`
// 		}
// 		menucss, err1 := os.ReadFile(dir + "/res/common/rightmenu.css")
// 		if err1 != nil {
// 			fmt.Println("doGetComponent get rightmenu.css error")
// 			return ""
// 		}
// 		menu, err2 := os.ReadFile(dir + "/res/common/rightmenu.html")
// 		if err2 != nil {
// 			fmt.Println("doGetComponent get rightmenu.html error")
// 			return ""
// 		}

// 		menujs, err3 := os.ReadFile(dir + "/res/common/rightmenu.js")
// 		if err3 != nil {
// 			fmt.Println("doGetComponent get rightmenu.js error")
// 			return ""
// 		}
// 		result := menujson{string(menucss), string(menu), string(menujs)}
// 		re, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println("doGetComponent json error")
// 			return ""
// 		}
// 		return string(re)
// 	} else {
// 		type pagejson struct {
// 			Pagecss  string `json:"css"`
// 			Pagehtml string `json:"html"`
// 			Pagejs   string `json:"js"`
// 		}
// 		Pagecss, err1 := os.ReadFile(dir + "/res/" + com + ".css")
// 		if err1 != nil {
// 			fmt.Println("doGetComponent get " + com + ".css error")
// 		}
// 		Pagejs, err3 := os.ReadFile(dir + "/res/" + com + ".js")
// 		if err3 != nil {
// 			fmt.Println("doGetComponent get " + com + ".js error")
// 		}
// 		Pagehtml, err2 := os.ReadFile(dir + "/res/" + com + ".html")
// 		if err2 != nil {
// 			fmt.Println("doGetComponent get " + com + ".html error")
// 			return ""
// 		}

// 		result := pagejson{string(Pagecss), string(Pagehtml), string(Pagejs)}
// 		re, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println("doGetComponent json error")
// 			return ""
// 		}
// 		return string(re)
// 	}
// }

func DoGetComponent(com string) string {
	type pagejson struct {
		Css  string `json:"css"`
		Html string `json:"html"`
		Js   string `json:"js"`
		Res  string `json:"res"`
	}

	result := pagejson{"", "", "", ""}
	var cssfile, htmlfile, jsfile string

	if com == "rightmenu" {
		cssfile = "/common/rightmenu.css"
		htmlfile = "/common/rightmenu.html"
		jsfile = "/common/rightmenu.js"
	} else {
		cssfile = "/" + com + ".css"
		jsfile = "/" + com + ".js"
		htmlfile = "/" + com + ".html"
	}

	result.Css = string(common.G_Files[cssfile])
	result.Js = string(common.G_Files[jsfile])
	result.Html = string(common.G_Files[htmlfile])
	language := common.G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := common.G_Res[language][htmlfile]
	if len(res) != 0 {
		tmp, _ := json.Marshal(res)
		result.Res = string(tmp)
	}

	re, err := json.Marshal(result)
	if err != nil {
		fmt.Println("doGetComponent json error")
		return ""
	}
	return string(re)
}

func DoConfig(configs string) string {
	re := common.GetResult(common.Msg_success, "")

	// 首先解析整个配置为map[string]interface{}
	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(configs), &configMap)
	if err != nil {
		fmt.Println("解析配置失败:", err)
		return string(re)
	}

	// 处理字符串类型的配置
	for key, value := range configMap {
		if strValue, ok := value.(string); ok {
			common.G_Config[key] = strValue
		} else if key == "customNetworks" {
			// 处理customNetworks字段，将其转换为JSON字符串后保存
			customNetworksJSON, err := json.Marshal(value)
			if err == nil {
				common.G_Config[key] = string(customNetworksJSON)
			}
		}
	}

	db, _ := leveldb.OpenFile(common.G_DataDir+"db", nil)
	defer db.Close()
	tmp, _ := json.Marshal(common.G_Config)
	db.Put([]byte("config"), tmp, nil)

	// 配置更新后刷新链接
	common.NullConnect(common.Nettype_ethereum)
	common.NullConnect(common.Nettype_polygon)
	common.NullConnect(common.Nettype_solana)
	common.NullConnect(common.Nettype_pharos)
	return string(re)
}

func DoGetConfig(configs string) string {
	return common.GetConfig(configs)
}

func DoGetMsg(code string) string {
	return common.GetMsg(code)
}

func DoGetRes(com string) string {
	var re string = ""
	htmlfile := "/" + com + ".html"
	language := common.G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := common.G_Res[language][htmlfile]
	if len(res) != 0 {
		tmp, _ := json.Marshal(res)
		re = string(tmp)
	}
	return re
}

func getNetType(nettype string) string {
	tmpType := "-1"
	switch nettype {
	case "Ethereum":
		tmpType = common.Nettype_ethereum
	case "Polygon":
		tmpType = common.Nettype_polygon
	case "Solana":
		tmpType = common.Nettype_solana
	case "Pharos":
		tmpType = common.Nettype_pharos
	}
	if tmpType == "-1" {
		tmpType = nettype
	}

	return tmpType
}

func DoChangeNet(nettype string) string {
	var re []byte
	tmpType := getNetType(nettype)

	if common.G_NetType != tmpType {
		if tmpType == common.Nettype_solana {
			//re = GetResult(Msg_fail, "not support")
			re = common.GetResult(common.Msg_success, "")
			common.G_NetType = tmpType
			tmp := solchain.GetSolChain()
			serverapi.SetServer(tmp)
		} else {
			common.G_NetType = tmpType
			common.NullConnect(common.Nettype_ethereum)
			tmp := ethchain.GetEthChain()
			serverapi.SetServer(tmp)
			re = common.GetResult(common.Msg_success, "")
		}
	} else {
		re = common.GetResult(common.Msg_success, "")
	}
	return string(re)
}
