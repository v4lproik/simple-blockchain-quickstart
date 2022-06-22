package services

import (
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type TransactionService interface {
	AddTransaction(models.State, *models.Transaction) (*models.Hash, error)
}

type FileTransactionService struct {
}

func NewFileTransactionService() FileTransactionService {
	return FileTransactionService{}
}

func (a FileTransactionService) AddTransaction(state models.State, tx *models.Transaction) (*models.Hash, error) {
	if state == nil {
		return nil, fmt.Errorf("cannot add transaction to a nil state")
	}

	if tx == nil {
		return nil, fmt.Errorf("cannot add transaction with a nil transaction")
	}

	//add transaction to state
	err := state.Add(*tx)
	if err != nil {
		return nil, fmt.Errorf("cannot add transaction to state: %v", err)
	}

	//persist new state to disk
	hash, err := state.Persist()
	if err != nil {
		return nil, fmt.Errorf("cannot persist state to disk: %v", err)
	}

	return &hash, nil
}
