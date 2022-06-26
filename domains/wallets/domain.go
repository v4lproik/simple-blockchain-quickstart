package wallets

import (
	"github.com/gin-gonic/gin"
)

const WALLETS_DOMAIN_URL = "/api/wallets"

func RunDomain(r *gin.Engine, walletsEnv *WalletsEnv, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(WALLETS_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	WalletsRegister(v1.Group("/"), walletsEnv)
}
