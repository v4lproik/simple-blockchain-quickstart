package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"io"
	"os"
	"sync"
)

type BlockService interface {
	GetNextBlocksFromHash(models.Hash) ([]models.Block, error)
}

type FileBlockService struct {
	mu sync.Mutex
	db *os.File
}

func NewFileBlockService(transactionFilePath string) (*FileBlockService, error) {
	db, err := os.OpenFile(transactionFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction file path")
	}
	return &FileBlockService{
		db: db,
	}, nil
}

func (a *FileBlockService) GetNextBlocksFromHash(from models.Hash) ([]models.Block, error) {
	a.mu.Lock()

	blocks := make([]models.Block, 0)
	hasFoundHash := false

	//for each block found in database
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

//
//func (a *FileBlockService) AddNextBlocksFromHash(from models.Hash) ([]models.BlockDB, error) {
//	blocks := make([]models.BlockDB, 0)
//	hasFoundHash := false
//
//	//for each block found in database
//	scanner := bufio.NewScanner(a.db)
//	for scanner.Scan() {
//		if err := scanner.Err(); err != nil {
//			return blocks, err
//		}
//
//		blockFsJson := scanner.Bytes()
//		var blockDB models.BlockDB
//		err := json.Unmarshal(blockFsJson, &blockDB)
//		if err != nil {
//			return blocks, err
//		}
//
//		if hasFoundHash {
//			blocks = append(blocks, blockDB)
//			continue
//		}
//
//		if from == blockDB.Hash {
//			hasFoundHash = true
//		}
//	}
//
//	return blocks, nil
//}
