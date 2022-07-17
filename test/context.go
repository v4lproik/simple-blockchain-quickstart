package test

import (
	"sync"

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
	ErrorBuilder = common.NewErrorBuilder()
)

func InitTestContext() {
	setContextSafeGuard.Do(func() {
		// test env, not prod
		isProd := false
		// stdout
		logPath := ""
		log.InitLogger(isProd, logPath)
	})
}
