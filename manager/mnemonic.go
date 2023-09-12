package manager

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/rain-zxn/transfer-address/hdutil"
)

func GetIndexAddrWithMnemonic(masterIndex uint32, subIndex uint32, mnemonic string) (*AddrKey, error) {
	extendedKey, err := hdutil.MnemonicToExtendkey(mnemonic)
	if err != nil {
		log.Error("hdutil.MnemonicToExtendkey err:", "err", err)
	}
	masterExtendedKey, err := hdutil.DeriveAccountFromMaster(extendedKey, masterIndex)
	if err != nil {
		log.Error("hdutil.DeriveAccountFromMaster err:", "err", err, "masterindex", masterIndex)
	}
	subExtendedKey, err := hdutil.DeriveSubAccountFromAccount(masterExtendedKey, subIndex)
	if err != nil {
		log.Error("hdutil.DeriveSubAccountFromAccount err:", "err", err, "masterindex", masterIndex, "subindex", subIndex)
	}
	//subprivateKey
	privateKey, err := hdutil.ExtendedKeyToHex(subExtendedKey)
	if err != nil {
		log.Error("hdutil.ExtendedKeyToHex err:", "err", err, "masterindex", masterIndex, "subindex", subIndex)
	}
	addr, err := hdutil.ExtendedKeyToAddress(subExtendedKey)
	if err != nil {
		log.Error("hdutil.ExtendedKeyToAddress err:", "err", err, "masterindex", masterIndex, "subindex", subIndex)
	}
	return &AddrKey{
		addr,
		privateKey,
	}, nil
}
