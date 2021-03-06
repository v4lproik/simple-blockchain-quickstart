package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	parser "github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
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
)

type Domain string

const (
	AUTH         Domain = "AUTH"
	BALANCES     Domain = "BALANCES"
	HEALTHZ      Domain = "HEALTHZ"
	NODES        Domain = "NODES"
	TRANSACTIONS Domain = "TRANSACTIONS"
	WALLETS      Domain = "WALLETS"
)

var apiConf = ApiConf{}

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

	serverOpts := apiConf.Server.Options
	server := &http.Server{
		Addr:    net.JoinHostPort(apiConf.Server.Address, fmt.Sprintf("%d", apiConf.Server.Port)),
		Handler: r.Handler(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if serverOpts.IsSsl {
			Logger.Info("runHttpServer: start server with tls")
			server.ListenAndServeTLS(serverOpts.CertFile, serverOpts.KeyFile)
		} else {
			Logger.Info("runHttpServer: start server without tls")
			server.ListenAndServe()
		}
	}()

	gracefullyShutdownServer(ctx, server)
}

func bindFunctionalDomains(r *gin.Engine) {
	// TODO: extract business logic and put it in a state service
	state, err := models.NewStateFromFile(opts.GenesisFilePath, opts.TransactionsFilePath)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot initialise the state: %s", err)
	}
	// initiate services
	fileTransactionService := services.NewFileTransactionService()
	keystoreService, err := services.NewEthKeystore(opts.KeystoreDirPath)
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

	miningAccount, _ := models.NewAccount(opts.MinerAddress)
	blockService, err := services.NewFileBlockService(
		opts.TransactionsFilePath,
		apiConf.Consensus.Complexity,
		miningAccount,
	)
	if err != nil {
		Logger.Fatalf("bindFunctionalDomains: cannot create block service: %s", err)
	}

	// initiate middlewares
	auto401 := apiConf.Auth.IsAuthenticationActivated
	authMiddleware := middleware.AuthWebSessionMiddleware(auto401, jwtService)

	// run domains
	for _, domain := range apiConf.Domains.ToStart {
		switch Domain(domain) {
		case AUTH:
			auth.RunDomain(r, jwtService, &passwordService, userService, apiConf.Auth.IsJwksEndpointActivated)
		case BALANCES:
			balances.RunDomain(r, balances.NewBalancesEnv(state), authMiddleware)
		case HEALTHZ:
			healthz.RunDomain(r)
		case NODES:
			if err := nodes.RunDomain(
				r,
				nodeService,
				state,
				fileTransactionService,
				blockService,
				apiConf.Synchronisation.RefreshIntervalInSeconds,
				apiConf.Consensus.CreateNewBlockIntervalInSeconds,
			); err != nil {
				Logger.Fatalf("bindFunctionalDomains: cannot start the node domain: %w", err)
			}
		case TRANSACTIONS:
			transactions.RunDomain(r, state, fileTransactionService, authMiddleware)
		case WALLETS:
			wallets.RunDomain(r, &wallets.WalletsEnv{
				Keystore: keystoreService,
			}, authMiddleware)
		default:
			Logger.Fatalf("bindFunctionalDomains: the functional domain %s is unknown", domain)
		}
	}
}

func gracefullyShutdownServer(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	Logger.Infof("runHttpServer: trying to gracefully close http server...")

	timeoutShutdown := 2
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutShutdown)*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		Logger.Fatalf("runHttpServer: couldn't gracefully stop http server: %s", err)
	}

	<-ctx.Done()
	Logger.Info("runHttpServer: http server closed")
}
