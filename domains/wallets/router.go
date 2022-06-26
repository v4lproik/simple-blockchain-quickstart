package wallets

import (
	"github.com/gin-gonic/gin"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	log "go.uber.org/zap"
	"net/http"
)

const CREATE_WALLET_ENDPOINT = "/"

type WalletsEnv struct {
	keystore     *KeystoreService
	errorBuilder ErrorBuilder
}

func WalletsRegister(router *gin.RouterGroup, env *WalletsEnv) {
	router.PUT(CREATE_WALLET_ENDPOINT, env.CreateWallet)
}

type CreateWalletParams struct {
	Password string `json:"password" binding:"required,password"`
}

func (env WalletsEnv) CreateWallet(c *gin.Context) {
	params := &CreateWalletParams{}
	//check params
	if err := ShouldBind(c, env.errorBuilder, "wallet cannot be created", params); err != nil {
		AbortWithError(c, *err)
		return
	}

	acc, err := env.keystore.NewKeystoreAccount(params.Password)
	if err != nil {
		log.S().Errorf("cannot generate a new wallet account %v", err)
		AbortWithError(c, *env.errorBuilder.New(http.StatusInternalServerError, "cannot generate a new wallet account"))
		return
	}

	//render
	c.JSON(http.StatusOK, gin.H{"balances": acc})
	return
}
