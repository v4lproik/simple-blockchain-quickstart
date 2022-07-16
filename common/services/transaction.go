package services

import (
	"errors"
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
		return nil, errors.New("AddTransaction: nil state")
	}

	if tx == nil {
		return nil, errors.New("AddTransaction: nil transaction")
	}

	//add transaction to state
	err := state.Add(*tx)
	if err != nil {
		return nil, fmt.Errorf("AddTransaction: failed to add tx: %s", err)
	}

	//persist new state to disk
	hash, err := state.Persist()
	if err != nil {
		return nil, fmt.Errorf("AddTransaction: failed to persist state: %s", err)
	}

	return &hash, nil
}
