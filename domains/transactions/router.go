package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	"net/http"
)

const ADD_TRANSACTIONS_ENDPOINT = "/"

type TransactionsEnv struct {
	stateService       services.StateService
	transactionService services.TransactionService
	errorBuilder       ErrorBuilder
}

func TransactionsRegister(router *gin.RouterGroup, env *TransactionsEnv) {
	router.POST(ADD_TRANSACTIONS_ENDPOINT, env.AddTransaction)
}

func (env TransactionsEnv) AddTransaction(c *gin.Context) {
	state, err := env.stateService.GetState()
	if err != nil || state == nil {
		//TODO add type of error for NewStateFromFile
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}

	if len(state.Balances()) == 0 {
		AbortWithError(c, *env.errorBuilder.New(404, "balances could not be found"))
		return
	}

	env.transactionService.AddTransaction(nil, nil)

	//map state with state response
	serializer := BalancesSerializer{state.Balances()}

	//render
	c.JSON(http.StatusOK, gin.H{"balances": serializer.Response()})
	return
}
