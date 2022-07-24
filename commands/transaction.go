package commands

import (
	"errors"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type TransactionCommands struct {
	List ListTransactionCommand `command:"list" description:"List all transactions"`
}

type TransactionCommandsOpts struct {
	GenesisFilePath      string
	TransactionsFilePath string
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
