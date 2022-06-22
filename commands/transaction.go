package commands

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	log "go.uber.org/zap"
)

type TransactionCommands struct {
	Add  AddTransactionCommand  `command:"add" description:"Add a new transaction"`
	List ListTransactionCommand `command:"list" description:"List all transactions"`
}

type TransactionCommandsOpts struct {
	GenesisFilePath      string
	TransactionsFilePath string
}

type AddTransactionCommand struct {
	opts   TransactionCommandsOpts
	From   string `short:"f" long:"from" description:"Transaction account to get tokens from" required:"true"`
	To     string `short:"t" long:"to" description:"Transaction account to send tokens to" required:"true"`
	Value  uint   `short:"v" long:"value" description:"Value of the transaction" required:"true"`
	Reason string `short:"r" long:"reason" description:"Reason of the transaction" required:"false"`
}

func NewAddTransactionCommand(genesisFilePath string, transactionsFilePath string) *AddTransactionCommand {
	add := new(AddTransactionCommand)
	add.opts.GenesisFilePath = genesisFilePath
	add.opts.TransactionsFilePath = transactionsFilePath

	return add
}

func (c *AddTransactionCommand) Execute(args []string) error {
	//get command line information
	tx := models.NewTransaction(models.NewAccount(c.From), models.NewAccount(c.To), c.Value, c.Reason)

	//load state
	state, err := models.NewStateFromFile(c.opts.GenesisFilePath, c.opts.TransactionsFilePath)
	if err != nil {
		log.S().Fatalf("cannot get blockchain state: %v", err)
	}
	defer state.Close()

	//add transaction to state
	err = state.Add(*tx)
	if err != nil {
		log.S().Fatalf("cannot add transaction to state: %v", err)
	}

	//persist new state to disk
	_, err = state.Persist()
	if err != nil {
		log.S().Fatalf("cannot persist state to disk: %v", err)
	}

	state.Print()

	return nil
}

type ListTransactionCommand struct {
	opts TransactionCommandsOpts
}

func NewListTransactionCommand(genesisFilePath string, transactionsFilePath string) *ListTransactionCommand {
	list := new(ListTransactionCommand)
	list.opts.GenesisFilePath = genesisFilePath
	list.opts.TransactionsFilePath = transactionsFilePath

	return list
}

func (c *ListTransactionCommand) Execute(args []string) error {
	state, err := models.NewStateFromFile(c.opts.GenesisFilePath, c.opts.TransactionsFilePath)
	if err != nil {
		log.S().Fatalf("cannot get blockchain state: %v", err)
	}
	defer state.Close()

	state.Print()

	return nil
}
