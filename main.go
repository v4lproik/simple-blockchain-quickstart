package main

import (
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

var env EnvVal

// @title Simple Blockchain Quickstart
// @version 1.0
// @description About
// This is an experimental repository which aims at shipping a decent skeleton for anyone who wants to get into the blockchain world via the Golang programming language

// @contact.name API Support
// @contact.email rousseau.joel@gmail.com

// @host localhost:8080
func main() {
	// parse cli arguments
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)

	// parse command line arguments which set up the application
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	// check environment
	checkArgs()

	// init Logger
	Logger.InitLogger(env.isProd(), opts.LogFilePath)
	defer Logger.Sync()

	// display program configuration
	displayAppConfiguration()

	// run as a node
	if opts.RunAsHttpserver {
		// add specific json validators used by endpoints
		validatorService := new(services.ValidatorService)
		validatorService.AddValidators()

		// run the http server
		runHttpServer()
	} else {
		// run as client
		// add commands and subcommands. eg. ./bin transaction list
		err = addCommands(parser)
		if err != nil {
			panic(err)
		}

		// parse command line for the requested actions (eg. transaction list, transaction add, etc...)
		_, err = parser.Parse()
		if err != nil {
			panic(err)
		}
	}
}

func funcName(arr []string, done chan<- bool, c chan<- map[string]int, wg *sync.WaitGroup) {
	for i, add := range arr {
		go func(i int, add string) {
			defer wg.Done()
			Logger.Infof("adding i=%d", i)
			c <- map[string]int{add: i}
		}(i, add)
	}
	wg.Wait()
	Logger.Infof("calling done <- true")
	done <- true
}
