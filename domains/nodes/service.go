package nodes

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"io/ioutil"
)

type NodeService struct {
	nodeDatabasePath string
}

type NetworkNodeFromDB struct {
	Nodes map[string]struct {
		Address      string
		Port         uint64
		Is_bootstrap bool
		Is_active    bool
	}
}

func NewNodeService(nodeDatabasePath string) (*NodeService, error) {
	service := &NodeService{nodeDatabasePath: nodeDatabasePath}
	//check if the file can be opened and list the nodes
	if _, err := service.List(); err != nil {
		return nil, nil
	}
	return service, nil
}

// Get user if found, nil otherwise
func (u *NodeService) List() ([]NetworkNode, error) {
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
	nodes := make([]NetworkNode, nodesNb)
	i := 0
	for _, node := range networkNodes.Nodes {
		nodes[i] = NetworkNode{
			Ip:          node.Address,
			Port:        node.Port,
			IsBootstrap: node.Is_bootstrap,
			IsActive:    node.Is_active,
		}
		i++
	}
	return nodes, nil
}
