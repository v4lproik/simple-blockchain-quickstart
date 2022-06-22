package balances

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type BalancesSerializer struct {
	balances map[models.Account]uint
}

type BalanceResponse struct {
	Account Account `json:"account"`
	Value   uint    `json:"value"`
}

type Account string

func (t BalancesSerializer) Response() []BalanceResponse {
	balances := t.balances

	response := make([]BalanceResponse, len(balances))
	i := 0
	for balance, val := range balances {
		response[i] = BalanceResponse{
			Account: Account(balance),
			Value:   val,
		}
		i++
	}
	return response
}
