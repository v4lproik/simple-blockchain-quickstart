package services

import (
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
)

type StateService interface {
	GetState() (error, models.State)
}

type FileStateService struct {
	fileDatabaseConf conf.FileDatabaseConf
}

func NewFileStateService(fileDatabaseConf conf.FileDatabaseConf) FileStateService {
	return FileStateService{
		fileDatabaseConf,
	}
}

func (a FileStateService) GetState() (error, models.State) {
	state, err := models.NewStateFromFile(a.fileDatabaseConf.GenesisFilePath(), a.fileDatabaseConf.TransactionFilePath())
	if err != nil {
		return fmt.Errorf("cannot get blockchain state: %v", err), nil
	}
	defer state.Close()

	return nil, state
}
