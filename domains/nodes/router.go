package nodes

import (
	"github.com/gin-gonic/gin"
	. "github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	"net/http"
)

const STATUS_NODE_ENDPOINT = "/status"

type NodesEnv struct {
	errorBuilder ErrorBuilder
	nodeService  *NodeService
	stateService services.StateService
}

func NodesRegister(router *gin.RouterGroup, env *NodesEnv) {
	router.GET(STATUS_NODE_ENDPOINT, env.NodeStatus)
}

func (env NodesEnv) NodeStatus(c *gin.Context) {
	//get all the nodes
	nodes, err := env.nodeService.List()
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(http.StatusInternalServerError, "nodes could not be found"))
		return
	}

	//get state
	state, err := env.stateService.GetState()
	if err != nil {
		AbortWithError(c, *env.errorBuilder.New(http.StatusInternalServerError, "state could not be found"))
		return
	}

	//init serializer
	serializer := &NodeSerializer{
		State: state,
		nodes: nodes,
	}

	//render
	c.JSON(http.StatusOK, gin.H{"status": serializer.Response()})
	return
}
