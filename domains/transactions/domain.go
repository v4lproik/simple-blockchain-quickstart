package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"
)

const TRANSACTIONS_DOMAIN_URL = "/api/transactions"

func RunDomain(r *gin.Engine, state models.State, transactionService services.TransactionService, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(TRANSACTIONS_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	TransactionsRegister(v1.Group("/"), &TransactionsEnv{
		state:              state,
		transactionService: transactionService,
		errorBuilder:       utils.NewErrorBuilder(),
	})
}
