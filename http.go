package main

import (
	"fmt"
	"time"

	parser "github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/middleware"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/auth"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/balances"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/healthz"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/nodes"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/transactions"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/wallets"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
	"github.com/v4lproik/simple-blockchain-quickstart/utils"
)

type Domain string

const (
	AUTH         Domain = "AUTH"
	BALANCES            = "BALANCES"
	HEALTHZ             = "HEALTHZ"
	NODES               = "NODES"
	TRANSACTIONS        = "TRANSACTIONS"
	WALLETS             = "WALLETS"
)

var apiConf = utils.ApiConf{}

func runHttpServer() {
	if err := parser.Parse(&apiConf); err != nil {
		Logger.Fatal(err)
	}

	// init router
	serverCorsOpts := apiConf.Server.HttpCors
	if env.isProd() {
		gin.SetMode("release")
	}
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(cors.New(cors.Config{
		AllowMethods:  serverCorsOpts.AllowedMethods,
		AllowHeaders:  serverCorsOpts.AllowedHeaders,
		AllowOrigins:  serverCorsOpts.AllowedOrigins,
		ExposeHeaders: []string{"Content-Length"},
	}))

	// extend Logger to handlers exposed by Gin
	r.Use(ginzap.Ginzap(Logger.Desugar(), time.RFC3339, true))

	// logs all panic to error log
	// stack opt means whether to append to the output the stack info
	r.Use(ginzap.RecoveryWithZap(Logger.Desugar(), true))
	r.Use(gin.Recovery())

	// start the functional domains
	bindFunctionalDomains(r)

	// start server according to the configuration passed in parameter or EnvVal variables
	serverOpts := apiConf.Server.Options
	serverPort := apiConf.Server.Port
	serverAddress := apiConf.Server.Address
	if serverOpts.IsSsl {
		Logger.Info("runHttpServer: start server with tls")
		r.RunTLS(fmt.Sprintf(":%d", serverPort), serverOpts.CertFile, serverOpts.KeyFile)
	} else {
		Logger.Info("runHttpServer: start server without tls")
		r.Run(fmt.Sprintf("%s:%d", serverAddress, serverPort))
	}
}

func bindFunctionalDomains(r *gin.Engine) {
	// TODO: extract business logic and put it in a state service
	state, err := models.NewStateFromFile(opts.GenesisFilePath, opts.TransactionsFilePath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot initialise the state: %s", err)
	}
	// initiate services
	errorBuilder := common.NewErrorBuilder()
	fileTransactionService := services.NewFileTransactionService()
	keystoreService, err := wallets.NewEthKeystore(opts.KeystoreDirPath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot create keystore service: %s", err)
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
		Logger.Fatalf("bindFunctionalDomains: cannot create jwt service: %s", err)
	}
	passwordService := services.NewDefaultPasswordService()
	userService, err := services.NewUserService(opts.UsersFilePath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot create user service: %s", err)
	}

	nodeService, err := nodes.NewNodeService(opts.NodesFilePath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot create node service: %s", err)
	}

	blockService, err := services.NewFileBlockService(opts.TransactionsFilePath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot create block service: %s", err)
	}

	// initiate middlewares
	auto401 := apiConf.Auth.IsAuthenticationActivated
	authMiddleware := middleware.AuthWebSessionMiddleware(auto401, errorBuilder, jwtService)

	// run domains
	for _, domain := range apiConf.Domains.ToStart {
		switch Domain(domain) {
		case AUTH:
			auth.RunDomain(r, jwtService, &passwordService, userService, apiConf.Auth.IsJwksEndpointActivated)
		case BALANCES:
			balances.RunDomain(r, balances.NewBalancesEnv(state, errorBuilder), authMiddleware)
		case HEALTHZ:
			healthz.RunDomain(r)
		case NODES:
			nodes.RunDomain(r, nodeService, state, blockService)
		case TRANSACTIONS:
			transactions.RunDomain(r, state, fileTransactionService, authMiddleware)
		case WALLETS:
			wallets.RunDomain(r, &wallets.WalletsEnv{
				Keystore:     keystoreService,
				ErrorBuilder: errorBuilder,
			}, authMiddleware)
		default:
			Logger.Fatalf("bindFunctionalDomains: the functional domain %s is unknown", domain)
		}
	}
}
