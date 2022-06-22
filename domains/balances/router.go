package balances

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	"net/http"
)

const LIST_BALANCES_ENDPOINT = "/"

type BalancesEnv struct {
	bddConf      conf.FileDatabaseConf
	errorBuilder ErrorBuilder
}

func BalancesRegister(router *gin.RouterGroup, env *BalancesEnv) {
	router.POST(LIST_BALANCES_ENDPOINT, env.ListBalances)
}

func (env BalancesEnv) ListBalances(c *gin.Context) {
	state, err := models.NewStateFromFile(env.bddConf.GenesisFilePath(), env.bddConf.TransactionFilePath())
	if err != nil || state == nil {
		//TODO add type of error for NewStateFromFile
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}

	if len(state.Balances) == 0 {
		AbortWithError(c, *env.errorBuilder.New(404, "balances could not be found"))
		return
	}

	//map state with state response
	serializer := BalancesSerializer{state.Balances}

	//render
	c.JSON(http.StatusOK, gin.H{"balances": serializer.Response()})
	return
}
