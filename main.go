package main

import (
	"github.com/jessevdk/go-flags"
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
		runHttpServer()
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
