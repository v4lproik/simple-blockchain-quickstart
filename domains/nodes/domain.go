package nodes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
)

const NODES_DOMAIN_URL = "/api/nodes"

func RunDomain(r *gin.Engine, nodeService *NodeService, stateService services.StateService, middlewares ...gin.HandlerFunc) {
	v1 := r.Group(NODES_DOMAIN_URL)
	for _, middleware := range middlewares {
		v1.Use(middleware)
	}

	//register http endpoints
	NodesRegister(v1.Group("/"), &NodesEnv{
		nodeService:  nodeService,
		stateService: stateService,
		errorBuilder: common.NewErrorBuilder(),
	})

	//run background tasks
	manager := NewNodeTaskManager(5, nodeService, stateService)
	ctx := context.Background()
	go manager.Run(ctx)
}
