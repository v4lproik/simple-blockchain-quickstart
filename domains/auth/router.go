package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/gin-jwks-rsa"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	log "go.uber.org/zap"
	"net/http"
)

const LOGIN_ENDPOINT = "/login"
const VALIDATE_ENDPOINT = "/validate"
const JWKS_ENDPOINT = "/.well-known/jwks.json"

type AuthEnv struct {
	errorBuilder           ErrorBuilder
	jwtService             *services.JwtService
	isActivateJwksEndpoint bool
}

func AuthRegister(router *gin.RouterGroup, env *AuthEnv) {
	router.GET(LOGIN_ENDPOINT, env.Login)
	router.GET(VALIDATE_ENDPOINT, env.Validate)

	jwksConf, err := gin_jwks_rsa.NewConfigBuilder().
		ImportPrivateKey().
		WithKeyId(env.jwtService.PrivateKeyId()).
		WithPath(env.jwtService.PrivateKeyPath()).
		Build()
	if err != nil {
		log.S().Fatalf("cannot start jwks service %v", err)
	}
	router.GET(JWKS_ENDPOINT, gin_jwks_rsa.Jkws(*jwksConf))
}

func (env AuthEnv) Login(c *gin.Context) {
	signedToken, err := env.jwtService.SignToken("test")
	if err != nil {
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}
	//render
	c.JSON(http.StatusOK, gin.H{"auth": signedToken})
	return
}

func (env AuthEnv) Validate(c *gin.Context) {
	token, err := env.jwtService.VerifyToken("test")
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(401, "token is not valid"))
		return
	}

	//render
	c.JSON(http.StatusOK, gin.H{"auth": token})
	return
}
