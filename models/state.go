package models

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type State struct {
	Balances         map[Account]uint
	transactionsPool []Transaction
	dbFile           *os.File
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
	state := &State{balances, make([]Transaction, 0), f}
	//for each transaction found in database
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		//extract a transaction
		var tx Transaction
		err := json.Unmarshal(scanner.Bytes(), &tx)
		if err != nil {
			return nil, err
		}

		//recreate state
		//check business logic of each transaction
		if err := state.apply(tx); err != nil {
			return nil, err
		}
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

func (s *State) Persist() error {
	txPool := make([]Transaction, len(s.transactionsPool))
	copy(txPool, s.transactionsPool)
	for i := 0; i < len(txPool); i++ {
		txJson, err := json.Marshal(txPool[i])
		if err != nil {
			return err
		}
		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}
		s.transactionsPool = s.transactionsPool[1:]
	}
	return nil
}

func (s *State) apply(tx Transaction) error {
	if tx.Reason == SELF_REWARD {
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
