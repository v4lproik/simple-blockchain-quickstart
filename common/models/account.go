package models

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type Account string

func NewAccount(account string) (Account, error) {
	if !common.IsHexAddress(account) {
		return "", errors.New("NewAccount: account " + account + " variable is not a valid ethereum account")
	}
	return Account(account), nil
}

func (acc *Account) isSameAccount(toCompare Account) bool {
	return fmt.Sprintf("%s", *acc) == fmt.Sprintf("%s", toCompare)
}
