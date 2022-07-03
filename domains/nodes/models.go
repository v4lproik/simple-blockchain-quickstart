package nodes

import (
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type NetworkNodeStatus struct {
	Hash         models.Hash
	Height       uint64
	NetworkNodes map[NetworkNodeAddress]NetworkNode
}

type NetworkNodeAddress struct {
	ip   string
	port uint64
}

func (n NetworkNodeAddress) String() string {
	return fmt.Sprintf("%s:%d", n.ip, n.port)
}

func (n NetworkNodeAddress) Ip() string {
	return n.ip
}

func (n NetworkNodeAddress) Port() uint64 {
	return n.port
}

func NewNetworkNodeAddress(ip string, port uint64) NetworkNodeAddress {
	return NetworkNodeAddress{
		ip,
		port,
	}
}

type NetworkNode struct {
	Name        string
	IsBootstrap bool
	IsActive    bool
}
