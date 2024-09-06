package main

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	//"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
)

type ErrorResult struct {
	Code int
	Err  string
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
	fmt.Println("wei:" + wei)
	if ballength > unittag {
		var first []byte = make([]byte, ballength-unittag)
		copy(first, []byte(wei)[0:ballength-unittag])
		//fmt.Println("first:" + string(first))
		var second []byte = make([]byte, unittag)
		copy(second, []byte(wei)[ballength-unittag:])
		//fmt.Println("second:" + string(second))
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

//判断字符串是否为10进制数字，，每个字符必须为数字【1--9】，最多有一个小数点”.“
func isNumber(value string) bool {
	//todo
	result, err := strconv.ParseFloat(value, 64)
	fmt.Println(result)
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
		connect, err := ethclient.Dial(G_Config["nodeAddress"])
		if err != nil {
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
		fmt.Println(t.String())
		add := common.HexToAddress(value)
		return add
	case "*big.Int":
		fmt.Println(t.String())
		ret := new(big.Int)
		ret.SetString(value, 10)
		return ret
	}

	return nil
}

func valueToString(value interface{}) string {
	var ret string
	switch reflect.TypeOf(value).String() {
	case "common.Address":
		fmt.Println("common.Address")
		add := value.(common.Address)
		ret = add.String()
	case "*big.Int":
		fmt.Println("*big.Int")
		v := new(big.Int)
		v = value.(*big.Int)
		ret = v.String()
	}
	return ret
}
