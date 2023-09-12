package manager

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/rain-zxn/transfer-address/hdutil"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
)

/***
BatchCreateAccount
*/
func BatchCreateAccount(ctx *cli.Context) (err error) {
	log.Info("start BatchCreateAccount ...")
	count := ctx.Uint("count")

	log.Info("BatchCreateAccount:", "masterCount", count)

	var mnemonic string
	file, err := os.Open(mnemonicfilename)
	if err != nil {
		log.Info("Open(mnemonicfilename)", "err", err)
	} else {
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			panic(fmt.Errorf("ReadAll(file)err:%v", err))
		} else {
			mnemonic = string(content)
		}
	}
	log.Info("file mnemonic", "mnemonic", mnemonic)

	if mnemonic == "" {
		mnemonic, err = hdutil.GenerateMnemonic()
		if err != nil {
			log.Error("hdutil.GenerateMnemonic", "err", err)
		}
	}

	if _, err := os.Stat(accountFile); os.IsNotExist(err) {
		os.MkdirAll(accountFile, os.ModePerm)
	}

	err = ioutil.WriteFile(mnemonicfilename, []byte(mnemonic), os.ModePerm)
	if err != nil {
		log.Error("WriteFile masteraddrs err:", "err", err)
	}

	log.Info("BatchCreateAccount:", "mnemonic", mnemonic)

	//write mnemonic
	extendedKey, err := hdutil.MnemonicToExtendkey(mnemonic)
	if err != nil {
		log.Error("hdutil.MnemonicToExtendkey err:", "err", err)
	}

	subAddrs := make([]*AddrKey, 0)

	//master addr
	masterExtendedKey, err := hdutil.DeriveAccountFromMaster(extendedKey, 0)
	if err != nil {
		log.Error("hdutil.DeriveAccountFromMaster err:", "err", err, "masterindex", 0)
	}
	for i := uint32(0); i < uint32(count); i++ {
		//subaddr
		subExtendedKey, err := hdutil.DeriveSubAccountFromAccount(masterExtendedKey, i)
		if err != nil {
			log.Error("hdutil.DeriveSubAccountFromAccount err:", "err", err, "masterindex", 0, "subindex", i)
		}
		//subprivateKey
		privateKey, err := hdutil.ExtendedKeyToHex(subExtendedKey)
		if err != nil {
			log.Error("hdutil.ExtendedKeyToHex err:", "err", err, "masterindex", 0, "subindex", i)
		}
		addr, err := hdutil.ExtendedKeyToAddress(subExtendedKey)
		if err != nil {
			log.Error("hdutil.ExtendedKeyToAddress err:", "err", err, "masterindex", 0, "subindex", i)
		}
		//wraite addr privateKey
		subAddrs = append(subAddrs, &AddrKey{
			addr,
			privateKey,
		})
	}

	jsonSub0Addrs, _ := json.MarshalIndent(subAddrs, "", "	")
	err = ioutil.WriteFile(accountFile+"_sub0addrs.json", jsonSub0Addrs, os.ModePerm)
	if err != nil {
		log.Error("WriteFile subAddrs err:", "err", err)
	}
	return nil
}
