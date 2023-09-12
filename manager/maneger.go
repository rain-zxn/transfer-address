package manager

import "github.com/ethereum/go-ethereum/common"

const (
	accountFile      = "./accounts/"
	mnemonicfilename = accountFile + "_mnemonic.json"
	USDT             = "usdt"
	MATIC            = "matic"
	Token_USDT       = "0xc2132D05D31c914a87C6611C10748AEb04B58e8F"
	node             = "https://polygon.llamarpc.com"
)

type AddrKey struct {
	Addr       common.Address
	PrivateKey string
}

type IndexAddr struct {
	Addr        common.Address
	MasterIndex uint32
	SubIndex    uint32
	Subi0       common.Address
}

type TxTransfer struct {
	Src    string
	Dst    string
	Token  string
	Amount string
	Hash   string
}
