package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "go.uber.org/zap"
	"io/ioutil"
	"os"
	"time"
)

type State struct {
	Balances         map[Account]uint
	transactionsPool []Transaction
	dbFile           *os.File
	latestBlockHash  Hash
}

type GenesisFile struct {
	Time     time.Time       `json:"genesis_time"`
	ChainId  string          `json:"chain_id"`
	Balances map[string]uint `json:"balances"`
}

func NewState(genesisFilePath string, transactionFilePath string) (*State, error) {
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
		balances[NewAccount(account)] = balance
	}

	//read transactions database
	f, err := os.OpenFile(transactionFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Transaction, 0), f, Hash{}}

	//for each block found in database
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockDB BlockDB
		err = json.Unmarshal(blockFsJson, &blockDB)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockDB.Value)
		if err != nil {
			return nil, err
		}

		//the hash reflecting the state is now the latest block being added to the database
		state.latestBlockHash = blockDB.Key
	}

	return state, nil
}

func (s *State) Add(tx Transaction) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.transactionsPool = append(s.transactionsPool, tx)
	return nil
}

func (s *State) Persist() (Hash, error) {
	hash := Hash{}

	//create a new Block only with the new transactions
	block := NewBlock(
		s.latestBlockHash,
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
	blockDBJson, err := json.Marshal(blockDB)
	if err != nil {
		return hash, err
	}

	//add to the DB the new block as well as a new line
	_, err = s.dbFile.Write(append(blockDBJson, '\n'))
	if err != nil {
		return hash, err
	}

	//latest block of the state is now the hash of the latest block inserted into the database
	s.latestBlockHash = blockHash

	//empty the transactions pool as it should only transactions that haven't been written to database yet
	s.transactionsPool = []Transaction{}

	return s.latestBlockHash, nil

}

func (s *State) apply(tx Transaction) error {
	if tx.Reason == SELF_REWARD {
		//refuse the transaction if it's a self reward with different from/to address
		if !tx.To.isSameAccount(tx.From) {
			return fmt.Errorf("from and to accounts should be the same as self-reward as been specified as a reason for the transaction")
		}
		s.Balances[tx.To] += tx.Value
		return nil
	}
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

//getters
func (s *State) GetLatestBlockHash() Hash {
	return s.latestBlockHash
}

//print
func (s *State) Print() {
	log.S().Infof("#####################")
	log.S().Infof("# Accounts balances #")
	log.S().Infof("#####################")
	log.S().Infof("State: %x", s.GetLatestBlockHash())
	log.S().Infof("---------------------")
	for account, balance := range s.Balances {
		log.S().Infof("%s: %d", account, balance)
	}
	log.S().Infof("---------------------")
}