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
	AddTx(*models.Transaction) error
	GetTxs() map[models.TransactionId]models.Transaction
	RemoveTx(models.TransactionId)
	RemoveTxs([]models.TransactionId)
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

	return nil
}

func (a *FileTransactionService) addTx(tx models.Transaction) error {
	hash, err := tx.Hash()
	if err != nil {
		return fmt.Errorf("addTx: %w: %s", ErrMarshalTx, err.Error())
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.pendingTxPool[hash]; ok {
		return fmt.Errorf("addTx: %w", ErrTxAlreadyInPool)
	}

	a.pendingTxPool[hash] = tx

	return nil
}

// GetTxs get pending transactions
func (a *FileTransactionService) GetTxs() map[models.TransactionId]models.Transaction {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.pendingTxPool
}

// RemoveTx remove transaction from pool
func (a *FileTransactionService) RemoveTx(id models.TransactionId) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.pendingTxPool, id)
}

// RemoveTx remove transactions from pool
func (a *FileTransactionService) RemoveTxs(ids []models.TransactionId) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, id := range ids {
		delete(a.pendingTxPool, id)
	}
}
