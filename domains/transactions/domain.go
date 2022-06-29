package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
)

const BALANCES_DOMAIN_URL = "/api/transactions"

func RunDomain(r *gin.Engine, stateService services.StateService, transactionService services.TransactionService, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(BALANCES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	TransactionsRegister(v1.Group("/"), &TransactionsEnv{
		stateService:       stateService,
		transactionService: transactionService,
		errorBuilder:       common.NewErrorBuilder(),
	})
}
