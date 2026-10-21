package solchain

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"balawallet/accounts/keystore"
	"balawallet/accounts/sol"
	balacommon "balawallet/common"

	"github.com/gagliardetto/solana-go"
)

func signTx(address, auth string, tx *solana.Transaction) (*solana.Transaction, error) {
	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return tx, err
	}
	if len(tx.Message.AccountKeys) == 0 {
		return tx, fmt.Errorf("accountkeys is null")
	}

	if address != string(tx.Message.AccountKeys[0].String()) {
		return tx, fmt.Errorf("account is not the required signer by tx")
	}

	if len(tx.Signatures) == 0 {
		tx.Signatures = make([]solana.Signature, 1)
	}

	signData, err := ks.SignMessageWithPassphrase(address, auth, messageContent)

	if err != nil {
		return nil, err
	}

	var signature solana.Signature
	copy(signature[:], signData)
	tx.Signatures[0] = signature

	return tx, nil
}

func solTolam(sol string) (uint64, error) {
	// 要调用isnumber判断是否为正确数字格式
	if balacommon.IsNumber(sol) {
		strs := strings.Split(sol, ".")
		var balstr string
		if len(strs) == 1 {
			if balstr != "0" {
				balstr = strs[0] + "000000000"
			}
		} else {
			lens := len(strs[1])

			for i := 0; i < 9-lens; i++ {
				strs[1] = strs[1] + "0"
			}

			if "0" == strs[0] {
				balstr = strings.TrimLeft(strs[1], "0")
			} else {
				balstr = strs[0] + strs[1]
			}
		}

		lam, err := strconv.ParseUint(balstr, 10, 64)
		if err != nil {
			return 0, err
		}

		return lam, nil
	}
	return 0, fmt.Errorf("bad number")

}

func getKey(add string, pa string) (*ed25519.PrivateKey, error) {
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
		ethAdapter := &sol.SolChainAdapter{}
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
			if prikey, ok := pri.(*ed25519.PrivateKey); ok {
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
