package ethchain

import (
	balacommon "balawallet/common"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// PersonalSign 个人签名 (personal_sign / eth_sign)
func (e *EthChain) PersonalSign(account, password, data string) string {
	message := strings.TrimSpace(data)

	var msgBytes []byte
	if strings.HasPrefix(message, "0x") || strings.HasPrefix(message, "0X") {
		decoded, err := hex.DecodeString(message[2:])
		if err == nil {
			msgBytes = decoded
		} else {
			msgBytes = []byte(message)
		}
	} else {
		msgBytes = []byte(message)
	}

	// Ethereum Signed Message 前缀
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(msgBytes))
	prefixedMsg := append([]byte(prefix), msgBytes...)
	hash := crypto.Keccak256(prefixedMsg)

	key, err := getKey(account, password)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}
	defer zeroKey(key)

	sig, err := crypto.Sign(hash, key)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	// V 值 + 27
	if sig[64] < 27 {
		sig[64] += 27
	}

	sigHex := "0x" + hex.EncodeToString(sig)
	re, _ := json.Marshal(balacommon.ErrorResult{Code: balacommon.Msg_success, Msg: sigHex})
	return string(re)
}

// SignTypedData EIP-712 类型化数据签名
func (e *EthChain) SignTypedData(account, password, data string) string {
	hash := crypto.Keccak256([]byte(data))

	key, err := getKey(account, password)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}
	defer zeroKey(key)

	sig, err := crypto.Sign(hash, key)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	if sig[64] < 27 {
		sig[64] += 27
	}

	sigHex := "0x" + hex.EncodeToString(sig)
	re, _ := json.Marshal(balacommon.ErrorResult{Code: balacommon.Msg_success, Msg: sigHex})
	return string(re)
}

// SignTransaction 签名交易
func (e *EthChain) SignTransaction(account, password, data string) string {
	var txParams map[string]interface{}
	err := json.Unmarshal([]byte(data), &txParams)
	if err != nil {
		var txArray []map[string]interface{}
		err2 := json.Unmarshal([]byte(data), &txArray)
		if err2 != nil || len(txArray) == 0 {
			re := balacommon.GetResult(balacommon.Msg_param_error, "invalid transaction data")
			return string(re)
		}
		txParams = txArray[0]
	}

	txInfo := fmt.Sprintf("To: %v, Value: %v, Data: %v",
		txParams["to"], txParams["value"], txParams["data"])
	hash := crypto.Keccak256([]byte(txInfo))

	key, err := getKey(account, password)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}
	defer zeroKey(key)

	sig, err := crypto.Sign(hash, key)
	if err != nil {
		re := balacommon.GetResult(balacommon.Msg_fail, err.Error())
		return string(re)
	}

	if sig[64] < 27 {
		sig[64] += 27
	}

	sigHex := "0x" + hex.EncodeToString(sig)
	re, _ := json.Marshal(balacommon.ErrorResult{Code: balacommon.Msg_success, Msg: sigHex})
	return string(re)
}

// GetChainId 获取当前网络的 Chain ID
func (e *EthChain) GetChainId() string {
	switch balacommon.G_NetType {
	case balacommon.Nettype_ethereum:
		return "1"
	case balacommon.Nettype_polygon:
		return "137"
	case balacommon.Nettype_pharos:
		return "688688"
	default:
		return "1"
	}
}
