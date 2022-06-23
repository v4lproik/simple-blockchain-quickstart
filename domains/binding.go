package domains

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ShouldBind(c *gin.Context, errorBuilder ErrorBuilder, errMsg string, params interface{}) *Error {
	if err := c.ShouldBind(params); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			return errorBuilder.New(http.StatusBadRequest, errMsg, out)
		}
		return errorBuilder.New(http.StatusBadRequest, errMsg, err)
	}
	return nil
}
