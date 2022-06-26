package wallets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const WALLETS_DOMAIN_URL = "/api/wallets"

func RunDomain(r *gin.Engine, keystoreDataDirPath string, middlewares ...gin.HandlerFunc) error {
	v1 := r.Group(WALLETS_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	keystoreService, err := NewKeystore(keystoreDataDirPath)
	if err != nil {
		return fmt.Errorf("couldn't initiate keystore service %v", err)
	}
	WalletsRegister(v1.Group("/"), &WalletsEnv{
		keystore:     keystoreService,
		errorBuilder: domains.NewErrorBuilder(),
	})

	return nil
}
