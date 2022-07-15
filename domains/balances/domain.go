package balances

import (
	"github.com/gin-gonic/gin"
)

const BALANCES_DOMAIN_URL = "/api/balances"

func RunDomain(r *gin.Engine, balancesEnv *BalancesEnv, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(BALANCES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	BalancesRegister(v1.Group("/"), balancesEnv)
}
