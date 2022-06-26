package balances

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/test"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"
)

var (
	testBlockchainFileDatabaseConf = test.NewTestBlockchainFileDatabaseConf(test.GenesisFilePath, test.BlocksFilePath)
	serviceState                   = services.NewFileStateService(testBlockchainFileDatabaseConf)
)

func setTestBlockchainFileDatabaseConf(genesisFilePath string, transactionFilePath string, isWrongGenesisFilePath bool, isWrongTransactionFilePath bool) {
	testBlockchainFileDatabaseConf.SetTestBlockchainFileDatabaseConf(genesisFilePath, transactionFilePath, isWrongGenesisFilePath, isWrongTransactionFilePath)
	serviceState = services.NewFileStateService(testBlockchainFileDatabaseConf)
}

type TestBalanceResponse struct {
	Response []BalanceResponse `json:"balances"`
}

var ListBalancesDomainTests = []struct {
	init           func(*http.Request)
	url            string
	method         string
	bodyData       []byte
	expectedCode   int
	jsonResponse   string
	validationFunc func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions)
	msg            string
	after          func(*http.Request)
}{
	//---------------------   Test suit for balance endpoints   ---------------------
	{
		init: func(req *http.Request) {
			setTestBlockchainFileDatabaseConf(test.GenesisFilePath, test.BlocksFilePath, false, false)
		},
		url:          BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		method:       "POST",
		expectedCode: http.StatusOK,
		jsonResponse: `{"balances":[{"account":"0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf","value":1003000},{"account":"0x7b65a12633dbe9a413b17db515732d69e684ebe2","value":998000}]}`,
		validationFunc: func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
			var balances TestBalanceResponse
			err := json.Unmarshal([]byte(wBodyA), &balances)
			if err != nil {
				fmt.Printf("%v", err)
			}
			sort.Slice(balances.Response, func(i, j int) bool {
				return balances.Response[i].Value < balances.Response[j].Value
			})

			var balances2 TestBalanceResponse
			err = json.Unmarshal([]byte(wBodyE), &balances2)
			if err != nil {
				asserts.Error(err, "%v")
			}
			sort.Slice(balances2.Response, func(i, j int) bool {
				return balances2.Response[i].Value < balances2.Response[j].Value
			})

			asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
			asserts.Equal(reflect.DeepEqual(balances, balances2), true, "Response Content - "+testName)
		},
		msg:   "request balances list should return a list of balances",
		after: func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			setTestBlockchainFileDatabaseConf(test.GenesisFilePath, test.BlocksFilePath, false, true)
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusInternalServerError,
		`{"error":{"code":500,"status":"Internal Server Error","message":"","context":[]}}`,
		test.StandardHttpValidationFunc,
		"request balances list with wrong genesis path should return code 500",
		func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			setTestBlockchainFileDatabaseConf(test.GenesisFilePath, test.BlocksFilePath, true, true)

		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusInternalServerError,
		`{"error":{"code":500,"status":"Internal Server Error","message":"","context":[]}}`,
		test.StandardHttpValidationFunc,
		"request balances list with wrong transaction path should return code 500",
		func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			setTestBlockchainFileDatabaseConf(test.GenesisFilePath, test.BlocksFilePath, true, true)
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusInternalServerError,
		`{"error":{"code":500,"status":"Internal Server Error","message":"","context":[]}}`,
		test.StandardHttpValidationFunc,
		"request balances list with state to nil should return code 500",
		func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			setTestBlockchainFileDatabaseConf(test.EmptyGenesisFilePath, test.EmptyBlocksFilePath, false, false)
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusNotFound,
		`{"error":{"code":404,"status":"Not Found","message":"balances could not be found","context":[]}}`,
		test.StandardHttpValidationFunc,
		"request balances list with wrong transaction path should return code 500",
		func(req *http.Request) {},
	},
}

func TestBalancesEnv_ListBalances(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()
	initServer(r)

	for _, testData := range ListBalancesDomainTests {
		var req *http.Request
		var err error
		if testData.bodyData != nil {
			req, err = http.NewRequest(testData.method, testData.url, bytes.NewBuffer(testData.bodyData))
		} else {
			req, err = http.NewRequest(testData.method, testData.url, nil)
		}
		req.Header.Set("Content-Type", "application/json")
		asserts.NoError(err)

		testData.init(req)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		testData.after(req)

		testData.validationFunc(testData.expectedCode, w.Code, testData.msg, testData.jsonResponse, w.Body.String(), asserts)
	}
}

func initServer(r *gin.Engine) {
	RunDomain(r, serviceState)
}
