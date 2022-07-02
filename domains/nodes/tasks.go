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
	refreshInterval uint64
	nodeService     NodeService
	stateService    services.StateService
}

func NewNodeTaskManager(refreshInterval uint64, nodeService NodeService) *NodeTaskManager {
	return &NodeTaskManager{
		refreshInterval: refreshInterval,
		nodeService:     nodeService,
	}
}

func (n *NodeTaskManager) Run(ctx context.Context) {
	ticker := time.NewTimer(time.Second * time.Duration(n.refreshInterval))

	for {
		select {
		case <-ticker.C:
			log.S().Debugf("looking for new nodes within the network")
			err := n.getOtherNodesViaNodeStatus()
			if err != nil {
				return
			}
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (n *NodeTaskManager) getOtherNodesViaNodeStatus() error {
	networkNodes, err := n.nodeService.List()
	if err != nil {
		return fmt.Errorf("error listing nodes: %s", err)
	}

	if len(networkNodes) == 0 {
		log.S().Debugf("no network nodes found... no sync...")
	}

	state, err := n.stateService.GetState()
	if err != nil {
		return fmt.Errorf("couldn't retrieve blockchain state %v", err)
	}

	_, err = n.nodeService.List()
	if err != nil {
		return fmt.Errorf("couldn't retrieve known nodes %v", err)
	}

	for ip, node := range networkNodes {
		status, err := getNodeStatus(ip, node)
		if err != nil {
			log.S().Errorf("unable to get node %s:%s status", ip, node.Port)
			continue
		}

		currentHeight := state.GetLatestBlockHeight()
		if currentHeight < status.Height {
			missingBlockCount := status.Height - currentHeight
			log.S().Debugf("new blocks (%d) needs to be added", missingBlockCount)
		}
	}
	return nil
}

func getNodeStatus(nodeIp NetworkNodeIp, node NetworkNode) (StatusNode, error) {
	url := fmt.Sprintf("http://%s/%s/%s", nodeIp, NODES_DOMAIN_URL, STATUS_NODE_ENDPOINT)

	res, err := http.Get(url)
	if err != nil {
		return StatusNode{}, err
	}

	return getStatusNode(res)
}

func getStatusNode(r *http.Response) (StatusNode, error) {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return StatusNode{}, fmt.Errorf("unable to read response body %s", err)
	}
	defer r.Body.Close()

	var response *NetworkNodesResponse
	err = json.Unmarshal(reqBodyJson, response)
	if err != nil {
		return StatusNode{}, fmt.Errorf("unable to unmarshal response body %s", err)
	}

	if response == nil {
		return StatusNode{}, fmt.Errorf("node response body is nil")
	}

	statusNode := StatusNode{}
	statusNode.Hash = response.Hash
	statusNode.Height = response.Height

	statusNode.NetworkNodes = make(map[NetworkNodeIp]NetworkNode, len(response.NetworkNodeResponse))
	for _, nodeResponse := range response.NetworkNodeResponse {
		statusNode.NetworkNodes[NetworkNodeIp(nodeResponse.Ip)] = NetworkNode{
			Name:        nodeResponse.Name,
			Port:        nodeResponse.Port,
			IsBootstrap: nodeResponse.IsBootstrap,
			IsActive:    nodeResponse.IsActive,
		}
	}

	return statusNode, nil
}
