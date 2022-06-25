package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
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
	router.PUT(ADD_TRANSACTIONS_ENDPOINT, env.AddTransaction)
}

type AddTransactionParams struct {
	From   string        `json:"from" binding:"required,gte=2"`
	To     string        `json:"to" binding:"required,gte=2"`
	Value  uint          `json:"value" binding:"required"`
	Reason models.Reason `json:"reason" binding:"omitempty,enum"`
}

func (env TransactionsEnv) AddTransaction(c *gin.Context) {
	params := &AddTransactionParams{}
	//check params
	if err := ShouldBind(c, env.errorBuilder, "transaction cannot be added", params); err != nil {
		AbortWithError(c, *err)
		return
	}

	//create a transaction
	tx := models.NewTransaction(
		models.Account(params.From),
		models.Account(params.To),
		params.Value,
		string(params.Reason),
	)

	//get state
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

	//add to state
	hash, err := env.transactionService.AddTransaction(state, tx)
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(500, "transaction cannot be added", err))
		return
	}

	//map state with state response
	serializer := TransactionSerializer{*hash, *tx}

	//render
	c.JSON(http.StatusOK, gin.H{"transaction": serializer.Response()})
	return
}
