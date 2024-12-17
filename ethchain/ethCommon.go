package ethchain

import (
	. "balawallet/common"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	//"github.com/ethereum/go-ethereum/crypto/sha3"

	"github.com/ethereum/go-ethereum/ethclient"
)

func newDynamicTx(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tipsPrice *big.Int, chain *big.Int) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: tipsPrice,
		ChainID:   chain,
	})
}

func newDeployTx(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tipsPrice *big.Int, chain *big.Int, data []byte) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		Value:     amount,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: tipsPrice,
		ChainID:   chain,
		Data:      data,
	})
}

func newCallContractTx(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, tipsPrice *big.Int, chain *big.Int, data []byte) *types.Transaction {
	return types.NewTx(&types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       gasLimit,
		GasFeeCap: gasPrice,
		GasTipCap: tipsPrice,
		ChainID:   chain,
		Data:      data,
	})
}

func weiTox(unittag int, wei string) string {
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

func weiToether(wei string) string {
	return weiTox(18, wei)
}

func weiTogwei(wei string) string {
	return weiTox(9, wei)
}

func etherTowei(ether string) (big.Int, error) {
	// 要调用isnumber判断是否为正确数字格式
	if isNumber(ether) {
		strs := strings.Split(ether, ".")
		var balstr string
		if len(strs) == 1 {
			if balstr != "0" {
				balstr = strs[0] + "000000000000000000"
			}
		} else {
			lens := len(strs[1])

			for i := 0; i < 18-lens; i++ {
				strs[1] = strs[1] + "0"
			}

			if "0" == strs[0] {
				balstr = strings.TrimLeft(strs[1], "0")
			} else {
				balstr = strs[0] + strs[1]
			}
		}

		var wei big.Int = *big.NewInt(0)
		wei.SetString(balstr, 10)
		return wei, nil
	}
	return *big.NewInt(0), fmt.Errorf("bad number")

}

// 判断字符串是否为10进制数字，，每个字符必须为数字【1--9】，最多有一个小数点”.“
func isNumber(value string) bool {
	//todo
	_, err := strconv.ParseFloat(value, 64)
	if err != nil && err.(*strconv.NumError).Err == strconv.ErrSyntax {
		fmt.Println("bad number:" + err.Error())
		return false
	}
	return true
}

func getChainInfo(connect *ethclient.Client, from string) (uint64, *big.Int, *big.Int, *big.Int, error) {
	ctx := context.Background()
	nounce, err := connect.NonceAt(ctx, common.HexToAddress(from), nil)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	gasPrice, err := connect.SuggestGasPrice(ctx) ///big.NewInt(20000000000)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	tipsPrice, err := connect.SuggestGasTipCap(ctx)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	chainId, err := connect.ChainID(ctx)
	if err != nil {
		return nounce, nil, nil, nil, err
	}
	return nounce, gasPrice, tipsPrice, chainId, nil

}

func getConnection() *ethclient.Client {
	if G_Etherclient == nil {
		var connect *ethclient.Client
		var err error
		if Nettype_ethereum == G_NetType {
			fmt.Println("net type eth")
			connect, err = ethclient.Dial(G_Config["ethAddress"])
		} else if Nettype_polygon == G_NetType {
			fmt.Println("net type sol:" + G_Config["polAddress"])
			connect, err = ethclient.Dial(G_Config["polAddress"])
		}

		if err != nil {
			fmt.Println("net connet fail")
			G_Etherclient = nil

		} else {
			G_Etherclient = connect
		}
	}
	return G_Etherclient
}

func stringToValue(t reflect.Type, value string) interface{} {

	switch t.String() {
	case "common.Address":
		add := common.HexToAddress(value)
		return add
	case "*big.Int":
		ret := new(big.Int)
		ret.SetString(value, 10)
		return ret
	case "string":
		return value
	}
	return nil
}

func valueToString(value interface{}) string {
	var ret string
	switch reflect.TypeOf(value).String() {
	case "common.Address":
		add := value.(common.Address)
		ret = add.String()
	case "*big.Int":
		v := new(big.Int)
		v = value.(*big.Int)
		ret = v.String()
	case "string":
		ret = value.(string)
	case "bool":
		v := value.(bool)
		if v == true {
			ret = "true"
		} else {
			ret = "false"
		}
	}
	return ret
}

func getFiles(dir string, suffix string) ([]string, error) {
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

func zeroKey(key *ecdsa.PrivateKey) {
	b := key.D.Bits()
	clear(b)
}

func getKey(add string, pa string) (*ecdsa.PrivateKey, error) {
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
			return nil, err
		}
		key, err := keystore.DecryptKey(keyjson, pa)
		if err != nil {
			return nil, err
		}
		// 确保导出的是需要私钥
		if key.Address.String() != add {
			fmt.Errorf("key content mismatch: have account %x, want %s", key.Address, add)
			return nil, errors.New(GetMsg(Msg_addressprivatekey_mismatch))
		} else {
			return key.PrivateKey, nil
		}
	} else {
		fmt.Println("can not find account:" + add)
	}
	return nil, errors.New(GetMsg(Msg_addressprivatekey_mismatch))
}
