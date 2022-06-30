package nodes

import "github.com/v4lproik/simple-blockchain-quickstart/common/models"

type StatusNode struct {
	Hash         models.Hash
	Number       uint64
	NetworkNodes []NetworkNode
}

type NetworkNode struct {
	Ip          string
	Port        uint64
	IsBootstrap bool
	IsActive    bool
}
