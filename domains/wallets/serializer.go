package wallets

import (
	"github.com/ethereum/go-ethereum/common"
)

type WalletSerializer struct {
	account common.Address
}

type WalletResponse struct {
	Account common.Address `json:"account"`
}

func (t WalletSerializer) Response() WalletResponse {
	return WalletResponse{
		Account: t.account,
	}
}
