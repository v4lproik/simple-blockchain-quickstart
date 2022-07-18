package test

import (
	"sync"

	"github.com/v4lproik/simple-blockchain-quickstart/domains/wallets"

	"github.com/stretchr/testify/assert"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/log"
)

var (
	// the init function is called in all test files, we don't need to reapply default configuration over and over
	setContextSafeGuard sync.Once

	// var which set the app context
	GenesisFilePath      = "../../test/testdata/genesis_test.json"
	EmptyGenesisFilePath = "../../test/testdata/genesis_empty.json"
	BlocksFilePath       = "../../test/testdata/blocks_test.db"
	EmptyBlocksFilePath  = "../../test/testdata/blocks_empty.db"
	KeystoreDirPath      = "../../test/testdata/keystore/"

	// functions that are used to verify whether a test is valid or not
	StandardHttpValidationFunc = func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
		asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
		asserts.Equal(wBodyE, wBodyA, "Response Content - "+testName)
	}
	RegexpHttpValidationFunc = func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
		asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
		asserts.Regexp(wBodyE, wBodyA, "Response Content - "+testName)
	}

	// services used across the entire application
	ErrorBuilder    common.ErrorBuilder
	KeyStoreService wallets.KeystoreService
)

func InitTestContext() {
	setContextSafeGuard.Do(func() {
		// init logger
		isProd := false
		// stdout
		logPath := ""
		log.InitLogger(isProd, logPath)

		// init services
		var err error
		ErrorBuilder = common.NewErrorBuilder()
		KeyStoreService, err = wallets.NewEthKeystore(KeystoreDirPath)
		if err != nil {
			log.Fatalf("InitTestContext: cannot init eth keystore service: %s", err)
		}
	})
}
