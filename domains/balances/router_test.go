package balances

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/test"
)

var (
	state, _   = models.NewStateFromFile(test.GenesisFilePath, test.BlocksFilePath)
	tState     = &testState{}
	balanceEnv *BalancesEnv
)

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
	//{
	//	init: func(req *http.Request) {
	//		setBalanceEnv(state, test.ErrorBuilder)
	//	},
	//	url:          BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
	//	method:       "POST",
	//	expectedCode: http.StatusOK,
	//	jsonResponse: `{"balances":[{"account":"0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf","value":998923},{"account":"0x7b65a12633dbe9a413b17db515732d69e684ebe2","value":1001077}]}`,
	//	validationFunc: func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions) {
	//		var balances TestBalanceResponse
	//		err := json.Unmarshal([]byte(wBodyA), &balances)
	//		if err != nil {
	//			fmt.Printf("%v", err)
	//		}
	//		sort.Slice(balances.Response, func(i, j int) bool {
	//			return balances.Response[i].Value < balances.Response[j].Value
	//		})
	//
	//		var balances2 TestBalanceResponse
	//		err = json.Unmarshal([]byte(wBodyE), &balances2)
	//		if err != nil {
	//			asserts.Error(err, "%v")
	//		}
	//		sort.Slice(balances2.Response, func(i, j int) bool {
	//			return balances2.Response[i].Value < balances2.Response[j].Value
	//		})
	//
	//		asserts.Equal(wCodeE, wCodeA, "Response Status - "+testName)
	//		asserts.Equal(reflect.DeepEqual(balances, balances2), true, "Response Content - "+testName)
	//	},
	//	msg:   "request balances list should return a list of balances",
	//	after: func(req *http.Request) {},
	//},
	{
		func(req *http.Request) {
			balanceEnv.state = tState
		},
		BALANCES_DOMAIN_URL + LIST_BALANCES_ENDPOINT,
		"POST",
		nil,
		http.StatusNotFound,
		`{"error":{"code":404,"status":"Not Found","message":"balances could not be found","context":[]}}`,
		test.StandardHttpValidationFunc,
		"request balances list with balances=0 should return error not found",
		func(req *http.Request) {},
	},
}

func TestBalancesEnv_ListBalances(t *testing.T) {
	test.InitTestContext()
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
	if balanceEnv == nil {
		balanceEnv = NewBalancesEnv(state, test.ErrorBuilder)
	}
	RunDomain(r, balanceEnv)
}

// utils
type testState struct{}

func (t testState) Add(transaction models.Transaction) error {
	// TODO implement me
	panic("implement me")
}

func (t testState) AddBlock(block models.Block) error {
	// TODO implement me
	panic("implement me")
}

func (t testState) AddBlocks(blocks []models.Block) error {
	// TODO implement me
	panic("implement me")
}

func (t testState) Balances() map[models.Account]uint {
	return make(map[models.Account]uint, 0)
}

func (t testState) Persist() (models.Hash, error) {
	// TODO implement me
	panic("implement me")
}

func (t testState) Close() error {
	// TODO implement me
	panic("implement me")
}

func (t testState) GetLatestBlockHash() models.Hash {
	// TODO implement me
	panic("implement me")
}

func (t testState) GetLatestBlockHeight() uint64 {
	// TODO implement me
	panic("implement me")
}

func (t testState) Print() {
	// TODO implement me
	panic("implement me")
}
