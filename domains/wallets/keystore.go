package wallets

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

//fine to leave eth object as we are not planning on implementing another keystore
//only for test/muck purposes
type KeystoreService interface {
	NewKeystoreAccount(password string) (common.Address, error)
}

type EthKeystoreService struct {
	keystore *keystore.KeyStore
}

func NewEthKeystore(keystoreDataDirPath string) (*EthKeystoreService, error) {
	var path os.FileInfo
	var err error
	if path, err = os.Stat(keystoreDataDirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("NewEthKeystore: error os.State with keystoreDataDirPath %s: %w", keystoreDataDirPath, err)
	}

	if !path.IsDir() {
		return nil, errors.New("NewEthKeystore: " + keystoreDataDirPath + " is not a folder")
	}

	ks := keystore.NewKeyStore(keystoreDataDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		return nil, errors.New("NewEthKeystore: failed initiating keystore with configuration keystoreDataDirPath " + keystoreDataDirPath)
	}

	return &EthKeystoreService{
		keystore: ks,
	}, nil
}

func (k *EthKeystoreService) NewKeystoreAccount(password string) (common.Address, error) {
	acc, err := k.keystore.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return acc.Address, nil
}
