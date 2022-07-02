package nodes

import "github.com/v4lproik/simple-blockchain-quickstart/common/models"

type StatusNode struct {
	Hash         models.Hash
	Height       uint64
	NetworkNodes map[NetworkNodeIp]NetworkNode
}

type NetworkNodeIp string
type NetworkNode struct {
	Name        string
	Port        uint64
	IsBootstrap bool
	IsActive    bool
}

//utils
