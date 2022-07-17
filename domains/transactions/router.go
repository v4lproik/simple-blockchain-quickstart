package transactions

import (
	"github.com/gin-gonic/gin"
	. "github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	"net/http"
)

const ADD_TRANSACTIONS_ENDPOINT = "/"

type TransactionsEnv struct {
	state              models.State
	transactionService services.TransactionService
	errorBuilder       ErrorBuilder
}

func TransactionsRegister(router *gin.RouterGroup, env *TransactionsEnv) {
	router.PUT(ADD_TRANSACTIONS_ENDPOINT, env.AddTransaction)
}

type AddTransactionParams struct {
	From   string        `json:"from" binding:"required,account"`
	To     string        `json:"to" binding:"required,account"`
	Value  uint          `json:"value" binding:"required,gte=1"`
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
	from, _ := models.NewAccount(params.From)
	to, _ := models.NewAccount(params.To)
	tx := models.NewTransaction(
		from,
		to,
		params.Value,
		string(params.Reason),
	)

	state := env.state
	if len(state.Balances()) == 0 {
		AbortWithError(c, *env.errorBuilder.New(http.StatusNotFound, "balances could not be found"))
		return
	}

	//add to state
	hash, err := env.transactionService.AddTransaction(state, tx)
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(http.StatusInternalServerError, "transaction cannot be added", err))
		return
	}

	//map state with state response
	serializer := TransactionSerializer{*hash, *tx}

	//render
	c.JSON(http.StatusOK, gin.H{"transaction": serializer.Response()})
}
