package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const AUTH_DOMAIN_URL = "/api/auth"

func RunDomain(r *gin.Engine, jwtService *services.JwtService, isActivateJwksEndpoint bool) {
	v1 := r.Group(AUTH_DOMAIN_URL)
	AuthRegister(v1.Group("/"), &AuthEnv{
		errorBuilder:           domains.NewErrorBuilder(),
		jwtService:             jwtService,
		isActivateJwksEndpoint: isActivateJwksEndpoint,
	})
}
