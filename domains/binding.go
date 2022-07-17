package domains

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/v4lproik/simple-blockchain-quickstart/common"
)

// gin-Gonic Bind Format Error
// Inspired from https://www.convictional.com/blog/gin-validation
type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "enum":
		return "The value cannot be submitted"
	case "password":
		return "The password doesn't comply with the policy (min 8 char with min 1 upper, 1 number and 1 symbol)"
	case "account":
		return "The account is not an Ethereum style account (eg. 0xa6aa1c9106f0c0d0895bb72f40cfc830180ebeaf)"
	case "hash":
		return "The hash should be a 32 byte array"
	}
	return "Unknown error"
}

func ShouldBind(c *gin.Context, errorBuilder common.ErrorBuilder, errMsg string, params interface{}) *common.Error {
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
