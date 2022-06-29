package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Block struct {
	Header BlockHeader   // metadata (parent block hash + time)
	Txs    []Transaction // new transactions only (payload)
}

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Height uint64 `json:"height"`
	Time   uint64 `json:"time"`
}

type BlockDB struct {
	Key   Hash  `json:"hash"`
	Value Block `json:"block"`
}

func NewBlock(parent Hash, height uint64, time uint64, txs []Transaction) Block {
	return Block{
		BlockHeader{
			parent,
			height,
			time,
		},
		txs,
	}
}

func (b Block) Hash() (Hash, error) {
	blockJson, err := json.Marshal(b)
	if err != nil {
		return Hash{}, err
	}
	return sha256.Sum256(blockJson), nil
}
