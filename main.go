package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/balances"
	"github.com/v4lproik/simple-blockchain-quickstart/domains/healthz"
	"github.com/v4lproik/simple-blockchain-quickstart/utils"
	"time"
)

var (
	apiConf = utils.ApiConf{}
)

func main() {
	//parse cli arguments
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)

	//parse command line arguments which set up the application
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	//init logger
	InitLogger(opts.LogFilePath)
	defer logger.Sync()

	//display program configuration
	displayAppConfiguration()

	//run as a node
	if opts.RunAsHttpserver {
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
		//TODO: Enumerate which domains need to start at bootstrap
		healthz.RunDomain(r)
		balances.RunDomain(r, conf.NewBlockchainFileDatabaseConf(opts.GenesisFilePath, opts.TransactionsFilePath))

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
	} else {
		//run as client
		//add commands and subcommands. eg. ./bin transaction list
		err = addCommands(parser)
		if err != nil {
			panic(err)
		}

		//parse command line for the requested actions (eg. transaction list, transaction add, etc...)
		_, err = parser.Parse()
		if err != nil {
			panic(err)
		}
	}
}
