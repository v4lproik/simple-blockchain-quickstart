package balances

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	"net/http"
)

type BalancesEnv struct {
	genesisFilePath     string
	transactionFilePath string
	errorBuilder        ErrorBuilder
}

func BalancesRegister(router *gin.RouterGroup, env *BalancesEnv) {
	router.POST("/", env.ListBalances)
}

func (env BalancesEnv) ListBalances(c *gin.Context) {
	state, err := models.NewStateFromFile(env.genesisFilePath, env.transactionFilePath)
	if err != nil {
		//TODO get which type of error happened and map it to http error for clarity
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}

	if state == nil {
		AbortWithError(c, *env.errorBuilder.New(404, "state could not be found"))
		return
	}

	balances := state.Balances
	if len(balances) == 0 {
		AbortWithError(c, *env.errorBuilder.New(404, "balances could not be found"))
		return
	}

	//map state with state response
	serializer := BalancesSerializer{balances}

	//render
	c.JSON(http.StatusOK, gin.H{"balances": serializer.Response()})
	return
}
