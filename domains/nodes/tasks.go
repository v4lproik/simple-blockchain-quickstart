package nodes

import (
	"context"
	"encoding/json"
	"fmt"
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
}

func NewNodeTaskManager(refreshInterval uint64, nodeService *NodeService, stateService services.StateService) *NodeTaskManager {
	return &NodeTaskManager{
		refreshIntervalInSeconds: refreshInterval,
		nodeService:              nodeService,
		stateService:             stateService,
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
		log.S().Errorf("trying to get node status %s status", address.String())
		status, err := getNodeStatus(address)
		if err != nil {
			log.S().Errorf("unable to get node %s status %v", address.String(), err)
			continue
		}

		currentHeight := state.GetLatestBlockHeight()
		if currentHeight < status.Height {
			missingBlockCount := status.Height - currentHeight
			log.S().Debugf("new blocks (%d) needs to be added", missingBlockCount)
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

	//TODO Do not use default http client
	res, err := http.Get(url)
	if err != nil {
		return NetworkNodeStatus{}, err
	}

	return getStatusNode(res)
}

func getStatusNode(r *http.Response) (NetworkNodeStatus, error) {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unable to read response body %s", err)
	}
	defer r.Body.Close()

	var response *NetworkNodesResponse
	err = json.Unmarshal(reqBodyJson, response)
	if err != nil {
		return NetworkNodeStatus{}, fmt.Errorf("unable to unmarshal response body %s", err)
	}

	if response == nil {
		return NetworkNodeStatus{}, fmt.Errorf("node response body is nil")
	}

	statusNode := NetworkNodeStatus{}
	statusNode.Hash = response.Hash
	statusNode.Height = response.Height

	statusNode.NetworkNodes = make(map[NetworkNodeAddress]NetworkNode, len(response.NetworkNodeResponse))
	for _, nodeResponse := range response.NetworkNodeResponse {
		statusNode.NetworkNodes[NewNetworkNodeAddress(nodeResponse.Ip, nodeResponse.Port)] = NetworkNode{
			Name:        nodeResponse.Name,
			IsBootstrap: nodeResponse.IsBootstrap,
			IsActive:    nodeResponse.IsActive,
		}
	}

	return statusNode, nil
}
