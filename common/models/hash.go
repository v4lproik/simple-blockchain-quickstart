package models

import (
	"encoding/hex"
	"reflect"
)

type Hash [32]byte

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(h.Hex()), nil
}

func (h *Hash) UnmarshalText(data []byte) error {
	_, err := hex.Decode(h[:], data)
	return err
}

func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

func CompareBlockHash(h1 Hash, h2 Hash) bool {
	return reflect.DeepEqual(h1, h2)
}
