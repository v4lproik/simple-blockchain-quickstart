package balances

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

func RunDomain(r *gin.Engine, genesisFilePath string, transactionFilePath string, middlewares ...gin.HandlerFunc) {
	v1 := r.Group("/api/balances")
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	BalancesRegister(v1.Group("/"), &BalancesEnv{
		genesisFilePath:     genesisFilePath,
		transactionFilePath: transactionFilePath,
		errorBuilder:        domains.NewErrorBuilder(),
	})
}
