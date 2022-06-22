package balances

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const BALANCES_DOMAIN_URL = "/api/balances"

func RunDomain(r *gin.Engine, bddConf conf.FileDatabaseConf, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(BALANCES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	BalancesRegister(v1.Group("/"), &BalancesEnv{
		bddConf:      bddConf,
		errorBuilder: domains.NewErrorBuilder(),
	})
}
