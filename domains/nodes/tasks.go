package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	log "go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"
)

type NodeTaskManager struct {
	refreshIntervalInSeconds uint64
	nodeService              *NodeService
	stateService             services.StateService
	blockService             services.BlockService
}

func NewNodeTaskManager(refreshInterval uint64, nodeService *NodeService, stateService services.StateService,
	blockService services.BlockService) *NodeTaskManager {
	return &NodeTaskManager{
		refreshIntervalInSeconds: refreshInterval,
		nodeService:              nodeService,
		stateService:             stateService,
		blockService:             blockService,
	}
}

func (n *NodeTaskManager) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(n.refreshIntervalInSeconds))

	for {
		select {
		case <-ticker.C:
			log.S().Debugf("looking for new nodes within the network")
			err := n.getOtherNodesViaNodeStatus()
			if err != nil {
				log.S().Errorf("error looking for new nodes within the network %v", err)
			}
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *NodeTaskManager) getOtherNodesViaNodeStatus() error {
	knownNetworkNodes, err := n.nodeService.List()
	if err != nil {
		return fmt.Errorf("error listing nodes: %s", err)
	}

	if len(knownNetworkNodes) == 0 {
		log.S().Debugf("no network nodes found... no sync...")
	}

	state, err := n.stateService.GetState()
	if err != nil {
		return fmt.Errorf("couldn't retrieve blockchain state %v", err)
	}

	for address, _ := range knownNetworkNodes {
		log.S().Errorf("trying to get node status %s", address.String())
		status, err := getNodeStatus(address)
		if err != nil {
			log.S().Errorf("unable to get node %s status %v", address.String(), err)
			continue
		}

		currentHeight := state.GetLatestBlockHeight()
		if currentHeight < status.Height {
			missingBlockCount := status.Height - currentHeight
			currentHash := state.GetLatestBlockHash()

			log.S().Debugf("new blocks (%d) needs to be added", missingBlockCount)
			//sync database from that node
			// get the blocks from other node
			blocks, err := getNextNodeBlocksFromHash(address, currentHash)
			if err != nil {
				return err
			}

			// insert the new block into our own database

		}

		for networkNodeIp, newNode := range status.NetworkNodes {
			_, isKnownNode := knownNetworkNodes[networkNodeIp]
			if !isKnownNode {
				log.S().Debugf("found new node with address %s", networkNodeIp)
				knownNetworkNodes[networkNodeIp] = newNode
			}
		}
	}

	return nil
}

func getNodeStatus(nodeAddress NetworkNodeAddress) (NetworkNodeStatus, error) {
	url := fmt.Sprintf("http://%s%s%s", nodeAddress.String(), NODES_DOMAIN_URL, STATUS_NODE_ENDPOINT)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	//TODO Do not use default http client
	cc := &http.Client{}
	response, err := cc.Do(req)
	if err != nil {
		return NetworkNodeStatus{}, err
	}

	return getStatusNode(response)
}

type NodeGetStatusResponse struct {
	Status NetworkNodesResponse `json:"status"`
}

func getStatusNode(r *http.Response) (NetworkNodeStatus, error) {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unable to read response body %s", err)
	}
	defer r.Body.Close()

	var response NodeGetStatusResponse
	err = json.Unmarshal(reqBodyJson, &response)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unable to unmarshal response body %s", err)
	}

	statusNode := NetworkNodeStatus{}
	statusNode.Hash = response.Status.Hash
	statusNode.Height = response.Status.Height

	statusNode.NetworkNodes = make(map[NetworkNodeAddress]NetworkNode, len(response.Status.NetworkNodeResponse))
	for _, nodeResponse := range response.Status.NetworkNodeResponse {
		statusNode.NetworkNodes[NewNetworkNodeAddress(nodeResponse.Ip, nodeResponse.Port)] = NetworkNode{
			Name:        nodeResponse.Name,
			IsBootstrap: nodeResponse.IsBootstrap,
			IsActive:    nodeResponse.IsActive,
		}
	}

	return statusNode, nil
}

func getNextNodeBlocksFromHash(nodeAddress NetworkNodeAddress, hash models.Hash) ([]models.BlockDB, error) {
	// generate url
	url := fmt.Sprintf("http://%s%s%s", nodeAddress.String(), NODES_DOMAIN_URL, BLOCKS_NODE_ENDPOINT)

	// generate payload
	listBlocksParam := ListBlocksParam{}
	hashStr, _ := hash.MarshalText()
	listBlocksParam.From = string(hashStr)

	// marshall payload
	body, _ := json.Marshal(listBlocksParam)

	//TODO Do not use default http client
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	//TODO Do not use default http client
	cc := &http.Client{}
	res, err := cc.Do(req)
	if err != nil {
		return []models.BlockDB{}, err
	}

	return getBlocks(res)
}

type NodeListBlocksResponse struct {
	Blocks []BlockWrapperResponse `json:"blocks"`
}

func getBlocks(r *http.Response) ([]models.BlockDB, error) {
	var blocks []models.BlockDB

	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return blocks, fmt.Errorf("unable to read response body %s", err)
	}
	defer r.Body.Close()
	log.S().Infof("%s", reqBodyJson)
	var response NodeListBlocksResponse
	err = json.Unmarshal(reqBodyJson, &response)
	if err != nil {
		return blocks, fmt.Errorf("unable to unmarshal response body %s", err)
	}

	var blocksRes NodeListBlocksResponse
	err = json.Unmarshal(reqBodyJson, &blocksRes)
	if err != nil {
		return blocks, err
	}

	blocks = make([]models.BlockDB, len(blocksRes.Blocks))
	for i, block := range blocksRes.Blocks {
		// init block structure
		blockDB := models.BlockDB{
			Hash: block.Hash,
			Block: models.Block{
				Header: models.BlockHeader{
					Parent: block.Block.Header.Parent,
					Height: block.Block.Header.Height,
					Time:   block.Block.Header.Time,
				},
			},
		}

		// add transactions
		txs := make([]models.Transaction, len(block.Block.Txs))
		for y, tx := range block.Block.Txs {
			txs[y] = models.Transaction{
				From:   tx.From,
				To:     tx.To,
				Value:  tx.Value,
				Reason: tx.Reason,
			}
		}
		blockDB.Block.Txs = txs

		// add to array of blocks
		blocks[i] = blockDB
	}

	return blocks, nil
}
