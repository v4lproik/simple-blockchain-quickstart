package transactions

import (
	"errors"
	"net/http"

	. "github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const ADD_TRANSACTIONS_ENDPOINT = "/"

type TransactionsEnv struct {
	state              models.State
	transactionService services.TransactionService
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
	// check params
	if err := ShouldBind(c, "transaction cannot be added", params); err != nil {
		AbortWithError(c, err)
		return
	}

	// check if the node is available to accept new transactions
	// the buffered channel is limited to a certain amount of bytes
	newPendingTxs := env.transactionService.NewPendingTxs()
	if len(newPendingTxs) == cap(newPendingTxs) {
		AbortWithError(c, NewError(http.StatusServiceUnavailable, "transaction cannot be added"))
		return
	}

	// create a transaction
	from, _ := models.NewAccount(params.From)
	to, _ := models.NewAccount(params.To)
	tx := models.NewTransaction(
		from,
		to,
		params.Value,
		string(params.Reason),
		DefaultTimeService.UnixUint64(),
	)

	state := env.state
	if len(state.Balances()) == 0 {
		AbortWithError(c, NewError(http.StatusNotFound, "balances could not be found"))
		return
	}

	// add to state
	err := env.transactionService.AddTx(tx)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.As(err, &services.ErrMarshalTx) {
			code = http.StatusConflict
		}
		AbortWithError(c, NewError(code, "transaction cannot be added"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}
