package services

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type UserService struct {
	userDb *toml.Tree
}

func NewUserService(userDatabasePath string) (*UserService, error) {
	file, err := toml.LoadFile(userDatabasePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load toml file %v", err)
	}
	return &UserService{userDb: file}, nil
}

// Get user if found, nil otherwise
func (u *UserService) Get(userName string) (*models.User, error) {
	p := u.userDb.Get(fmt.Sprintf("Users.%s", userName)).(*toml.Tree)
	if p == nil {
		return nil, fmt.Errorf("user %s cannot be found", userName)
	}
	return &models.User{
		Name: userName,
		Hash: p.Get("password").(string),
	}, nil
}
