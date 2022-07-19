package services

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type BlockService interface {
	GetNextBlocksFromHash(models.Hash) ([]models.Block, error)
}

type FileBlockService struct {
	mu sync.Mutex
	db *os.File

	miningComplexity uint32
}

func NewFileBlockService(transactionFilePath string, miningComplexity uint32) (*FileBlockService, error) {
	db, err := os.OpenFile(transactionFilePath, os.O_APPEND|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("NewFileBlockService: cannot open txs database: %w", err)
	}
	return &FileBlockService{
		mu:               sync.Mutex{},
		db:               db,
		miningComplexity: miningComplexity,
	}, nil
}

func (a *FileBlockService) GetNextBlocksFromHash(from models.Hash) ([]models.Block, error) {
	a.mu.Lock()

	blocks := make([]models.Block, 0)
	hasFoundHash := false

	// for each block found in database
	scanner := bufio.NewScanner(a.db)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return blocks, fmt.Errorf("GetNextBlocksFromHash: error while scanning: %w", err)
		}

		blockFsJson := scanner.Bytes()
		var blockDB models.BlockDB
		err := json.Unmarshal(blockFsJson, &blockDB)
		if err != nil {
			return blocks, err
		}

		if hasFoundHash {
			blocks = append(blocks, blockDB.Block)
			continue
		}

		if from == blockDB.Hash {
			hasFoundHash = true
		}
	}

	_, err := a.db.Seek(0, io.SeekStart)
	if err != nil {
		return blocks, fmt.Errorf("GetNextBlocksFromHash: couldn't reset pointer on dbfile: %w", err)
	}

	a.mu.Unlock()

	return blocks, nil
}

type PendingBlock struct {
	parent models.Hash
	height uint64
	time   uint64
	miner  models.Account
	txs    []models.Transaction
}

// Mine mines a pending block meaning that it'll try to find a valid nonce
// so it can create a block in the blockchain
func (a *FileBlockService) Mine(ctx context.Context, pb PendingBlock) (*models.Block, error) {
	var block *models.Block
	if len(pb.txs) == 0 {
		return nil, errors.New("Mine: cannot mine block with empty transaction")
	}

	count := uint32(0)
	for {
		select {
		case <-ctx.Done():
			return block, fmt.Errorf("Mine: mining task has been shutdown")
		default:
		}

		nonce := utils.GenerateNonce()
		block = &models.Block{
			Header: models.BlockHeader{
				Parent: pb.parent,
				Height: pb.height,
				Nonce:  nonce,
				Time:   pb.time,
			},
			Txs: pb.txs,
		}

		blockHash, err := block.Hash()
		if err != nil {
			// notest
			return block, fmt.Errorf("Mine: failed to get block hash: %w", err)
		}

		var i uint32
		isValid := true
		for i = 0; i < a.miningComplexity; i++ {
			if fmt.Sprintf("%x", blockHash[i]) != "0" {
				isValid = false
				break
			}
		}
		printAttempts(count)
		count++

		if isValid {
			Logger.Infof("Mine: attempt %d found a nonce=%d, block hash=%s", count, nonce, blockHash.Hex())
			return block, nil
		}
	}
}

func printAttempts(i uint32) {
	if i%1000000 == 0 {
		Logger.Debugf("attempt: %d", i+1)
	}
}
