package conf

type FileDatabaseConf interface {
	GenesisFilePath() string
	TransactionFilePath() string
}

type BlockchainFileDatabaseConf struct {
	genesisFilePath     string
	transactionFilePath string
}

func NewBlockchainFileDatabaseConf(genesisFilePath string, transactionFilePath string) BlockchainFileDatabaseConf {
	return BlockchainFileDatabaseConf{
		genesisFilePath:     genesisFilePath,
		transactionFilePath: transactionFilePath,
	}
}

func (f BlockchainFileDatabaseConf) GenesisFilePath() string {
	return f.genesisFilePath
}

func (f BlockchainFileDatabaseConf) TransactionFilePath() string {
	return f.transactionFilePath
}
