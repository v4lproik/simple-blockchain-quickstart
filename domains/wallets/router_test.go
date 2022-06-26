package wallets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
	"github.com/v4lproik/simple-blockchain-quickstart/test"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	keystoreService KeystoreService
	walletsEnv      *WalletsEnv
	faultyKeystore  *FaultyKeystore
)

func setKeystoreService(f KeystoreService) {
	walletsEnv.Keystore = f
}

var CreateWalletsDomainTests = []struct {
	init           func(*http.Request)
	url            string
	method         string
	bodyData       func() ([]byte, error)
	expectedCode   int
	jsonResponse   string
	validationFunc func(wCodeE int, wCodeA int, testName string, wBodyE string, wBodyA string, asserts *assert.Assertions)
	msg            string
	after          func(*http.Request)
}{
	//---------------------   Test suit for wallet endpoints   ---------------------
	{
		init: func(req *http.Request) {
			setKeystoreService(keystoreService)
		},
		url:            WALLETS_DOMAIN_URL + CREATE_WALLET_ACC_ENDPOINT,
		method:         "PUT",
		bodyData:       func() ([]byte, error) { return json.Marshal(CreateWalletParams{Password: "P@assword123!"}) },
		expectedCode:   http.StatusCreated,
		jsonResponse:   `{"wallet":{"account":"0x[a-fA-F0-9]{40}"}}`,
		validationFunc: test.RegexpHttpValidationFunc,
		msg:            "request creation wallet with accepted password should return wallet account",
		after:          func(req *http.Request) {},
	},
	{
		init: func(req *http.Request) {
			setKeystoreService(keystoreService)
		},
		url:            WALLETS_DOMAIN_URL + CREATE_WALLET_ACC_ENDPOINT,
		method:         "PUT",
		bodyData:       func() ([]byte, error) { return json.Marshal(CreateWalletParams{Password: "aaaaaaaaaa!"}) },
		expectedCode:   http.StatusBadRequest,
		jsonResponse:   `{"error":{"code":400,"status":"Bad Request","message":"wallet cannot be created","context":[[{"field":"Password","message":"The password doesn't comply with the policy (min 8 char with min 1 upper, 1 number and 1 symbol)"}]]}}`,
		validationFunc: test.StandardHttpValidationFunc,
		msg:            "request creation wallet with not valid password should return error",
		after:          func(req *http.Request) {},
	},
	{
		init: func(req *http.Request) {
			setKeystoreService(faultyKeystore)
		},
		url:            WALLETS_DOMAIN_URL + CREATE_WALLET_ACC_ENDPOINT,
		method:         "PUT",
		bodyData:       func() ([]byte, error) { return json.Marshal(CreateWalletParams{Password: "P@assword123!"}) },
		expectedCode:   http.StatusInternalServerError,
		jsonResponse:   `{"error":{"code":500,"status":"Internal Server Error","message":"cannot generate a new wallet account","context":["cannot generate private key"]}}`,
		validationFunc: test.StandardHttpValidationFunc,
		msg:            "request creation wallet with error generating account should return error",
		after:          func(req *http.Request) {},
	},
}

func TestBalancesEnv_ListBalances(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()
	initServer(r)

	for _, testData := range CreateWalletsDomainTests {
		var req *http.Request
		var err error
		if testData.bodyData != nil {
			jsonPayload, err := testData.bodyData()
			asserts.NoError(err)
			req, err = http.NewRequest(testData.method, testData.url, bytes.NewBuffer(jsonPayload))
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
	services.ValidatorService{}.AddValidators()
	faultyKeystore = NewFaultyKeystore()
	keystoreService, _ = NewEthKeystore(test.KeystoreDirPath)
	walletsEnv = &WalletsEnv{
		Keystore:     keystoreService,
		ErrorBuilder: domains.NewErrorBuilder(),
	}

	RunDomain(r, walletsEnv)
}

//test models
type FaultyKeystore struct{}

func NewFaultyKeystore() *FaultyKeystore {
	return &FaultyKeystore{}
}

func (f *FaultyKeystore) NewKeystoreAccount(password string) (common.Address, error) {
	return common.Address{}, fmt.Errorf("cannot generate private key")
}