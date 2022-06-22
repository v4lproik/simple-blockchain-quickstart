package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const BALANCES_DOMAIN_URL = "/api/transactions"

func RunDomain(r *gin.Engine, stateService services.StateService, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(BALANCES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	TransactionsRegister(v1.Group("/"), &TransactionsEnv{
		stateService: stateService,
		errorBuilder: domains.NewErrorBuilder(),
	})
}
