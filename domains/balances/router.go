package balances

import (
	"net/http"

	. "github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
)

const LIST_BALANCES_ENDPOINT = "/"

type BalancesEnv struct {
	state        models.State
	errorBuilder ErrorBuilder
}

func NewBalancesEnv(state models.State, builder ErrorBuilder) *BalancesEnv {
	return &BalancesEnv{
		state:        state,
		errorBuilder: builder,
	}
}

func BalancesRegister(router *gin.RouterGroup, env *BalancesEnv) {
	router.POST(LIST_BALANCES_ENDPOINT, env.ListBalances)
}

func (env *BalancesEnv) ListBalances(c *gin.Context) {
	state := env.state

	if len(state.Balances()) == 0 {
		AbortWithError(c, *env.errorBuilder.New(http.StatusNotFound, "balances could not be found"))
		return
	}

	// map state with state response
	serializer := BalancesSerializer{state.Balances()}

	// render
	c.JSON(http.StatusOK, gin.H{"balances": serializer.Response()})
}
