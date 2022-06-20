package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/v4lproik/simple-blockchain-quickstart/commands"
)

var opts struct {
	RunAsHttpserver      bool   `short:"r" long:"run_as_http_server" description:"Run the application as an http server" required:"false"`
	GenesisFilePath      string `short:"g" long:"genesis_file_path" description:"Genesis file path" required:"true"`
	TransactionsFilePath string `short:"d" long:"transactions_file_path" description:"Transactions file path" required:"true"`
	LogFilePath          string `short:"l" long:"log_file_path" description:"Where application logs will be written. If this value is not specified, the logs will be displayed to the console." required:"false"`
}

func displayAppConfiguration() {
	logger.Infof("Transactions file: %s", opts.TransactionsFilePath)
	logger.Infof("Genesis file: %s", opts.GenesisFilePath)
	//logger.Infof("Run as http server: %s", opts.RunAsHttpserver)
	if opts.LogFilePath != "" {
		logger.Infof("Output in log file: %s", opts.LogFilePath)
	} else {
		logger.Infof("Output: console")
	}
}

//general commands
func addCommands(parser *flags.Parser) error {
	return addTransactionCommands(parser)
}

//transaction
func addTransactionCommands(parser *flags.Parser) error {
	_, err := parser.AddCommand(
		"transaction",
		"transaction utility commands including: add, list",
		"Utilities developed to ease the operations and debugging of transactions.",
		&commands.TransactionCommands{
			Add:  *commands.NewAddTransactionCommand(opts.GenesisFilePath, opts.TransactionsFilePath),
			List: *commands.NewListTransactionCommand(opts.GenesisFilePath, opts.TransactionsFilePath),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
