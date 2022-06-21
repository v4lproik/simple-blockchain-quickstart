package models

import (
	"fmt"
)

type Account string

func NewAccount(name string) Account {
	return Account(name)
}

func (acc *Account) isSameAccount(toCompare Account) bool {
	return fmt.Sprintf("%s", *acc) == fmt.Sprintf("%s", toCompare)
}
