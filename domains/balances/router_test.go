package balances

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testBlockchainFileDatabaseConf = &TestBlockchainFileDatabaseConf{"../../databases/genesis.json", "../../databases/blocks.db", false, false}

var ListConnectorTests = []struct {
	init           func(*http.Request)
	url            string
	method         string
	bodyData       []byte
	expectedCode   int
	responseRegexg string
	msg            string
	after          func(*http.Request)
}{
	//---------------------   Test suit for balance endpoints   ---------------------
	{
		func(req *http.Request) {
			testBlockchainFileDatabaseConf.isWrongGenesisFilePath = false
			testBlockchainFileDatabaseConf.isWrongTransactionFilePath = false
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusOK,
		`{"balances":\[{"account":"v4lproik","value":998000},{"account":"cloudvenger","value":1003000}\]}`,
		"request balances list should return a list of balances",
		func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			testBlockchainFileDatabaseConf.isWrongGenesisFilePath = true
			testBlockchainFileDatabaseConf.isWrongTransactionFilePath = false
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusInternalServerError,
		`{"error":{"code":500,"status":"Internal Server Error","message":"","context":\[\]}}`,
		"request balances list with wrong genesis path should return code 500",
		func(req *http.Request) {},
	},
	{
		func(req *http.Request) {
			testBlockchainFileDatabaseConf.isWrongGenesisFilePath = true
			testBlockchainFileDatabaseConf.isWrongTransactionFilePath = true
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusInternalServerError,
		`{"error":{"code":500,"status":"Internal Server Error","message":"","context":\[\]}}`,
		"request balances list with wrong transaction path should return code 500",
		func(req *http.Request) {},
	},
}

func TestBalancesEnv_ListBalances(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()
	initServer(r)

	for _, testData := range ListConnectorTests {
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

		asserts.Equal(testData.expectedCode, w.Code, "Response Status - "+testData.msg)
		asserts.Regexp(testData.responseRegexg, w.Body.String(), "Response Content - "+testData.msg)
	}
}

func initServer(r *gin.Engine) {
	RunDomain(r, testBlockchainFileDatabaseConf)
}

//test models
type TestBlockchainFileDatabaseConf struct {
	genesisFilePath            string
	transactionFilePath        string
	isWrongGenesisFilePath     bool
	isWrongTransactionFilePath bool
}

func (f TestBlockchainFileDatabaseConf) GenesisFilePath() string {
	if f.isWrongGenesisFilePath {
		return f.genesisFilePath[:1]
	}
	return f.genesisFilePath
}

func (f TestBlockchainFileDatabaseConf) TransactionFilePath() string {
	if f.isWrongTransactionFilePath {
		return f.transactionFilePath[:1]
	}
	return f.transactionFilePath
}
