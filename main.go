package main

import (
	"github.com/jessevdk/go-flags"
)

func main() {
	//parse cli arguments
	//ignoreunknown is important here as we parse the cli several times as we declare commands and subcommands after the first cli "parsing"
	//it means that some arguments are declared after the first cli parsing and so the command will be rejected if we don't specify the
	//ignore unknown flag.
	parser := flags.NewParser(&opts, flags.IgnoreUnknown)

	//parse command line arguments
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	//init logger
	InitLogger(opts.LogFilePath)
	defer logger.Sync()

	//display program configuration
	displayAppConfiguration()

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
