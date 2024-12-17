package common

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
)

var G_Timeout time.Duration = 5 * time.Second
var G_Config = make(map[string]string)
var G_Etherclient *ethclient.Client = nil
var G_BlockchainConfig bool = true
var G_DataDir string
var G_NetType int

const (
	Nettype_ethereum int = 0
	Nettype_polygon  int = 1
	Nettype_solana   int = 2
)

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
	DoGetConfig("")
	if G_Config["ethAddress"] != "" {
		G_Etherclient, err = ethclient.Dial(G_Config["ethAddress"])
		if err != nil {
			G_Etherclient = nil
			fmt.Println("can not connect to node")
		}
	} else {
		fmt.Println("have not config blockchain node")
		G_BlockchainConfig = false
	}
}

type ErrorResult struct {
	Code uint
	Msg  string
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

func GetMsg(code uint) string {
	lang := G_Config["language"]
	if lang == "" {
		lang = "zh-CN"
	}
	msg := G_Msg[lang][code]
	return msg
}
