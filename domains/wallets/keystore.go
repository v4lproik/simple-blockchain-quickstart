package wallets

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

type KeystoreService struct {
	keystore *keystore.KeyStore
}

func NewKeystore(keystoreDataDirPath string) (*KeystoreService, error) {
	var path os.FileInfo
	var err error
	if path, err = os.Stat(keystoreDataDirPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot initiate a keystore with configuration: keystoreDataDirPath %v", keystoreDataDirPath)
	}

	if !path.IsDir() {
		return nil, fmt.Errorf("cannot initiate a keystore with keystoreDataDirPath %v not a folder", keystoreDataDirPath)
	}

	ks := keystore.NewKeyStore(keystoreDataDirPath, keystore.StandardScryptN, keystore.StandardScryptP)
	if ks == nil {
		return nil, fmt.Errorf("cannot initiate a keystore with configuration: keystoreDataDirPath %v", keystoreDataDirPath)
	}

	return &KeystoreService{
		keystore: ks,
	}, nil
}

func (k KeystoreService) NewKeystoreAccount(password string) (common.Address, error) {
	acc, err := k.keystore.NewAccount(password)
	if err != nil {
		return common.Address{}, err
	}

	return acc.Address, nil
}
