package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/middleware"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/auth"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/balances"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/healthz"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/transactions"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/wallets"
	"github.com/v4lproik/simple-blockchain-quickstart/utils"
	"go.uber.org/zap"
	"time"
)

var (
	log     *zap.SugaredLogger
	apiConf = utils.ApiConf{}
)

func runHttpServer() {
	if err := env.Parse(&apiConf); err != nil {
		logger.Fatal(err)
	}

	//init router
	serverCorsOpts := apiConf.Server.HttpCors
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowMethods:  serverCorsOpts.AllowedMethods,
		AllowHeaders:  serverCorsOpts.AllowedHeaders,
		AllowOrigins:  serverCorsOpts.AllowedOrigins,
		ExposeHeaders: []string{"Content-Length"},
	}))

	//extend logger to handlers exposed by Gin
	r.Use(ginzap.Ginzap(logger.Desugar(), time.RFC3339, true))

	//logs all panic to error log
	//stack opt means whether to append to the output the stack info
	r.Use(ginzap.RecoveryWithZap(logger.Desugar(), true))
	r.Use(gin.Recovery())

	//start the functional domains
	bindFunctionalDomains(r)

	//start server according to the configuration passed in parameter or env variables
	serverOpts := apiConf.Server.Options
	serverPort := apiConf.Server.Port
	serverAddress := apiConf.Server.Address
	if serverOpts.IsSsl {
		logger.Info("start server with tls")
		r.RunTLS(fmt.Sprintf(":%d", serverPort), serverOpts.CertFile, serverOpts.KeyFile)
	} else {
		logger.Info("start server without tls")
		r.Run(fmt.Sprintf("%s:%d", serverAddress, serverPort))
	}
}

//TODO: Enumerate which domains need to start at bootstrap
func bindFunctionalDomains(r *gin.Engine) {
	//initiate services
	errorBuilder := common.NewErrorBuilder()
	fileStateService := services.NewFileStateService(conf.NewBlockchainFileDatabaseConf(opts.GenesisFilePath, opts.TransactionsFilePath))
	fileTransactionService := services.NewFileTransactionService()
	keystoreService, err := wallets.NewEthKeystore(opts.KeystoreDirPath)
	if err != nil {
		log.Fatalf("cannot create keystore service %v", err)
	}
	jwtOpts := apiConf.Auth.Jwt
	jwtService, err := services.NewJwtService(
		services.NewVerifyingConf(
			jwtOpts.Verifying.JkmsUrl,
			jwtOpts.Verifying.JkmsRefreshCacheIntervalInMin,
			jwtOpts.Verifying.JkmsRefreshCacheRateLimitInMin,
			jwtOpts.Verifying.JkmsRefreshCacheTimeoutInSec,
		),
		services.NewSigningConf(
			jwtOpts.Signing.Algo,
			jwtOpts.Signing.Audience,
			jwtOpts.Signing.Domain,
			jwtOpts.Signing.ExpiresIn,
			jwtOpts.Signing.Issuer,
			jwtOpts.Signing.KeyPath,
			jwtOpts.Signing.KeyId,
		),
	)
	if err != nil {
		log.Fatalf("cannot create jwt service %v", err)
	}
	passwordService := services.NewDefaultPasswordService()
	userService, err := services.NewUserService(opts.UsersFilePath)
	if err != nil {
		log.Fatalf("cannot create user service %v", err)
	}
	//initiate middlewares
	auto401 := apiConf.Auth.IsAuthenticationActivated
	authMiddleware := middleware.AuthWebSessionMiddleware(auto401, errorBuilder, jwtService)

	//run domains
	healthz.RunDomain(r)
	balances.RunDomain(r, fileStateService, authMiddleware)
	transactions.RunDomain(r, fileStateService, fileTransactionService, authMiddleware)
	auth.RunDomain(r, jwtService, &passwordService, userService, apiConf.Auth.IsJwksEndpointActivated)
	wallets.RunDomain(r, &wallets.WalletsEnv{
		Keystore:     keystoreService,
		ErrorBuilder: errorBuilder,
	}, authMiddleware)
}
