package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"
)

const AUTH_DOMAIN_URL = "/api/auth"

func RunDomain(r *gin.Engine,
	jwtService *services.JwtService,
	passwordService *services.PasswordService,
	userService *services.UserService,
	isActivateJwksEndpoint bool,
) {
	v1 := r.Group(AUTH_DOMAIN_URL)
	AuthRegister(v1.Group("/"), &AuthEnv{
		errorBuilder:           utils.NewErrorBuilder(),
		jwtService:             jwtService,
		userService:            userService,
		passwordService:        passwordService,
		isActivateJwksEndpoint: isActivateJwksEndpoint,
	})
}
