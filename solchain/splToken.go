package solchain

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
	case "mintInfo":
		return mintInfo(mint, client)
	default:
	}

	return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_param_error))
}

func InvokeTokenContract(mint solana.PublicKey, method abi.Method, params []string, payer solana.PublicKey, pass string, client *rpc.Client) (string, error) {
	switch method.Name {
	case "createTokenAccount":
		owner, err := solana.PublicKeyFromBase58(params[0])
		if err != nil {
			return "", err
		}
		return createTokenAccount(mint, owner, payer, pass, client)
	case "transfer":
		to, err := solana.PublicKeyFromBase58(params[0])
		if err != nil {
			return "", err
		}
		amount, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return "", err
		}
		return transfer(mint, payer, to, amount, pass, client)
	case "mintTo":
		to, err := solana.PublicKeyFromBase58(params[0])
		if err != nil {
			return "", err
		}
		amount, err := strconv.ParseUint(params[1], 10, 64)
		if err != nil {
			return "", err
		}
		return mintTo(mint, to, payer, amount, pass, client)
	default:
	}

	return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_param_error))
}

func mintInfo(mint solana.PublicKey, client *rpc.Client) (string, error) {
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	out, err := client.GetAccountInfo(ctx, mint)
	var result []string
	if err != nil {
		return "", err
	} else {
		fmt.Println("mint account:" + out.Value.Owner.String())
		min := new(token.Mint)
		data2 := out.Value.Data.GetBinary()
		dec2 := bin.NewBinDecoder(data2)
		err = dec2.Decode(&min)
		if err != nil {
			return "", err
		}
		result = append(result, min.MintAuthority.String())
		result = append(result, strconv.FormatUint(min.Supply, 10))
		result = append(result, strconv.FormatUint(uint64(min.Decimals), 10))
		if min.IsInitialized {
			result = append(result, "true")
		} else {
			result = append(result, "false")
		}
		result = append(result, min.FreezeAuthority.String())
	}
	if len(result) != 5 {
		return "", fmt.Errorf(balacommon.GetMsgByCode(balacommon.Msg_account_noexist))
	}
	outs, err := json.Marshal(result)
	return string(outs), err
}

func balanceOf(mint, account solana.PublicKey, client *rpc.Client) (string, error) {
	var result []string
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	out, err := client.GetAccountInfo(ctx, account)
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

func transfer(mint, from, to solana.PublicKey, amount uint64, pass string, client *rpc.Client) (string, error) {

	var conf rpc.GetTokenAccountsConfig
	conf.Mint = &mint
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	out, err := client.GetTokenAccountsByOwner(ctx, from, &conf, nil)
	if err != nil {
		return "", err
	}
	var fromAcc solana.PublicKey
	if len(out.Value) >= 1 {
		fromAcc = out.Value[0].Pubkey
	} else {
		return "", fmt.Errorf("from owner dont have token account")
	}
	ctx1, fuc1 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc1()
	out, err = client.GetTokenAccountsByOwner(ctx1, to, &conf, nil)
	if err != nil {
		return "", err
	}
	var toAcc solana.PublicKey
	if len(out.Value) >= 1 {
		toAcc = out.Value[0].Pubkey
	} else {
		return "", fmt.Errorf("to owner dont have token account")
	}

	transferIns, err := token.NewTransferInstruction(
		amount,
		fromAcc,
		toAcc,
		from,
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return "", err
	}
	ctx2, fuc2 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc2()
	recent, err := client.GetLatestBlockhash(
		ctx2,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return "", err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{transferIns},
		recent.Value.Blockhash,
		solana.TransactionPayer(from),
	)
	if err != nil {
		panic(err)
	}

	// 对交易进行签名
	_, err = signTx(from.String(), pass, tx)
	if err != nil {
		return "", err
	}

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
	ctx3, fuc3 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc3()
	sig, err := client.SendTransaction(
		ctx3,
		tx,
	)
	if err != nil {
		return "", err
	}
	return sig.String(), nil
}

func getTokenAccount(owner solana.PublicKey, mint solana.PublicKey, client *rpc.Client) (string, error) {

	var conf rpc.GetTokenAccountsConfig
	conf.Mint = &mint
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	out, err := client.GetTokenAccountsByOwner(ctx, owner, &conf, nil)
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

func createTokenAccount(mint, owner, payer solana.PublicKey, pass string, client *rpc.Client) (string, error) {

	createins, err := associatedtokenaccount.NewCreateInstruction(
		payer,
		owner,
		mint,
	).ValidateAndBuild()
	if err != nil {
		return "", err
	}

	// 创建交易
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	recent, err := client.GetLatestBlockhash(
		ctx,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return "", err
	}
	tx, err := solana.NewTransaction([]solana.Instruction{createins}, recent.Value.Blockhash, solana.TransactionPayer(payer))
	if err != nil {
		return "", err
	}

	// 对交易进行签名
	_, err = signTx(payer.String(), pass, tx)
	if err != nil {
		return "", err
	}

	// 发送交易
	ctx1, fuc1 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc1()
	sig, err := client.SendTransaction(ctx1, tx)
	if err != nil {
		return "", err
	}

	return sig.String(), nil

}

func mintTo(mint, owner, payer solana.PublicKey, amount uint64, pass string, client *rpc.Client) (string, error) {
	//通过owner 获取代币关联账号
	var conf rpc.GetTokenAccountsConfig
	conf.Mint = &mint
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	out, err := client.GetTokenAccountsByOwner(ctx, owner, &conf, nil)
	if err != nil {
		return "", err
	}
	var to solana.PublicKey
	if len(out.Value) >= 1 {
		to = out.Value[0].Pubkey
	} else {
		return "", fmt.Errorf("to owner dont have token account")
	}
	// 铸造代币指令
	mintToInstruction, err := token.NewMintToInstruction(
		amount,
		mint,
		to,
		payer,
		[]solana.PublicKey{},
	).ValidateAndBuild()
	if err != nil {
		return "", err
	}

	// 创建交易
	ctx1, fuc1 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc1()
	recent, err := client.GetLatestBlockhash(
		ctx1,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return "", err
	}
	tx, err := solana.NewTransaction([]solana.Instruction{mintToInstruction}, recent.Value.Blockhash, solana.TransactionPayer(payer))
	if err != nil {
		return "", err
	}

	// 对交易进行签名
	// 对交易进行签名
	_, err = signTx(payer.String(), pass, tx)
	if err != nil {
		return "", err
	}

	// 发送交易
	fmt.Println("send tx:")
	ctx2, fuc2 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc2()
	sig, err := client.SendTransaction(ctx2, tx)
	if err != nil {
		return "", err
	}

	return sig.String(), nil
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
	ctx, fuc := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc()
	recent, err := client.GetLatestBlockhash(
		ctx,
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
	ctx1, fuc1 := context.WithTimeout(context.Background(), balacommon.G_Timeout)
	defer fuc1()
	sig, err := client.SendTransaction(ctx1, tx)
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
