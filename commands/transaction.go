package commands

import (
	"errors"
	"fmt"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
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
	state              models.State
	transactionService services.TransactionService
	From               string `short:"f" long:"from" description:"Transaction account to get tokens from" required:"true"`
	To                 string `short:"t" long:"to" description:"Transaction account to send tokens to" required:"true"`
	Value              uint   `short:"v" long:"value" description:"Value of the transaction" required:"true"`
	Reason             string `short:"r" long:"reason" description:"Reason of the transaction" required:"false"`
}

func NewAddTransactionCommand(state models.State) (*AddTransactionCommand, error) {
	if state == nil {
		return nil, errors.New("NewAddTransactionCommand: state cannot be nil")
	}
	add := new(AddTransactionCommand)
	add.state = state
	add.transactionService = services.NewFileTransactionService()

	return add, nil
}

func checkArgs(c AddTransactionCommand) (models.Account, models.Account, error) {
	from, err := models.NewAccount(c.From)
	if err != nil {
		return from, "", errors.New("checkArgs: from variable is not a valid ethereum account")
	}

	to, err := models.NewAccount(c.To)
	if err != nil {
		return to, from, errors.New("checkArgs: to variable is not a valid ethereum account")
	}

	if c.Value <= 0 {
		return from, to, errors.New("checksArgs: value needs to be a positive value")
	}
	return from, to, nil
}

func (c *AddTransactionCommand) Execute(_ []string) error {
	// check args
	from, to, err := checkArgs(*c)
	if err != nil {
		return fmt.Errorf("Execute: error checking args: %s", err)
	}
	// create transaction object
	tx := models.NewTransaction(from, to, c.Value, c.Reason)

	// get the state
	state := c.state
	_, err = c.transactionService.AddTransaction(state, tx)
	if err != nil {
		return fmt.Errorf("Execute: error adding tx: %s", err)
	}

	state.Print()
	return nil
}

type ListTransactionCommand struct {
	state models.State
}

func NewListTransactionCommand(state models.State) (*ListTransactionCommand, error) {
	if state == nil {
		return nil, errors.New("NewListTransactionCommand: state cannot be nil")
	}
	list := new(ListTransactionCommand)
	list.state = state

	return list, nil
}

func (c *ListTransactionCommand) Execute(_ []string) error {
	c.state.Print()
	return nil
}
