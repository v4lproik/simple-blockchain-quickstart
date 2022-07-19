package nodes

import (
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

// nodes
type NodeSerializer struct {
	models.State
	nodes      map[NetworkNodeAddress]NetworkNode
	pendingTxs []models.Transaction
}

type NetworkNodeResponse struct {
	Name        string `json:"name"`
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
	i := 0
	for address, node := range n.nodes {
		nodesResponse[i] = NetworkNodeResponse{
			Name:        node.Name,
			Ip:          address.ip,
			Port:        address.port,
			IsBootstrap: node.IsBootstrap,
			IsActive:    node.IsActive,
		}
		i++
	}
	response.NetworkNodeResponse = nodesResponse

	return *response
}

// blocks
type BlockSerializer struct {
	block models.Block
}

type BlocksSerializer struct {
	blocks []models.Block
}

type TransactionResponse struct {
	From   models.Account `json:"from"`
	To     models.Account `json:"to"`
	Value  uint           `json:"value"`
	Reason string         `json:"reason"`
}

type BlockHeaderResponse struct {
	Parent models.Hash `json:"parent"`
	Height uint64      `json:"height"`
	Time   uint64      `json:"time"`
}

type BlockResponse struct {
	Header BlockHeaderResponse   `json:"header"`
	Txs    []TransactionResponse `json:"transactions"`
}

type BlocksResponse struct {
	Blocks []BlockResponse `json:"blocks"`
}

func (n *BlockSerializer) Response() BlockResponse {
	response := BlockResponse{}
	block := n.block

	// add block metadata
	response = BlockResponse{
		Header: BlockHeaderResponse{
			Parent: block.Header.Parent,
			Height: block.Header.Height,
			Time:   block.Header.Time,
		},
	}

	// add block transactions
	txRes := make([]TransactionResponse, len(block.Txs))
	for i, tx := range block.Txs {
		txRes[i] = TransactionResponse{
			From:   tx.From,
			To:     tx.To,
			Value:  tx.Value,
			Reason: tx.Reason,
		}
	}
	response.Txs = txRes

	return response
}

func (n *BlocksSerializer) Response() BlocksResponse {
	response := make([]BlockResponse, len(n.blocks))
	for i, block := range n.blocks {
		serializer := BlockSerializer{block: block}
		response[i] = serializer.Response()
	}
	return BlocksResponse{response}
}
