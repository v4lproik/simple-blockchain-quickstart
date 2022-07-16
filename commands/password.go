package commands

import (
	"encoding/base64"
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

type PasswordCommands struct {
	Hash    HashAPasswordCommand `command:"hash" description:"Hash a password"`
	Compare CompareHashCommand   `command:"compare" description:"Compare a hash to a password"`
}

type HashAPasswordCommand struct {
	passwordService services.PasswordService
	Password        string `short:"p" long:"password" description:"Password to hash" required:"true"`
}

func NewHashAPasswordCommand() *HashAPasswordCommand {
	return &HashAPasswordCommand{
		passwordService: services.NewDefaultPasswordService(),
	}
}

func (c *HashAPasswordCommand) Execute(args []string) error {
	hash, err := c.passwordService.GenerateHash(c.Password)
	if err != nil {
		return fmt.Errorf("Execute: cannot generate a hash: %w", err)
	}
	Logger.Infof("hash: %s", hash)
	return nil
}

type CompareHashCommand struct {
	passwordService services.PasswordService
	Hash            string `short:"b" long:"hash" description:"Hash encoded in base64 to compare" required:"true"`
	Password        string `short:"p" long:"password" description:"Password to hash" required:"true"`
}

func NewCompareHashCommand() *CompareHashCommand {
	return &CompareHashCommand{
		passwordService: services.NewDefaultPasswordService(),
	}
}

func (c *CompareHashCommand) Execute(args []string) error {
	hash, err := base64.StdEncoding.DecodeString(c.Hash)
	if err != nil {
		return fmt.Errorf("Execute: error decoding base64 hash: %w", err)
	}

	isPassword, err := c.passwordService.ComparePasswordAndHash(c.Password, string(hash))
	if err != nil {
		return fmt.Errorf("Execute: error comparing hashes: %w", err)
	}
	verb := "do not"
	if isPassword {
		verb = "do"
	}
	Logger.Infof("the password and the hash %s match", verb)
	return nil
}
