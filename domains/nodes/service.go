package nodes

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
	"os"
	"sync"
)

type NodeService struct {
	nodeDatabasePath string
	mu               sync.Mutex
}

type NetworkNodeName string
type NetworkNodeRecord struct {
	Address      string
	Port         uint64
	Is_bootstrap bool
	Is_active    bool
}
type NetworkNodeFromDB struct {
	Nodes map[NetworkNodeName]NetworkNodeRecord
}

func NewNodeService(nodeDatabasePath string) (*NodeService, error) {
	service := &NodeService{nodeDatabasePath: nodeDatabasePath}
	//check if the file can be opened and list the nodes
	if _, err := service.List(); err != nil {
		return nil, nil
	}
	return service, nil
}

// List nodes in the network if found, nil otherwise
func (u *NodeService) List() (map[NetworkNodeAddress]NetworkNode, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	file, err := ioutil.ReadFile(u.nodeDatabasePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load toml file %v", err)
	}

	var networkNodes NetworkNodeFromDB
	err = toml.Unmarshal(file, &networkNodes)
	if err != nil {
		return nil, fmt.Errorf("network nodes couldn't be extracted %v", err)
	}

	nodesNb := len(networkNodes.Nodes)
	nodes := make(map[NetworkNodeAddress]NetworkNode, nodesNb)
	i := 0
	for nodeName, node := range networkNodes.Nodes {
		nodes[NewNetworkNodeAddress(node.Address, node.Port)] = NetworkNode{
			Name:        string(nodeName),
			IsBootstrap: node.Is_bootstrap,
			IsActive:    node.Is_active,
		}
		i++
	}
	return nodes, nil
}

// Add nodes in the database, return error otherwise
func (u *NodeService) Add(nodes map[NetworkNodeAddress]NetworkNode) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	mapToInsert := make(map[NetworkNodeName]NetworkNodeRecord, len(nodes))
	for address, node := range nodes {
		mapToInsert[NetworkNodeName(node.Name)] = NetworkNodeRecord{
			Address:      address.ip,
			Port:         address.port,
			Is_bootstrap: node.IsBootstrap,
			Is_active:    node.IsActive,
		}
	}

	nodesToInsert := NetworkNodeFromDB{mapToInsert}
	byteNodes, err := toml.Marshal(&nodesToInsert)
	if err != nil {
		return fmt.Errorf("network nodes couldn't be marshalled %v", err)
	}

	err = os.WriteFile(u.nodeDatabasePath, byteNodes, 0644)
	if err != nil {
		return fmt.Errorf("cannot write in file %v", err)
	}

	return nil
}
