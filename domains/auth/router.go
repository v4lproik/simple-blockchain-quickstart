package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/gin-jwks-rsa"
	. "github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	log "go.uber.org/zap"
	"net/http"
)

const LOGIN_ENDPOINT = "/login"
const JWKS_ENDPOINT = "/.well-known/jwks.json"

type AuthEnv struct {
	errorBuilder           ErrorBuilder
	jwtService             *services.JwtService
	userService            *services.UserService
	passwordService        *services.PasswordService
	isActivateJwksEndpoint bool
}

func AuthRegister(router *gin.RouterGroup, env *AuthEnv) {
	router.POST(LOGIN_ENDPOINT, env.Login)

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

type LoginParams struct {
	Username string `json:"username" binding:"required,gte=2"`
	Password string `json:"password" binding:"required,password"`
}

func (env AuthEnv) Login(c *gin.Context) {
	params := &LoginParams{}
	//check params
	if err := ShouldBind(c, env.errorBuilder, "login cannot occur", params); err != nil {
		AbortWithError(c, *err)
		return
	}

	//check if user is in bdd
	user, err := env.userService.Get(params.Username)
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(404, "user %s could not be found", params.Username))
		return
	}

	//check if passwords match
	_, err = env.passwordService.ComparePasswordAndHash(params.Password, user.Hash)
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(401, "password is not correct"))
		return
	}

	//create and sign an access token if passwords match
	token, err := env.jwtService.SignToken(*user)
	if err != nil {
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}

	//render
	c.JSON(http.StatusOK, gin.H{"access_token": token})
	return
}
