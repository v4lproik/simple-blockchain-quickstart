package nodes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/v4lproik/simple-blockchain-quickstart/common"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	. "github.com/v4lproik/simple-blockchain-quickstart/domains"
	log "go.uber.org/zap"
	"net/http"
)

const STATUS_NODE_ENDPOINT = "/status"
const BLOCKS_NODE_ENDPOINT = "/blocks"

type NodesEnv struct {
	errorBuilder ErrorBuilder
	nodeService  *NodeService
	stateService services.StateService
	blockService services.BlockService
}

func NodesRegister(router *gin.RouterGroup, env *NodesEnv) {
	router.GET(STATUS_NODE_ENDPOINT, env.NodeStatus)
	router.POST(BLOCKS_NODE_ENDPOINT, env.NodeListBlocks)
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

type ListBlocksParam struct {
	From string `json:"from" binding:"required,hash"`
}

// NodeListBlocks Get blocks from a specific hash specified in the payload
func (env NodesEnv) NodeListBlocks(c *gin.Context) {
	params := &ListBlocksParam{}
	//check params
	if err := ShouldBind(c, env.errorBuilder, "blocks cannot be listed", params); err != nil {
		AbortWithError(c, *err)
		return
	}

	//verified in parameter above
	hashFrom := models.Hash{}
	err := hashFrom.UnmarshalText([]byte(params.From))
	if err != nil {
		AbortWithError(c, *env.errorBuilder.NewUnknownError())
		return
	}
	log.S().Debugf("starting process of collecting blocks from hash=%s", params.From)

	blocks, err := env.blockService.GetNextBlocksFromHash(hashFrom)
	if err != nil {
		log.S().Error(fmt.Errorf("NodeListBlocks: couldn't retrieve blocks from DB: %w", err))
		AbortWithError(c, *env.errorBuilder.New(http.StatusInternalServerError, "blocks could not be retrieved"))
		return
	}

	//render
	serializer := BlocksSerializer{
		blocks: blocks,
	}
	c.JSON(http.StatusOK, gin.H{"blocks": serializer.Response()})
	return
}
