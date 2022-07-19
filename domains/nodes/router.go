package nodes

import (
	"fmt"
	"net/http"

	. "github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

const (
	STATUS_NODE_ENDPOINT = "/status"
	BLOCKS_NODE_ENDPOINT = "/blocks"
)

type NodesEnv struct {
	nodeService  *NodeService
	state        models.State
	blockService services.BlockService
}

func NodesRegister(router *gin.RouterGroup, env *NodesEnv) {
	router.GET(STATUS_NODE_ENDPOINT, env.NodeStatus)
	router.POST(BLOCKS_NODE_ENDPOINT, env.NodeListBlocks)
}

func (env NodesEnv) NodeStatus(c *gin.Context) {
	// get all the nodes
	nodes, err := env.nodeService.List()
	if err != nil {
		AbortWithError(c, NewError(http.StatusInternalServerError, "nodes could not be found"))
		return
	}

	// init serializer
	serializer := &NodeSerializer{
		State: env.state,
		nodes: nodes,
	}

	// render
	c.JSON(http.StatusOK, gin.H{"status": serializer.Response()})
}

type ListBlocksParam struct {
	From string `json:"from" binding:"required,hash"`
}

// NodeListBlocks Get blocks from a specific hash specified in the payload
func (env NodesEnv) NodeListBlocks(c *gin.Context) {
	params := &ListBlocksParam{}
	// check params
	if err := ShouldBind(c, "blocks cannot be listed", params); err != nil {
		AbortWithError(c, err)
		return
	}

	// verified in parameter above
	hashFrom := models.Hash{}
	err := hashFrom.UnmarshalText([]byte(params.From))
	if err != nil {
		AbortWithError(c, NewUnknownError())
		return
	}
	Logger.Debugf("starting process of collecting blocks from hash=%s", params.From)

	blocks, err := env.blockService.GetNextBlocksFromHash(hashFrom)
	if err != nil {
		Logger.Error(fmt.Errorf("NodeListBlocks: couldn't retrieve blocks from DB: %w", err))
		AbortWithError(c, NewError(http.StatusInternalServerError, "blocks could not be retrieved"))
		return
	}

	// render
	serializer := BlocksSerializer{
		blocks: blocks,
	}

	c.JSON(http.StatusOK, serializer.Response())
}
