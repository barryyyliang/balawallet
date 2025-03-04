package ethchain

import (
	balacommon "balawallet/common"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strings"

	"balawallet/accounts/eth"
	"balawallet/accounts/keystore"

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

func weiToether(wei string) string {
	return balacommon.ToBig(18, wei)
}

func weiTogwei(wei string) string {
	return balacommon.ToBig(9, wei)
}

func etherTowei(ether string) (big.Int, error) {
	// 要调用isnumber判断是否为正确数字格式
	if balacommon.IsNumber(ether) {
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

func zeroKey(key *ecdsa.PrivateKey) {
	b := key.D.Bits()
	clear(b)
}

func getKey(add string, pa string) (*ecdsa.PrivateKey, error) {
	selectAccountFile := ""
	acc := ks.Accounts()
	for i := 0; i < len(acc); i++ {
		if acc[i].GetAddress() == add {
			selectAccountFile = acc[i].GetUrl().Path
		}
	}
	if selectAccountFile != "" {
		keyjson, err := os.ReadFile(selectAccountFile)
		if err != nil {
			return nil, err
		}
		ethAdapter := &eth.EthChainAdapter{}
		key, err := keystore.DecryptKey(keyjson, pa, ethAdapter)
		if err != nil {
			return nil, err
		}
		// 确保导出的是需要私钥
		if key.GetStringAddress() != add {
			fmt.Errorf("key content mismatch: have account %x, want %s", key.GetStringAddress(), add)
			return nil, errors.New(balacommon.GetMsgByCode(balacommon.Msg_addressprivatekey_mismatch))
		} else {
			pri := key.GetPrivateKey()
			if prikey, ok := pri.(*ecdsa.PrivateKey); ok {
				return prikey, nil
			} else {
				return nil, errors.New(balacommon.GetMsgByCode(balacommon.Msg_readaccountfile_fail))
			}
		}
	} else {
		fmt.Println("can not find account:" + add)
	}
	return nil, errors.New(balacommon.GetMsgByCode(balacommon.Msg_addressprivatekey_mismatch))
}
