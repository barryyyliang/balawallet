package spltoken

import (
	balacommon "balawallet/common"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/text"
	//"github.com/ethereum/go-ethereum/crypto/sha3"
)

func ReadTokenContract(mint solana.PublicKey, method abi.Method, params []string, client *rpc.Client) (string, error) {
	switch method.Name {
	case "balanceOf":
		account, err := solana.PublicKeyFromBase58(params[0])
		if err != nil {
			return "", err
		}
		return balanceOf(mint, account, client)
	case "getTokenAccount":
		owner, err := solana.PublicKeyFromBase58(params[0])
		if err != nil {
			return "", err
		}
		return getTokenAccount(owner, mint, client)
	default:
	}

	return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_param_error))
}

func Decimals(mint solana.PublicKey, client *rpc.Client) string {

	return ""
}
func TotalSupply(mint solana.PublicKey, client *rpc.Client) string {

	return ""
}
func balanceOf(mint, account solana.PublicKey, client *rpc.Client) (string, error) {
	var result []string
	out, err := client.GetAccountInfo(context.TODO(), account)
	if err != nil {
		return "", err
	} else {
		var tk token.Account
		data := out.Value.Data.GetBinary()
		dec := bin.NewBinDecoder(data)
		err = dec.Decode(&tk)
		if err != nil {
			return "", err
		}
		if mint.Equals(tk.Mint) {
			result = append(result, account.String())
			result = append(result, tk.Mint.String())
			result = append(result, tk.Owner.String())
			result = append(result, strconv.FormatUint(tk.Amount, 10))
		} else {
			return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_account_mint_mismatch))
		}
	}
	if len(result) != 4 {
		return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_account_noexist))
	}
	outs, err := json.Marshal(result)
	return string(outs), err
}

func Transfer(from *solana.PrivateKey, to solana.PublicKey, amount uint64, client *rpc.Client) {
	recent, err := client.GetLatestBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		panic(err)
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				from.PublicKey(),
				to,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(from.PublicKey()),
	)
	if err != nil {
		panic(err)
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if from.PublicKey().Equals(key) {
				return from
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Errorf("unable to sign transaction: %w", err))
	}
	spew.Dump(tx)
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// // Send transaction, and wait for confirmation:
	// wsClient, err := ws.Connect(context.Background(), rpc.LocalNet_WS)
	// if err != nil {
	// 	panic(err)
	// }
	// sig, err := confirm.SendAndConfirmTransaction(
	// 	context.TODO(),
	// 	client,
	// 	wsClient,
	// 	tx,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// spew.Dump(sig)

	// Or just send the transaction WITHOUT waiting for confirmation:
	sig, err := client.SendTransaction(
		context.TODO(),
		tx,
	)
	if err != nil {
		panic(err)
	}
	spew.Dump(sig)
}

func getTokenAccount(owner solana.PublicKey, mint solana.PublicKey, client *rpc.Client) (string, error) {

	var conf rpc.GetTokenAccountsConfig
	conf.Mint = &mint
	out, err := client.GetTokenAccountsByOwner(context.TODO(), owner, &conf, nil)
	if err != nil {
		return "", err
	}

	var result []string
	for _, value := range out.Value {
		result = append(result, value.Pubkey.String())
		var token token.Account
		data := value.Account.Data.GetBinary()
		dec := bin.NewBinDecoder(data)
		err := dec.Decode(&token)
		if err != nil {
			fmt.Println(err.Error())
		}
		result = append(result, token.Mint.String())
		result = append(result, token.Owner.String())
		result = append(result, strconv.FormatUint(token.Amount, 10))
		break
	}
	if len(result) != 4 {
		return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_account_noexist))
	}
	outs, err := json.Marshal(result)
	return string(outs), err
}

func CreateTokenAccount(mint, owner solana.PublicKey, payer solana.PrivateKey, client *rpc.Client) (solana.PublicKey, error) {

	createTokenAccount, err := associatedtokenaccount.NewCreateInstruction(
		payer.PublicKey(),
		owner,
		mint,
	).ValidateAndBuild()
	if err != nil {
		fmt.Println("create instruction error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 创建交易
	recent, err := client.GetLatestBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		fmt.Println("get recent error")
		return solana.PublicKey{}, err
	}
	tx, err := solana.NewTransaction([]solana.Instruction{createTokenAccount}, recent.Value.Blockhash, solana.TransactionPayer(solana.PublicKey(payer.PublicKey())))
	if err != nil {
		fmt.Println("create NewTransaction error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 对交易进行签名
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(payer.PublicKey()) {
			return &payer
		}
		return nil
	})
	if err != nil {
		fmt.Println("sign error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 发送交易
	fmt.Println("send tx:")
	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		fmt.Println("send tx:" + err.Error())
		return solana.PublicKey{}, err
	}

	fmt.Println("get tx sig:" + sig.String())

	spew.Dump(tx)

	return solana.PublicKey{}, nil

}

func MintTo(mint, to solana.PublicKey, payer solana.PrivateKey, amount uint64, client *rpc.Client) {
	// 铸造代币指令
	mintToInstruction, err := token.NewMintToInstruction(
		amount, // 铸造 1 个代币，精度为 9
		mint,
		to,
		payer.PublicKey(),
		nil,
	).ValidateAndBuild()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 创建交易
	recent, err := client.GetLatestBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		fmt.Println("get recent error")
		return
	}
	tx, err := solana.NewTransaction([]solana.Instruction{mintToInstruction}, recent.Value.Blockhash, solana.TransactionPayer(solana.PublicKey(payer.PublicKey())))
	if err != nil {
		fmt.Println("create NewTransaction error:" + err.Error())
		return
	}

	// 对交易进行签名
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(payer.PublicKey()) {
			return &payer
		}
		return nil
	})
	if err != nil {
		fmt.Println("sign error:" + err.Error())
		return
	}

	// 发送交易
	fmt.Println("send tx:")
	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		fmt.Println("send tx:" + err.Error())
		return
	}

	fmt.Println("get tx sig:" + sig.String())

	spew.Dump(tx)
}

func CreateSPLToken(payer solana.PrivateKey, client *rpc.Client) (solana.PublicKey, error) {

	// 创建一个新的代币地址
	//mint, err := solana.NewRandomKey()
	account := newAccount()
	mint := account.PrivateKey
	fmt.Println("mint:" + mint.PublicKey().String())

	createAccount, err := system.NewCreateAccountInstruction(
		1461600,
		82,
		solana.TokenProgramID,
		//solana.SystemProgramID,
		payer.PublicKey(),
		mint.PublicKey(),
	).ValidateAndBuild()

	if err != nil {
		fmt.Println("create account instruction error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 创建创建代币的交易指令
	createMintInst, err := token.NewInitializeMint2Instruction(
		9,
		payer.PublicKey(),
		payer.PublicKey(),
		mint.PublicKey(),
	).ValidateAndBuild()
	if err != nil {
		fmt.Println("create instruction error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 创建交易
	recent, err := client.GetLatestBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		fmt.Println("get recent error")
		return solana.PublicKey{}, err
	}
	tx, err := solana.NewTransaction([]solana.Instruction{createAccount, createMintInst}, recent.Value.Blockhash, solana.TransactionPayer(solana.PublicKey(payer.PublicKey())))
	if err != nil {
		fmt.Println("create NewTransaction error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 对交易进行签名
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(payer.PublicKey()) {
			return &payer
		}
		if key.Equals(mint.PublicKey()) {
			return &mint
		}
		return nil
	})
	if err != nil {
		fmt.Println("sign error:" + err.Error())
		return solana.PublicKey{}, err
	}

	// 发送交易
	fmt.Println("send tx:")
	sig, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		fmt.Println("send tx:" + err.Error())
		return solana.PublicKey{}, err
	}

	fmt.Println("get tx sig:" + sig.String())

	spew.Dump(tx)
	// Pretty print the transaction:
	tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "mint spl"))

	return mint.PublicKey(), nil
}

func newAccount() *solana.Wallet {
	//Create a new account:
	account := solana.NewWallet()
	fmt.Println("account public key:", account.PublicKey())
	return account
}
