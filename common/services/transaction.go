package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

var (
	ErrMarshalTx       = errors.New("marshal error")
	ErrTxAlreadyInPool = errors.New("transaction is already in pool")
)

type TransactionService interface {
	AddPendingTx(models.Transaction) error
	GetPendingTxs() map[models.TransactionId]models.Transaction
	RemovePendingTx(models.TransactionId)
	RemovePendingTxs([]models.TransactionId)
}

type FileTransactionService struct {
	mu sync.Mutex

	pendingTxPool map[models.TransactionId]models.Transaction
}

// NewFileTransactionService default constructor
func NewFileTransactionService() *FileTransactionService {
	return &FileTransactionService{
		pendingTxPool: make(map[models.TransactionId]models.Transaction),
	}
}

// AddPendingTx adds a transaction to the pool
func (a *FileTransactionService) AddPendingTx(tx models.Transaction) error {
	// add transaction to the pool
	err := a.addPendingTxToPool(tx)
	if err != nil {
		return fmt.Errorf("AddPendingTx: failed to add transaction to mempool: %w", err)
	}

	return nil
}

// addPendingTxToPool check whether it's a valid transaction and eventually add it to the pool
func (a *FileTransactionService) addPendingTxToPool(tx models.Transaction) error {
	hash, err := tx.Hash()
	if err != nil {
		return fmt.Errorf("addPendingTxToPool: %w: %s", ErrMarshalTx, err.Error())
	}

	// TODO: Implement verifyTx(tx models.Transaction).
	// Needs models.state refacto

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.pendingTxPool[hash]; ok {
		return fmt.Errorf("addPendingTxToPool: %w", ErrTxAlreadyInPool)
	}

	a.pendingTxPool[hash] = tx

	return nil
}

// GetPendingTxs get pending transactions
func (a *FileTransactionService) GetPendingTxs() map[models.TransactionId]models.Transaction {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.pendingTxPool
}

// RemovePendingTx remove transaction from pool
func (a *FileTransactionService) RemovePendingTx(id models.TransactionId) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.pendingTxPool, id)
}

// RemovePendingTxs remove transactions from pool
func (a *FileTransactionService) RemovePendingTxs(ids []models.TransactionId) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, id := range ids {
		delete(a.pendingTxPool, id)
	}
}
