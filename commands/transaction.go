package commands

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models/conf"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
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
	stateService       services.StateService
	transactionService services.TransactionService
	From               string `short:"f" long:"from" description:"Transaction account to get tokens from" required:"true"`
	To                 string `short:"t" long:"to" description:"Transaction account to send tokens to" required:"true"`
	Value              uint   `short:"v" long:"value" description:"Value of the transaction" required:"true"`
	Reason             string `short:"r" long:"reason" description:"Reason of the transaction" required:"false"`
}

func NewAddTransactionCommand(genesisFilePath string, transactionsFilePath string) *AddTransactionCommand {
	add := new(AddTransactionCommand)
	add.stateService = services.NewFileStateService(conf.NewBlockchainFileDatabaseConf(genesisFilePath, transactionsFilePath))
	add.transactionService = services.NewFileTransactionService()

	return add
}

func (c *AddTransactionCommand) Execute(args []string) error {
	//create transaction object
	tx := models.NewTransaction(models.NewAccount(c.From), models.NewAccount(c.To), c.Value, c.Reason)

	//get the state
	state, err := c.stateService.GetState()
	if err != nil {
		return err
	}

	_, err = c.transactionService.AddTransaction(state, tx)
	if err != nil {
		return err
	}

	state.Print()
	return nil
}

type ListTransactionCommand struct {
	stateService services.StateService
}

func NewListTransactionCommand(genesisFilePath string, transactionsFilePath string) *ListTransactionCommand {
	list := new(ListTransactionCommand)
	list.stateService = services.NewFileStateService(conf.NewBlockchainFileDatabaseConf(genesisFilePath, transactionsFilePath))

	return list
}

func (c *ListTransactionCommand) Execute(args []string) error {
	state, err := c.stateService.GetState()
	if err != nil {
		return err
	}
	state.Print()

	return nil
}
