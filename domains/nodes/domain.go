package nodes

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
)

const NODES_DOMAIN_URL = "/api/nodes"

func RunDomain(
	r *gin.Engine,
	nodeService *NodeService,
	state models.State,
	transactionService services.TransactionService,
	blockService services.BlockService,
	taskManagerSyncInterval uint32,
	middlewares ...gin.HandlerFunc,
) error {
	v1 := r.Group(NODES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	// register http endpoints
	NodesRegister(v1.Group("/"), &NodesEnv{
		nodeService:  nodeService,
		state:        state,
		blockService: blockService,
	})

	// run background tasks
	manager, err := NewNodeTaskManager(
		taskManagerSyncInterval,
		nodeService,
		state,
		transactionService,
		blockService,
	)
	if err != nil {
		return fmt.Errorf("RunDomain: node task manager cannot start: %w", err)
	}

	ctx := context.Background()
	go manager.RunMine(ctx)

	ctx = context.Background()
	go manager.RunSync(ctx)
	return nil
}
