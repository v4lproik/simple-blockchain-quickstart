package models

import (
	"crypto/sha256"
	"encoding/json"
)

type Block struct {
	Header BlockHeader   `json:"header"`
	Txs    []Transaction `json:"transactions"`
}

type BlockHeader struct {
	Parent Hash   `json:"parent"`
	Height uint64 `json:"height"`
	Nonce  uint32 `json:"nonce"`
	Time   uint64 `json:"time"`
}

type BlockDB struct {
	Hash  Hash  `json:"hash"`
	Block Block `json:"block"`
}

func NewBlock(parent Hash, height uint64, nonce uint32, time uint64, txs []Transaction) Block {
	return Block{
		BlockHeader{
			parent,
			height,
			nonce,
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
