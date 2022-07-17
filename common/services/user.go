package services

import (
	"fmt"
	"io/ioutil"

	"github.com/pelletier/go-toml/v2"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type UserService struct {
	userDatabasePath string
}

// Name toml specific struct
type Name string

type UserFromDB struct {
	Nodes map[Name]struct {
		Hash string
	}
}

// NewUserService initiate the user service
func NewUserService(userDatabasePath string) (*UserService, error) {
	service := &UserService{userDatabasePath: userDatabasePath}
	// check if the file can be opened and list the nodes
	if _, err := service.List(); err != nil {
		return nil, nil
	}
	return service, nil
}

// Get user if found, nil otherwise
func (u *UserService) Get(userName string) (*models.User, error) {
	// toml v2 has removed the querying language
	// small file at the moment so this works atm - this should be addressed at some point cf. kanban board
	users, err := u.List()
	if err != nil {
		return nil, fmt.Errorf("Get: failed to list users: %w", err)
	}

	if val, ok := users[Name(userName)]; ok {
		return &val, nil
	}

	return nil, nil
}

// List all users, empty array otherwise
func (u *UserService) List() (map[Name]models.User, error) {
	file, err := ioutil.ReadFile(u.userDatabasePath)
	if err != nil {
		return nil, fmt.Errorf("List: failed to read user database: %w", err)
	}

	var usersFromDb UserFromDB
	err = toml.Unmarshal(file, &usersFromDb)
	if err != nil {
		return nil, fmt.Errorf("List: failed to unmarshal users: %w", err)
	}

	usersNb := len(usersFromDb.Nodes)
	users := make(map[Name]models.User, usersNb)
	i := 0
	for name, userProperties := range usersFromDb.Nodes {
		users[name] = models.User{
			Hash: userProperties.Hash,
		}
		i++
	}
	return users, nil
}
