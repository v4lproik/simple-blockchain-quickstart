package wallets

import (
	"github.com/gin-gonic/gin"
	. "github.com/v4lproik/simple-blockchain-quickstart/common"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
	"net/http"
)

type WalletsEnv struct {
	Keystore     KeystoreService
	ErrorBuilder ErrorBuilder
}

const CREATE_WALLET_ACC_ENDPOINT = "/"

func WalletsRegister(router *gin.RouterGroup, env *WalletsEnv) {
	router.PUT(CREATE_WALLET_ACC_ENDPOINT, env.CreateWallet)
}

type CreateWalletParams struct {
	Password string `json:"password" binding:"required,password"`
}

func (env *WalletsEnv) CreateWallet(c *gin.Context) {
	params := &CreateWalletParams{}
	errMsg := "wallet cannot be created"

	//check params
	if err := ShouldBind(c, env.ErrorBuilder, errMsg, params); err != nil {
		AbortWithError(c, *err)
		return
	}

	acc, err := env.Keystore.NewKeystoreAccount(params.Password)
	if err != nil {
		Logger.Errorf("CreateWallet: failed to generate a new wallet account: %s", err)
		AbortWithError(c, *env.ErrorBuilder.New(http.StatusInternalServerError, errMsg, err))
		return
	}

	//render
	c.JSON(http.StatusCreated, gin.H{"wallet": WalletSerializer{acc}.Response()})
	return
}
