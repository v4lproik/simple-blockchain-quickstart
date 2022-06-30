package nodes

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

type NodeSerializer struct {
	models.State
	nodes []NetworkNode
}

type NetworkNodeResponse struct {
	Ip          string `json:"ip"`
	Port        uint64 `json:"port"`
	IsBootstrap bool   `json:"is_bootstrap"`
	IsActive    bool   `json:"is_active"`
}

type NetworkNodesResponse struct {
	Hash                models.Hash           `json:"block_hash"`
	Height              uint64                `json:"block_height"`
	NetworkNodeResponse []NetworkNodeResponse `json:"network_nodes"`
}

func (n *NodeSerializer) Response() NetworkNodesResponse {
	response := new(NetworkNodesResponse)
	response.Hash = n.State.GetLatestBlockHash()
	response.Height = n.State.GetLatestBlockHeight()

	nodesResponse := make([]NetworkNodeResponse, len(n.nodes))
	for i, node := range n.nodes {
		nodesResponse[i] = NetworkNodeResponse{
			Ip:          node.Ip,
			Port:        node.Port,
			IsBootstrap: node.IsBootstrap,
			IsActive:    node.IsActive,
		}
	}
	response.NetworkNodeResponse = nodesResponse

	return *response
}
