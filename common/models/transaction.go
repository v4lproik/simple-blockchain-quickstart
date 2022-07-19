package models

import (
	"crypto/sha256"
	"encoding/json"
)

// Transaction
type TransactionId Hash

type Transaction struct {
	From   Account `json:"from"`
	To     Account `json:"to"`
	Value  uint    `json:"value"`
	Reason string  `json:"reason"`
	Time   uint64  `json:"time"`
}

func NewTransaction(from Account, to Account, value uint, reason string, time uint64) *Transaction {
	return &Transaction{
		From:   from,
		To:     to,
		Value:  value,
		Reason: string(getReason(reason)),
		Time:   time,
	}
}

func (t Transaction) Hash() (TransactionId, error) {
	txJson, err := json.Marshal(t)
	if err != nil {
		return TransactionId(Hash{}), err
	}

	return sha256.Sum256(txJson), nil
}

// Reason transaction reason
type Reason string

const (
	OTHER       Reason = ""
	SELF_REWARD        = "self-reward"
	BIRTHDAY           = "birthday"
	LOAN               = "loan"
)

func getReason(reason string) Reason {
	switch reason {
	case "self-reward":
		return SELF_REWARD
	case "birthday":
		return BIRTHDAY
	case "loan":
		return LOAN
	}
	return OTHER
}

func (s Reason) IsValid() bool {
	switch s {
	case OTHER, SELF_REWARD, BIRTHDAY, LOAN:
		return true
	}

	return false
}
