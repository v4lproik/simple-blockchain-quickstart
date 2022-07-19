package nodes

import (
	"context"

	"github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
)

const NODES_DOMAIN_URL = "/api/nodes"

func RunDomain(
	r *gin.Engine,
	nodeService *NodeService,
	state models.State,
	blockService services.BlockService,
	taskManagerSyncInterval uint32,
	middlewares ...gin.HandlerFunc,
) {
	v1 := r.Group(NODES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	// register http endpoints
	NodesRegister(v1.Group("/"), &NodesEnv{
		nodeService:  nodeService,
		state:        state,
		blockService: blockService,
		errorBuilder: utils.NewErrorBuilder(),
	})

	// run background tasks
	manager, _ := NewNodeTaskManager(taskManagerSyncInterval, nodeService, state, blockService)
	ctx := context.Background()
	go manager.Run(ctx)
}
