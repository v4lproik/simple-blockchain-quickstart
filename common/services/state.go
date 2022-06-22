package services

import (
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
)

type StateService interface {
	GetState() (models.State, error)
}

type FileStateService struct {
	fileDatabaseConf conf.FileDatabaseConf
}

func NewFileStateService(fileDatabaseConf conf.FileDatabaseConf) FileStateService {
	return FileStateService{
		fileDatabaseConf,
	}
}

func (a FileStateService) GetState() (models.State, error) {
	state, err := models.NewStateFromFile(a.fileDatabaseConf.GenesisFilePath(), a.fileDatabaseConf.TransactionFilePath())
	if err != nil {
		return nil, fmt.Errorf("cannot get blockchain state: %v", err)
	}

	return state, nil
}
