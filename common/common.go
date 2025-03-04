package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	solRPC "github.com/gagliardetto/solana-go/rpc"
	"github.com/syndtr/goleveldb/leveldb"
)

var G_Timeout time.Duration = 5 * time.Second
var G_Config = make(map[string]string)
var g_Etherclient *ethclient.Client = nil
var g_SolClient *solRPC.Client = nil
var G_BlockchainConfig bool = true
var G_DataDir string
var G_NetType int

func InitGlobal(dir string) {
	G_NetType = Nettype_ethereum
	G_DataDir = dir + "/data/"
	_, err := os.Stat(G_DataDir)
	if err != nil {
		if !os.IsExist(err) {
			os.Mkdir(G_DataDir, os.ModePerm)
		}
	}
	//初始化配置
	GetConfig("")
	if G_Config["ethAddress"] != "" {
		g_Etherclient, err = ethclient.Dial(G_Config["ethAddress"])
		if err != nil {
			g_Etherclient = nil
			fmt.Println("can not connect to ethereum node")
		}
	} else {
		fmt.Println("have not config ethereum node")
		G_BlockchainConfig = false
	}

	if G_Config["solAddress"] != "" {
		g_SolClient = solRPC.New(G_Config["solAddress"])
		if g_SolClient == nil {
			fmt.Println("can not connect to ethereum node")
		}
	} else {
		fmt.Println("have not config solana node")
	}
}

func IsExist(db *leveldb.DB, name string) bool {
	it := db.NewIterator(nil, nil)
	defer it.Release()
	for it.Next() {
		if string(it.Value()) == name {
			return true
		}
	}
	return false
}

func GetResult(code uint, addMsg string) []byte {
	var re []byte
	lang := G_Config["language"]
	if lang == "" {
		lang = "zh-CN"
	}
	msg := G_Msg[lang][code]
	if msg == "" {
		msg = "null"
	}
	if addMsg != "" {
		msg = msg + ":" + addMsg
	}
	re, _ = json.Marshal(ErrorResult{code, msg})
	return re
}

func GetMsgByCode(code uint) string {
	lang := G_Config["language"]
	if lang == "" {
		lang = "zh-CN"
	}
	msg := G_Msg[lang][code]
	return msg
}

func GetMsg(code string) string {
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

func GetConfig(configs string) string {
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

func GetSolConnection() *solRPC.Client {
	if g_SolClient == nil {
		if G_Config["solAddress"] != "" {
			g_SolClient = solRPC.New(G_Config["solAddress"])
			if g_SolClient == nil {
				fmt.Println("can not connect to ethereum node")
			}
		} else {
			fmt.Println("have not config solana node")
		}
	}
	return g_SolClient
}

func GetEthConnection() *ethclient.Client {
	if g_Etherclient == nil {
		var connect *ethclient.Client
		var err error
		if Nettype_ethereum == G_NetType {
			connect, err = ethclient.Dial(G_Config["ethAddress"])
		} else if Nettype_polygon == G_NetType {
			connect, err = ethclient.Dial(G_Config["polAddress"])
		}

		if err != nil {
			fmt.Println("net connet fail")
			g_Etherclient = nil

		} else {
			g_Etherclient = connect
		}
	}
	return g_Etherclient
}

func NullConnect(ntye int) {
	switch ntye {
	case Nettype_ethereum:
		g_Etherclient = nil
	case Nettype_polygon:
		g_Etherclient = nil
	case Nettype_solana:
		g_SolClient = nil
	}
}

/*
将一个小的单位的数转换为一个大的单位的数，
unittag 转换率，比如 “个” 转换为 “万” ，1万=10^4,所以参数为 4
*/

func ToBig(unittag int, wei string) string {
	ballength := len(wei)
	var balance string
	if ballength > unittag {
		var first []byte = make([]byte, ballength-unittag)
		copy(first, []byte(wei)[0:ballength-unittag])
		var second []byte = make([]byte, unittag)
		copy(second, []byte(wei)[ballength-unittag:])
		secondStr := strings.TrimRight(string(second), "0")
		if secondStr != "" {
			balance = string(first) + "." + secondStr
		} else {
			balance = string(first)
		}
	} else {
		wei = strings.TrimRight(wei, "0")
		if wei != "" {
			gap := unittag - ballength
			balance = "0."

			for i := 0; i < gap; i++ {
				balance = balance + "0"
			}
			balance = balance + wei
		} else {
			balance = "0"
		}
	}
	return balance
}

// 判断字符串是否为10进制数字，，每个字符必须为数字【1--9】，最多有一个小数点”.“
func IsNumber(value string) bool {
	//todo
	_, err := strconv.ParseFloat(value, 64)
	if err != nil && err.(*strconv.NumError).Err == strconv.ErrSyntax {
		fmt.Println("bad number:" + err.Error())
		return false
	}
	return true
}

func GetFilesName(dir string, suffix string) ([]string, error) {
	var re []string
	abiDir, err := os.Open(dir)
	if err == nil {
		files, err2 := abiDir.Readdir(-1)
		if err2 == nil {
			var name string
			var namelen int
			var suffixlen = len(suffix)
			for _, file := range files {
				name = file.Name()
				namelen = len(name)
				if file.IsDir() || namelen < suffixlen {
					continue
				}
				if name[namelen-suffixlen:namelen] == suffix {
					re = append(re, name)
				}
			}
		} else {
			return re, err2
		}
	} else {
		return re, err
	}
	return re, nil
}
