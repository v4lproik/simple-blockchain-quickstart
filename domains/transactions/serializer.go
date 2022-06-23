package transactions

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type TransactionSerializer struct {
	hash        models.Hash
	transaction models.Transaction
}

type TransactionResponse struct {
	Transaction struct {
		Hash   models.Hash    `json:"hash"`
		From   models.Account `json:"from"`
		To     models.Account `json:"to"`
		Value  uint           `json:"value"`
		Reason string         `json:"reason"`
	} `json:"transaction"`
}

func (t TransactionSerializer) Response() TransactionResponse {
	return TransactionResponse{
		Transaction: struct {
			Hash   models.Hash    `json:"hash"`
			From   models.Account `json:"from"`
			To     models.Account `json:"to"`
			Value  uint           `json:"value"`
			Reason string         `json:"reason"`
		}{
			t.hash,
			t.transaction.From,
			t.transaction.To,
			t.transaction.Value,
			t.transaction.Reason,
		},
	}
}
