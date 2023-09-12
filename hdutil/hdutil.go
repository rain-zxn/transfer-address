package hdutil

import (
	"encoding/json"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hd "github.com/miguelmota/go-ethereum-hdwallet"
)

var basePath string = "m/44'/60'/0'/0/0"

func GenerateMnemonic() (string, error) {
	entropy, err := hd.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return hd.NewMnemonicFromEntropy(entropy)
}

func MnemonicToExtendkey(mnemonic string) (*hdkeychain.ExtendedKey, error) {
	seed, err := hd.NewSeedFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	return hdkeychain.NewKeyFromString(masterKey.String())
}

// use masterKey to generate each multiSig account
// path: m/44'/60'/index'
func DeriveAccountFromMaster(masteKey *hdkeychain.ExtendedKey, index uint32) (*hdkeychain.ExtendedKey, error) {
	m_44H, err := masteKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil, err
	}
	m_44H_60H, err := m_44H.Derive(hdkeychain.HardenedKeyStart + 60)
	if err != nil {
		return nil, err
	}
	return m_44H_60H.Derive(hdkeychain.HardenedKeyStart + index)
}

// use multiSig account to generate key for each chain
// recommend to use chainId as subIndex
// path: m/0/subIndex
func DeriveSubAccountFromAccount(accountKey *hdkeychain.ExtendedKey, subIndex uint32) (*hdkeychain.ExtendedKey, error) {
	m_0, err := accountKey.Derive(0)
	if err != nil {
		return nil, err
	}
	return m_0.Derive(subIndex)
}

func ExtendedKeyToHex(key *hdkeychain.ExtendedKey) (string, error) {
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return "", err
	}
	privateKeyEcdsa := privateKey.ToECDSA()
	return hexutil.Encode(crypto.FromECDSA(privateKeyEcdsa))[2:], nil
}

func ExtendedKeyToAddress(key *hdkeychain.ExtendedKey) (common.Address, error) {
	privateKey, err := key.ECPrivKey()
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(privateKey.PublicKey), nil
}

func EncryptAnyKey(key string, auth string, scryptN, scryptP int) ([]byte, error) {
	keyBytes := []byte(key)
	cryptoStruct, err := keystore.EncryptDataV3(keyBytes, []byte(auth), scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	return json.Marshal(cryptoStruct)
}

func DecryptAnyKey(keyjson []byte, auth string) (string, error) {
	cryptoStruct := new(keystore.CryptoJSON)
	if err := json.Unmarshal(keyjson, cryptoStruct); err != nil {
		return "", err
	}
	key, err := keystore.DecryptDataV3(*cryptoStruct, auth)
	if err != nil {
		return "", err
	}
	return string(key), nil
}
