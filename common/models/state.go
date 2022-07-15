package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	log "go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
)

var (
	ErrNextBlockHeight = errors.New("latest block height doesn't match with next block (height + 1)")
	ErrNextBlockHash   = errors.New("latest block hash doesn't match with next block")
)

type GenesisFile struct {
	Time     time.Time       `json:"genesis_time"`
	ChainId  string          `json:"chain_id"`
	Balances map[string]uint `json:"balances"`
}

type (
	State interface {
		// Add adds a transaction
		Add(Transaction) error
		// AddBlock to the state
		AddBlock(Block) error
		// AddBlocks to the state
		AddBlocks([]Block) error
		// Balances return the balances as map
		Balances() map[Account]uint
		Persist() (Hash, error)
		Close() error
		GetLatestBlockHash() Hash
		GetLatestBlockHeight() uint64
		Print()
	}
)

type FromFileState struct {
	balances         map[Account]uint
	transactionsPool []Transaction
	dbFile           *os.File
	latestBlockHash  Hash
	latestBlock      Block
}

func NewStateFromFile(genesisFilePath string, transactionFilePath string) (*FromFileState, error) {
	//read genesis file
	file, err := ioutil.ReadFile(genesisFilePath)
	if err != nil {
		return nil, err
	}

	//extract genesis file information into struct
	data := GenesisFile{}
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)
	for account, balance := range data.Balances {
		acc, err := NewAccount(account)
		if err != nil {
			return nil, fmt.Errorf("error building state from file %v", err)
		}
		balances[acc] = balance
	}

	//read transactions database
	db, err := getTransactionsDb(transactionFilePath)
	if err != nil {
		return nil, err
	}

	state, err := getFileStateFromFile(balances, db)
	if err != nil {
		return state, err
	}
	return state, nil
}

func getFileStateFromFile(balances map[Account]uint, db *os.File) (*FromFileState, error) {
	state := &FromFileState{balances, make([]Transaction, 0), db, Hash{}, Block{}}

	//for each block found in database
	scanner := bufio.NewScanner(db)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockDB BlockDB
		err := json.Unmarshal(blockFsJson, &blockDB)
		if err != nil {
			return nil, err
		}

		// we do not call applyBlocks here
		// we are initiating the state from the initial database containing legit blocks, so it's
		// safe not to apply any business logic on the blocks themselves
		err = state.applyTxs(blockDB.Block.Txs)
		if err != nil {
			return nil, err
		}

		//keep a copy of the latest block and its hash,
		// so it can be exposed to the network
		state.latestBlockHash = blockDB.Hash
		state.latestBlock = blockDB.Block
	}
	return state, nil
}

func getTransactionsDb(transactionFilePath string) (*os.File, error) {
	return os.OpenFile(transactionFilePath, os.O_APPEND|os.O_RDWR, 0600)
}

func (s *FromFileState) Balances() map[Account]uint {
	return s.balances
}

func (s *FromFileState) Add(tx Transaction) error {
	if err := s.applyTx(tx); err != nil {
		return err
	}
	s.transactionsPool = append(s.transactionsPool, tx)
	return nil
}

func (s *FromFileState) AddBlock(block Block) error {
	// as we use the transaction pool to persist transactions, we have two choices here
	// either we force a flush by persisting every pending transaction or we copy the state,
	// we block the state until sync is done and then we re-establish the state and accept
	// new transactions which will be added to the freshly re-restablished state
	// let's go for the latter

	// TODO: create benchmark
	// our state is a pointer so we need to copy its value
	var copiedStateFromFile FromFileState
	err := copier.CopyWithOption(&copiedStateFromFile, s, copier.Option{DeepCopy: true})
	if err != nil {
		return fmt.Errorf("AddBlock: cannot copy the state: %w", err)
	}
	log.S().Debugf("this %d", s.latestBlock.Header.Height)
	log.S().Debugf("copy %d", copiedStateFromFile.latestBlock.Header.Height)

	// validate the block
	err = copiedStateFromFile.applyBlock(block)
	if err != nil {
		return err
	}

	// create a blockFS, ready to be added to the state
	blockHash, err := block.Hash()
	if err != nil {
		return err
	}

	blockDB := BlockDB{
		Hash:  blockHash,
		Block: block,
	}
	err = s.persistBlockToDB(blockDB)
	if err != nil {
		return err
	}

	// now the blocks have been written in the DB
	// the state (copy) balance has been updated each time a block has been inserted
	// in the database. As no error happened during the writing process, we
	// then need to update the state (original).
	s.balances = copiedStateFromFile.Balances()
	s.latestBlock = block
	s.latestBlockHash = blockHash

	return nil
}

func (s *FromFileState) AddBlocks(blocks []Block) error {
	for _, block := range blocks {
		err := s.AddBlock(block)
		if err != nil {
			return fmt.Errorf("AddBlocks: %w", err)
		}
	}
	return nil
}

func (s *FromFileState) Persist() (Hash, error) {
	hash := Hash{}

	//create a new Block only with the new transactions
	block := NewBlock(
		s.latestBlockHash,
		s.latestBlock.Header.Height+1,
		uint64(time.Now().Unix()),
		s.transactionsPool,
	)
	//generate block hash
	blockHash, err := block.Hash()
	if err != nil {
		return hash, err
	}

	//create database block which includes its hash and the transactions (block itself)
	blockDB := BlockDB{blockHash, block}
	err = s.persistBlockToDB(blockDB)
	if err != nil {
		return blockHash, err
	}

	//latest block of the state is now the hash of the latest block inserted into the database
	s.latestBlockHash = blockHash
	s.latestBlock = blockDB.Block

	//empty the transactions pool as it should only transactions that haven't been written to database yet
	s.transactionsPool = []Transaction{}

	return s.latestBlockHash, nil
}

func (s *FromFileState) persistBlockToDB(block BlockDB) error {
	blockDBJson, err := json.Marshal(block)
	if err != nil {
		return err
	}

	//add to the DB the new block as well as a new line
	_, err = s.dbFile.Write(append(blockDBJson, '\n'))
	if err != nil {
		return err
	}
	return nil
}

// applyBlock checks if a block can be added to the database
// also checks if the blocks which is trying to be added has previousBlock (or parentBlock)
// is block.height == previousBlock.height + 1 and that its previousBlock.parentHash points to block.hash
func (s *FromFileState) applyBlock(block Block) error {
	log.S().Debugf("%d == %d", block.Header.Height, s.latestBlock.Header.Height+1)
	if block.Header.Height != s.latestBlock.Header.Height+1 {
		return fmt.Errorf("applyBlock: %w", ErrNextBlockHeight)
	}

	log.S().Debugf("s.latestBlockHash %s == block.Header.Parent %s", s.latestBlockHash, block.Header.Parent)
	if !CompareBlockHash(s.latestBlockHash, block.Header.Parent) {
		return fmt.Errorf("applyBlock: %w", ErrNextBlockHash)
	}

	return s.applyTxs(block.Txs)
}

// applyTxs is a wrapper calling applyTx and propagate error if any
func (s *FromFileState) applyTxs(txs []Transaction) error {
	for _, tx := range txs {
		if err := s.applyTx(tx); err != nil {
			return err
		}
	}
	return nil
}

// applyTx checks if a transaction can be added to the blockchain
// also checks if the account has enough money as well as the transaction metadata is valid
func (s *FromFileState) applyTx(tx Transaction) error {
	if tx.Reason == SELF_REWARD {
		//refuse the transaction if it's a self reward with different from/to address
		if !tx.To.isSameAccount(tx.From) {
			return fmt.Errorf("from and to accounts should be the same as self-reward as been specified as a reason for the transaction")
		}
		s.balances[tx.To] += tx.Value
		return nil
	}
	if tx.Value > s.balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}
	s.balances[tx.From] -= tx.Value
	s.balances[tx.To] += tx.Value
	return nil
}

func (s *FromFileState) Close() error {
	return s.dbFile.Close()
}

func (s *FromFileState) GetLatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *FromFileState) GetLatestBlockHeight() uint64 {
	return s.latestBlock.Header.Height
}

func (s *FromFileState) Print() {
	log.S().Infof("#####################")
	log.S().Infof("# Accounts balances #")
	log.S().Infof("#####################")
	log.S().Infof("State: %x", s.GetLatestBlockHash())
	log.S().Infof("Height: %x", s.GetLatestBlockHeight())
	log.S().Infof("---------------------")
	for account, balance := range s.balances {
		log.S().Infof("%s: %d", account, balance)
	}
	log.S().Infof("---------------------")
}
