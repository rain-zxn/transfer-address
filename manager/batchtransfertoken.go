package manager

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blocktree/go-owcrypt/sha3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

var (
	addrNonce uint64
	SumLimit  uint64
	SumGas    uint64
)

/***
BatchTransferToken
*/
func BatchTransferToken(ctx *cli.Context) (err error) {
	log.Info("start BatchTransferToken ...")

	starttosubaccount := uint32(ctx.Int64("starttosubaccount"))
	endtosubaccount := uint32(ctx.Int64("endtosubaccount"))
	token := ctx.String("token")
	tokenAmount := ctx.String("tokenamount")
	estimate := ctx.Bool("estimate")

	log.Info("BatchTransferToken:", "starttosubaccount", starttosubaccount)
	log.Info("BatchTransferToken:", "endtosubaccount", endtosubaccount)
	log.Info("BatchTransferToken:", "token", token)
	log.Info("BatchTransferToken:", "tokenamount", tokenAmount)

	if !strings.EqualFold(token, USDT) && !strings.EqualFold(token, MATIC) {
		panic("token need usdt or matic")
	}

	exponent := new(big.Int)
	exponent.Exp(big.NewInt(10), big.NewInt(18), nil)
	x, ok := new(big.Int).SetString(tokenAmount, 10)
	if !ok {
		panic("tokenAmount SetString err")
	}
	amount := new(big.Int).Mul(x, exponent)

	var mnemonic string
	file, err := os.Open(mnemonicfilename)
	if err != nil {
		panic(fmt.Sprintf("Open(mnemonicfilename)err:%v", err))
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("ReadAll(file) err:%v", err))
	}
	mnemonic = string(content)
	if len(mnemonic) == 0 {
		panic(fmt.Errorf("mnemonic is nil"))
	}

	srcAddr, err := GetIndexAddrWithMnemonic(0, 0, mnemonic)
	if err != nil {
		panic(fmt.Errorf("GetIndexAddrWithMnemonic srcAddr err:%v, i:%v, j:%v", err, 0, 0))
	}
	targetAddrs := make([]*AddrKey, 0)
	for j := starttosubaccount; j <= endtosubaccount; j++ {
		targetAddr, err := GetIndexAddrWithMnemonic(0, j, mnemonic)
		if err != nil {
			panic(fmt.Errorf("GetIndexAddrWithMnemonic targetAddr err:%v, i:%v, j:%v", err, 0, j))
		}
		targetAddrs = append(targetAddrs, targetAddr)
	}
	return transferToken(srcAddr, targetAddrs, token, amount, estimate)

	return nil
}

func transferToken(srcAddr *AddrKey, targetAddrs []*AddrKey, token string, amount *big.Int, estimate bool) error {

	client, err := ethclient.Dial(node)
	if err != nil {
		panic(fmt.Sprintf("fail to dial client %s ", node))
	}

	txTransfers := make([]*TxTransfer, 0)
	for i, v := range targetAddrs {
		var hash string
		var err error
		switch {
		case strings.EqualFold(token, MATIC):
			hash, err = SendTxAccountOneNonce(client, srcAddr, v.Addr, amount, make([]byte, 0), estimate)
		case strings.EqualFold(token, USDT):
			hash, err = SendTxAccountOneNonce(client, srcAddr, common.HexToAddress(Token_USDT), nil, fillData(v.Addr, amount), estimate)
		}
		if err != nil {
			panic(fmt.Sprintf("SendTxAccount index:%v err: %v", i, err))
		}
		log.Info("TransferToken", "hash", hash, "token", "matic")
		txTransfers = append(txTransfers, &TxTransfer{
			srcAddr.Addr.String(),
			v.Addr.String(),
			"matic",
			amount.String(),
			hash,
		})

	}
	log.Info("SumGas", "SumGas", SumGas)

	transferTokenFile := "batchtransfertoken/" + token + "/"
	if _, err := os.Stat(transferTokenFile); os.IsNotExist(err) {
		os.MkdirAll(transferTokenFile, os.ModePerm)
	}

	jsontxTransferTokens, _ := json.MarshalIndent(txTransfers, "", "	")
	err = ioutil.WriteFile(transferTokenFile+"_"+srcAddr.Addr.String(), jsontxTransferTokens, os.ModePerm)
	if err != nil {
		log.Error("WriteFile batchtransfertoken err:", "err", err, "srcAddr", srcAddr.Addr.String())
	}

	return nil
}

func SendTxAccountOneNonce(client *ethclient.Client, senderAddr *AddrKey, to common.Address, amount *big.Int, data []byte, estimate bool) (hash string, err error) {
	balance, _ := client.BalanceAt(context.Background(), senderAddr.Addr, nil)
	log.Info("SendTxAccount EstimateGas", "addr", senderAddr.Addr, "balance", balance)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("SuggestGasPrice err: %v", err)
	}
	if addrNonce == 0 {
		nonce, err := client.NonceAt(context.Background(), senderAddr.Addr, nil)
		if err != nil {
			return "", fmt.Errorf("PendingNonceAt err, addr %s, err %v", senderAddr.Addr, err)
		}
		addrNonce = nonce
	}
	msg := ethereum.CallMsg{From: senderAddr.Addr, To: &to, Value: amount, Data: data}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return "", fmt.Errorf("EstimateGas err, addr %s, err %v", senderAddr.Addr, err)
	}
	log.Info("SendTxAccount EstimateGas", "gasLimit", gasLimit)
	log.Info("SendTxAccount addrNonce", "addrNonce", addrNonce)
	gasLimit = uint64(1.001 * float32(gasLimit))
	SumLimit += gasLimit
	log.Info("SendTxAccount EstimateGas", "gasLimit*1.001", gasLimit)
	log.Info("SendTxAccount gasPrice", "gasPrice", gasPrice)
	SumGas += new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))).Uint64()
	log.Info("SendTxAccount gas", "gas", new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit))))
	if estimate {
		return "", nil
	}

	tx := types.NewTransaction(addrNonce, to, amount, gasLimit, gasPrice, data)
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint("fail to chainid "), err)
	}
	signer := types.LatestSignerForChainID(chainId)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint("fail to signer ", err))

	}
	privateKey, err := crypto.HexToECDSA(senderAddr.PrivateKey)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("crypto HexToECDSA err:%v,addr:%v", err, senderAddr.Addr))
	}
	tx, err = types.SignTx(tx, signer, privateKey)
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprint("SendTransaction err", err))
	}
	log.Info("SendTxAccount", "srcaddr:", senderAddr.Addr, "toaddr", to, "amount", amount, "data:", hex.EncodeToString(data))
	addrNonce++
	return tx.Hash().String(), err
}

func fillData(toAddress common.Address, amount *big.Int) []byte {
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	return data
}
