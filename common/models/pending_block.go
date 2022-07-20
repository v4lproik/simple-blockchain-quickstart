package models

type PendingBlock struct {
	Parent       Hash
	Height       uint64
	Time         uint64
	MinerAddress Account
	Txs          []Transaction
}

func NewPendingBlock(parent Hash, height uint64, minerAddress Account, time uint64, txs []Transaction) PendingBlock {
	return PendingBlock{
		Parent:       parent,
		Height:       height,
		Time:         time,
		MinerAddress: minerAddress,
		Txs:          txs,
	}
}
