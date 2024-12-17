package common

import (
	"encoding/json"
	"fmt"
	"strconv"

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

	result.Css = string(G_Files[cssfile])
	result.Js = string(G_Files[jsfile])
	result.Html = string(G_Files[htmlfile])
	language := G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := G_Res[language][htmlfile]
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
	re := GetResult(Msg_success, "")
	var config = make(map[string]string)

	json.Unmarshal([]byte(configs), &config)

	for key, value := range config {
		G_Config[key] = value
	}

	db, _ := leveldb.OpenFile(G_DataDir+"db", nil)
	defer db.Close()
	tmp, _ := json.Marshal(G_Config)
	db.Put([]byte("config"), tmp, nil)

	/*
		for key, value := range G_Config {
			err = db.Put([]byte(key), []byte(value), nil)
			if err != nil {
				fmt.Println("put config error")
				re = getResult(Msg_writedb_fail, err.Error())
			} else {
				fmt.Println("put config:" + key + ":" + value)
			}
		}
	*/
	// 配置更新后刷新链接
	G_Etherclient = nil
	return string(re)
}

func DoGetConfig(configs string) string {
	var re string
	db, err := leveldb.OpenFile(G_DataDir+"db", nil)
	if err != nil {
		fmt.Println("open db error:" + err.Error())
		return ""
	}

	defer db.Close()

	data, err := db.Get([]byte("config"), nil)
	if err != nil {
		return ""
	}
	err = json.Unmarshal(data, &G_Config)
	if err != nil {
		return ""
	}

	if configs != "" {
		var configName []string
		json.Unmarshal([]byte(configs), &configName)
		var result = make(map[string]string)
		for _, value := range configName {
			result[value] = G_Config[value]
		}
		tmp, _ := json.Marshal(result)
		re = string(tmp)
	} else {
		tmp, _ := json.Marshal(G_Config)
		re = string(tmp)
	}
	return string(re)
}

func DoGetMsg(code string) string {
	n, _ := strconv.Atoi(code)
	lang := G_Config["language"]
	if lang == "" {
		lang = "zh-CN"
	}
	msg := G_Msg[lang][uint(n)]
	if msg == "" {
		msg = "null"
	}

	return msg
}

func DoGetRes(com string) string {
	var re string = ""
	htmlfile := "/" + com + ".html"
	language := G_Config["language"]
	if language == "" {
		language = "zh-CN"
	}
	res := G_Res[language][htmlfile]
	if len(res) != 0 {
		tmp, _ := json.Marshal(res)
		re = string(tmp)
	}
	return re
}

func DoChangeNet(nettype string) string {
	var re []byte
	tmpType := -1
	switch nettype {
	case "Ethereum":
		tmpType = Nettype_ethereum
	case "Polygon":
		tmpType = Nettype_polygon
	case "Solana":
		tmpType = Nettype_solana
	}
	if G_NetType != tmpType {
		if tmpType == Nettype_solana {
			re = GetResult(Msg_fail, "not support")
		} else {
			G_NetType = tmpType
			G_Etherclient = nil
			re = GetResult(Msg_success, "")
		}
	} else {
		re = GetResult(Msg_success, "")
	}
	return string(re)
}
