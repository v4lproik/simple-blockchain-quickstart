package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/v4lproik/simple-blockchain-quickstart/commands"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
	"os"
)

type EnvVal string

const (
	Dev  = "dev"
	Prod = "prod"
)

func (e EnvVal) isValid() bool {
	switch e {
	case Dev:
		return true
	case Prod:
		return true
	}
	return false
}

func (e EnvVal) isProd() bool {
	if e == Prod {
		return true
	}
	return false
}

var opts struct {
	RunAsHttpserver      bool   `short:"r" long:"run_as_http_server" description:"Run the application as an http server" required:"false"`
	UsersFilePath        string `short:"u" long:"users_file_path" description:"Users file path" required:"true"`
	GenesisFilePath      string `short:"g" long:"genesis_file_path" description:"Genesis file path" required:"true"`
	TransactionsFilePath string `short:"d" long:"transactions_file_path" description:"Transactions file path" required:"true"`
	NodesFilePath        string `short:"n" long:"nodes_file_path" description:"Nodes file path" required:"true"`
	KeystoreDirPath      string `short:"k" long:"keystore_dir_path" description:"Keystore dir path" required:"true"`
	LogFilePath          string `short:"l" long:"log_file_path" description:"Where application logs will be written. If this value is not specified, the logs will be displayed to the console" required:"false"`
	Environment          string `short:"e" long:"environment" description:"Set the environment variable. Accepted values are [dev, prod]" required:"false" default:"dev"`
}

func displayAppConfiguration() {
	Logger.Infof("Environment: %s", opts.Environment)
	Logger.Infof("Transactions file: %s", opts.TransactionsFilePath)
	Logger.Infof("Genesis file: %s", opts.GenesisFilePath)
	Logger.Infof("Users file: %s", opts.UsersFilePath)
	Logger.Infof("Nodes file: %s", opts.NodesFilePath)
	Logger.Infof("Keystore dir: %s", opts.KeystoreDirPath)
	if opts.LogFilePath != "" {
		Logger.Infof("Output in log file: %s", opts.LogFilePath)
	} else {
		Logger.Infof("Output: console")
	}
}

func checkArgs() {
	env = EnvVal(opts.Environment)
	if !env.isValid() {
		fmt.Println("environment " + opts.Environment + " is not accepted. Choose from [dev, prod]. Exiting")
		os.Exit(1)
	}
}

//general commands
func addCommands(parser *flags.Parser) error {
	err := addTransactionCommands(parser)
	if err != nil {
		return fmt.Errorf("cannot add transaction commands %v", err)
	}

	err = addPasswordCommands(parser)
	if err != nil {
		return fmt.Errorf("cannot add password commands %v", err)
	}

	return nil
}

//transaction
func addTransactionCommands(parser *flags.Parser) error {
	state, err := models.NewStateFromFile(opts.GenesisFilePath, opts.TransactionsFilePath)
	if err != nil {
		Logger.Fatalf("addTransactionCommands: %w", err)
	}

	addT, _ := commands.NewAddTransactionCommand(state)
	listT, _ := commands.NewListTransactionCommand(state)
	_, err = parser.AddCommand(
		"transaction",
		"transaction utility commands including: add, list",
		"Utilities developed to ease the operations and debugging of transactions.",
		&commands.TransactionCommands{
			Add:  *addT,
			List: *listT,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

//password
func addPasswordCommands(parser *flags.Parser) error {
	_, err := parser.AddCommand(
		"password",
		"password utility commands including: hash",
		"Utilities developed to ease the operations and debugging of password.",
		&commands.PasswordCommands{
			Hash:    *commands.NewHashAPasswordCommand(),
			Compare: *commands.NewCompareHashCommand(),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
