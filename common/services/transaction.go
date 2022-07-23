package services

import (
	"errors"
	"fmt"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

var (
	ErrMarshalTx       = errors.New("marshal error")
	ErrTxAlreadyInPool = errors.New("transaction is already in pool")
)

type TransactionService interface {
	AddTx(*models.Transaction) error
	addTx(models.Transaction) error
	RemoveTx(models.TransactionId)

	NewPendingTxs() chan models.Transaction
}

type FileTransactionService struct {
	pendingTxPool map[models.TransactionId]models.Transaction

	newPendingTxs chan models.Transaction
}

// NewFileTransactionService default constructor
func NewFileTransactionService() *FileTransactionService {
	return &FileTransactionService{
		pendingTxPool: make(map[models.TransactionId]models.Transaction),
		newPendingTxs: make(chan models.Transaction, 2),
	}
}

// AddTx adds a transaction to the pool
func (a *FileTransactionService) AddTx(tx *models.Transaction) error {
	if tx == nil {
		return errors.New("AddTx: nil transaction")
	}

	// add transaction to the pool
	err := a.addTx(*tx)
	if err != nil {
		return fmt.Errorf("AddTx: failed to add transaction to mempool: %w", err)
	}

	// send the event that there's a new tx ready to be mined
	a.newPendingTxs <- *tx
	return nil
}

func (a *FileTransactionService) addTx(tx models.Transaction) error {
	hash, err := tx.Hash()
	if err != nil {
		return fmt.Errorf("addTx: %w: %s", ErrMarshalTx, err.Error())
	}

	if _, ok := a.pendingTxPool[hash]; ok {
		return fmt.Errorf("addTx: %w", ErrTxAlreadyInPool)
	}

	a.pendingTxPool[hash] = tx

	return nil
}

// RemoveTx remove transaction from pool
func (a *FileTransactionService) RemoveTx(id models.TransactionId) {
	delete(a.pendingTxPool, id)
}

func (a *FileTransactionService) NewPendingTxs() chan models.Transaction {
	return a.newPendingTxs
}
