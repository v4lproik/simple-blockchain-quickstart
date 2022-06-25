package domains

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

//gin-Gonic Bind Format Error
//Inspired from https://www.convictional.com/blog/gin-validation
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
	}
	return "Unknown error"
}

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