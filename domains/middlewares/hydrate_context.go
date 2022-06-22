package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/v4lproik/simple-blockchain-quickstart/domains"
)

const ERROR_BUILDER = "ERROR_BUILDER"

func ErrorBuilderToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ERROR_BUILDER, domains.NewErrorBuilder())
	}
}
