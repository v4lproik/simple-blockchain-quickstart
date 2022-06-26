package test

import (
	"github.com/stretchr/testify/assert"
)

var (
	//var which set the app context
	GenesisFilePath      = "../../databases/testdata/genesis_test.json"
	EmptyGenesisFilePath = "../../databases/testdata/genesis_empty.json"
	BlocksFilePath       = "../../databases/testdata/blocks_test.db"
	EmptyBlocksFilePath  = "../../databases/testdata/blocks_empty.db"
	KeystoreDirPath      = "../../databases/testdata/keystore/"

	//functions that are used to verify whether a test is valid or not
	StandardHttpValidationFunc = func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
		asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
		asserts.Equal(wBodyE, wBodyA, "Response Content - "+testName)
	}
	RegexpHttpValidationFunc = func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
		asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
		asserts.Regexp(wBodyE, wBodyA, "Response Content - "+testName)
	}
)

//function which help you set the app context
type TestBlockchainFileDatabaseConf struct {
	genesisFilePath            string
	transactionFilePath        string
	isWrongGenesisFilePath     bool
	isWrongTransactionFilePath bool
}

func NewTestBlockchainFileDatabaseConf(genesisFilePath, transactionFilePath string) *TestBlockchainFileDatabaseConf {
	return &TestBlockchainFileDatabaseConf{
		genesisFilePath:            genesisFilePath,
		transactionFilePath:        transactionFilePath,
		isWrongGenesisFilePath:     false,
		isWrongTransactionFilePath: false,
	}
}

func (t TestBlockchainFileDatabaseConf) GenesisFilePath() string {
	if t.isWrongGenesisFilePath {
		return t.genesisFilePath[:1]
	}
	return t.genesisFilePath
}

func (t TestBlockchainFileDatabaseConf) TransactionFilePath() string {
	if t.isWrongTransactionFilePath {
		return t.transactionFilePath[:1]
	}
	return t.transactionFilePath
}

func (t *TestBlockchainFileDatabaseConf) SetTestBlockchainFileDatabaseConf(genesisFilePath string, transactionFilePath string, isWrongGenesisFilePath bool, isWrongTransactionFilePath bool) {
	t.genesisFilePath = genesisFilePath
	t.transactionFilePath = transactionFilePath
	t.isWrongGenesisFilePath = isWrongGenesisFilePath
	t.isWrongTransactionFilePath = isWrongTransactionFilePath
}
