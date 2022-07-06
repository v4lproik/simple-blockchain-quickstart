package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"os"
)

type BlockService interface {
	GetBlockNextBlocksFrom(hash models.Hash) ([]models.BlockDB, error)
}

type FileBlockService struct {
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

func (a *FileBlockService) GetBlockNextBlocksFrom(from models.Hash) ([]models.BlockDB, error) {
	blocks := make([]models.BlockDB, 0)
	hasFoundHash := false

	//for each block found in database
	scanner := bufio.NewScanner(a.db)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return blocks, err
		}

		blockFsJson := scanner.Bytes()
		var blockDB models.BlockDB
		err := json.Unmarshal(blockFsJson, &blockDB)
		if err != nil {
			return blocks, err
		}

		if hasFoundHash {
			blocks = append(blocks, blockDB)
			continue
		}

		if from == blockDB.Hash {
			hasFoundHash = true
		}
	}

	return blocks, nil
}
