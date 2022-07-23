package nodes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

type BlockHeight uint64

type NodeTaskManager struct {
	refreshIntervalInSeconds uint32

	state models.State

	nodeService        *NodeService
	transactionService services.TransactionService
	blockService       services.BlockService
}

// NewNodeTaskManager handles all the background tasks needed for a node to sync its status
// as well as mining new blocks
func NewNodeTaskManager(
	refreshInterval uint32,
	nodeService *NodeService,
	state models.State,
	transactionService services.TransactionService,
	blockService services.BlockService,
) (*NodeTaskManager, error) {
	if state == nil {
		return nil, errors.New("NewNodeTaskManager: state cannot be nil")
	}
	if nodeService == nil {
		return nil, errors.New("NewNodeTaskManager: node service cannot be nil")
	}

	return &NodeTaskManager{
		refreshIntervalInSeconds: refreshInterval,
		nodeService:              nodeService,
		state:                    state,
		transactionService:       transactionService,
		blockService:             blockService,
	}, nil
}

// RunMine starts mining a new block when a new transaction is being submitted
func (n *NodeTaskManager) RunMine(ctx context.Context) {
	Logger.Debugf("RunMine: Start mining...")
	for {
		select {
		case pendingTx := <-n.transactionService.NewPendingTxs():
			Logger.Debugf("RunMine: Received new transaction...")
			_, err := n.blockService.Mine(context.Background(), models.PendingBlock{
				Parent:       n.state.GetLatestBlockHash(),
				Height:       n.state.GetLatestBlockHeight() + 1,
				Time:         utils.DefaultTimeService.UnixUint64(),
				MinerAddress: n.blockService.ThisNodeMiningAddress(),
				Txs: []models.Transaction{
					*models.NewTransaction(pendingTx.From, pendingTx.To, pendingTx.Value, pendingTx.Reason, pendingTx.Time),
				},
			})
			if err != nil {
				Logger.Errorf("RunMine: ")
			}
		case <-ctx.Done():
			Logger.Debugf("RunMine: stop mining...")
			return
		}
	}
}

// RunSync starts the process of syncing the list of nodes within the network as well as this node's database
func (n *NodeTaskManager) RunSync(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(n.refreshIntervalInSeconds))

	for {
		select {
		case <-ticker.C:
			Logger.Debugf("RunSync: looking for new nodes within the network")

			// first fetch the nodes' status within the network
			// status contains block height and other peers in network
			nodeStatus, err := n.runFetchNodeStatus()
			if err != nil {
				Logger.Errorf("RunSync: failed to lookup to new nodes: %s", err)
			}

			// time to synchronise our database as we have other nodes' status block height
			err = n.runSyncNode(nodeStatus)
			if err != nil {
				Logger.Errorf("RunSync: failed to synchronise: %s", err)
			}
		case <-ctx.Done():
			Logger.Debugf("RunSync: Stop looking for new nodes within the network")
			ticker.Stop()
		}
	}
}

func (n *NodeTaskManager) runFetchNodeStatus() (map[NetworkNodeAddress]NetworkNodeStatus, error) {
	knownNetworkNodes, err := n.nodeService.List()
	if err != nil {
		return nil, fmt.Errorf("runFetchNodeStatus: failed to list nodes: %w", err)
	}
	if len(knownNetworkNodes) == 0 {
		Logger.Debugf("runFetchNodeStatus: no network nodes found... no sync...")
		return nil, nil
	}

	// this is the minimum amount of calls that we'll be making
	// log(n) calls is to be expected as each node has a X number of unknown nodes
	// that we might be calling
	wg := &sync.WaitGroup{}
	wg.Add(len(knownNetworkNodes))

	done := make(chan bool)
	c := make(chan map[NetworkNodeAddress]NetworkNodeStatus)

	// we use a buffered channel, no need for safe concurrent map
	nodeStatus := make(map[NetworkNodeAddress]NetworkNodeStatus, len(knownNetworkNodes))

	go fetchNodesHeights(done, wg, c, knownNetworkNodes)

waitLoop:
	for {
		select {
		case val := <-c:
			for i, s := range val {
				nodeStatus[i] = s
			}
		case <-done:
			close(c)
			break waitLoop
		}
	}
	close(done)

	return nodeStatus, nil
}

func fetchNodesHeights(done chan<- bool, wg *sync.WaitGroup, nodeStatus chan map[NetworkNodeAddress]NetworkNodeStatus, nodes map[NetworkNodeAddress]NetworkNode) {
	// asynchronously loop over each knownNode
	for address := range nodes {
		go func(address NetworkNodeAddress) {
			defer wg.Done()

			// call the node and get its status
			// send in channel node height or 0 if the node is not reachable
			status, err := getNodeStatus(address)
			if err != nil {
				nodeStatus <- map[NetworkNodeAddress]NetworkNodeStatus{address: status}
				Logger.Warnf("runFetchNodeStatus: failed to reach node: %s", err)
				return
			}
			nodeStatus <- map[NetworkNodeAddress]NetworkNodeStatus{address: status}

			// each node status returns its own knownNodes information
			// we also want to reach out to those nodes as their heights might be closer to the world state
			for networkNodeIp, newNode := range status.NetworkNodes {
				_, isKnownNode := nodes[networkNodeIp]
				// we only want to reach out to the node if it's not known and active
				if !isKnownNode && newNode.IsActive {
					nodes[networkNodeIp] = newNode

					// increment delta to one as we run a new goroutine
					wg.Add(1)
					go func(address NetworkNodeAddress) {
						defer wg.Done()

						// call the node and get its status
						// send in channel node height or 0 if the node is not reachable
						status, err := getNodeStatus(address)
						if err != nil {
							nodeStatus <- map[NetworkNodeAddress]NetworkNodeStatus{address: status}
							Logger.Warnf("runFetchNodeStatus: failed to reach node: %s", err)
							return
						}
						nodeStatus <- map[NetworkNodeAddress]NetworkNodeStatus{address: status}
					}(networkNodeIp)
				}
			}
		}(address)
	}

	wg.Wait()
	done <- true
}

func (n *NodeTaskManager) runSyncNode(nodeStatus map[NetworkNodeAddress]NetworkNodeStatus) error {
	Logger.Debugf("runSyncNode: synchronisation has started")

	// get current block height
	state := n.state
	currentHeight := state.GetLatestBlockHeight()

	// find the node with the highest block height and higher than ours
	// if condition met, then synchronise, otherwise do not do anything
	var nodeToSynchFrom map[NetworkNodeAddress]NetworkNodeStatus
	for address, status := range nodeStatus {
		if currentHeight < status.Height {
			if nodeToSynchFrom == nil {
				nodeToSynchFrom = make(map[NetworkNodeAddress]NetworkNodeStatus, 1)
			}
			nodeToSynchFrom[address] = status
		}
	}

	// skip sync if couldn't find any node with a higher block height
	if nodeToSynchFrom == nil {
		Logger.Debugf("runSyncNode: couldn't find a node with a higher block height than us")
		return nil
	}

	// start sync the node
	for address, status := range nodeToSynchFrom {
		missingBlockCount := status.Height - currentHeight
		currentHash := state.GetLatestBlockHash()

		Logger.Debugf("runSyncNode: starting synchronisation, new blocks (%d) needs from %s to be added", missingBlockCount, address.String())
		blocks, err := getNextNodeBlocksFromHash(address, currentHash)
		if err != nil {
			return fmt.Errorf("runSyncNode: failed at fetching blocks from node to sychronise from: %w", err)
		}

		// insert the new blocks into our database
		if err = state.AddBlocks(blocks); err != nil {
			return fmt.Errorf("runSyncNode: failed to add blocks into database: %w", err)
		}
	}

	Logger.Debugf("runSyncNode: synchronisation is over")
	return nil
}

func getNodeStatus(nodeAddress NetworkNodeAddress) (NetworkNodeStatus, error) {
	url := fmt.Sprintf("http://%s%s%s", nodeAddress.String(), NODES_DOMAIN_URL, STATUS_NODE_ENDPOINT)

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	// TODO Do not use default http client
	cc := &http.Client{}
	response, err := cc.Do(req)
	if err != nil {
		return NetworkNodeStatus{}, err
	}

	return unmarshalNodeStatus(response)
}

type NodeGetStatusResponse struct {
	Status NetworkNodesResponse `json:"status"`
}

func unmarshalNodeStatus(r *http.Response) (NetworkNodeStatus, error) {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unmarshalNodeStatus: failed to read response body: %w", err)
	}
	defer r.Body.Close()

	var response NodeGetStatusResponse
	err = json.Unmarshal(reqBodyJson, &response)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unmarshalNodeStatus: failed to unmarshall body: %w", err)
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

func getNextNodeBlocksFromHash(nodeAddress NetworkNodeAddress, hash models.Hash) ([]models.Block, error) {
	// generate url
	url := fmt.Sprintf("http://%s%s%s", nodeAddress.String(), NODES_DOMAIN_URL, BLOCKS_NODE_ENDPOINT)

	// generate payload
	listBlocksParam := ListBlocksParam{}
	hashStr, _ := hash.MarshalText()
	listBlocksParam.From = string(hashStr)

	// marshall payload
	body, _ := json.Marshal(listBlocksParam)

	// TODO Do not use default http client
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// TODO Do not use default http client
	cc := &http.Client{}
	res, err := cc.Do(req)
	if err != nil {
		return []models.Block{}, err
	}

	return getBlocks(res)
}

func getBlocks(r *http.Response) ([]models.Block, error) {
	var blocks []models.Block

	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return blocks, fmt.Errorf("unable to read response body %s", err)
	}
	defer r.Body.Close()

	var response BlocksResponse
	err = json.Unmarshal(reqBodyJson, &response)
	if err != nil {
		return blocks, fmt.Errorf("unable to unmarshal response body %s", err)
	}

	var blocksRes BlocksResponse
	err = json.Unmarshal(reqBodyJson, &blocksRes)
	if err != nil {
		return blocks, err
	}

	blocks = make([]models.Block, len(blocksRes.Blocks))
	for i, blockRes := range blocksRes.Blocks {
		// init block structure
		block := models.Block{
			Header: models.BlockHeader{
				Parent: blockRes.Header.Parent,
				Height: blockRes.Header.Height,
				Time:   blockRes.Header.Time,
			},
		}

		// add transactions
		txs := make([]models.Transaction, len(blockRes.Txs))
		for y, tx := range blockRes.Txs {
			txs[y] = models.Transaction{
				From:   tx.From,
				To:     tx.To,
				Value:  tx.Value,
				Reason: tx.Reason,
			}
		}
		block.Txs = txs

		// add to array of blocks
		blocks[i] = block
	}

	return blocks, nil
}
