package middleware

import (
	"encoding/json"

	. "github.com/v4lproik/simple-blockchain-quickstart/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/v4lproik/simple-blockchain-quickstart/common/models"
	"github.com/v4lproik/simple-blockchain-quickstart/common/services"
	Logger "github.com/v4lproik/simple-blockchain-quickstart/log"
)

const (
	AUTH_HEADER = "X-API-TOKEN"
)

// A helper to write user_id and user_model to the context
func UpdateUserContext(c *gin.Context, user models.User) {
	c.Set("my_user", user)
}

func AuthWebSessionMiddleware(auto401 bool, errorBuilder ErrorBuilder, jwtService *services.JwtService) gin.HandlerFunc {
	Logger.Debugf("authentication is %s", auto401)
	return func(c *gin.Context) {
		// if authentication not required
		if !auto401 {
			c.Next()
			return
		}

		// if authentication is required
		// extract token
		jwtToken := c.Request.Header.Get(AUTH_HEADER)
		if jwtToken == "" {
			AbortWithError(c, *errorBuilder.New(401, "authentication token cannot be found"))
			return
		}
		// verify the token if present
		authToken, err := jwtService.VerifyToken(jwtToken)
		if err != nil {
			context := ""
			if err.Error() == "Token is expired" {
				context = "token is expired"
			}
			AbortWithError(c, *errorBuilder.New(401, "authentication token is not valid", context))
			return
		}

		// extract info in claims
		if claims, ok := authToken.Claims.(jwt.MapClaims); ok && authToken.Valid {
			data := claims["dat"]
			if data == nil {
				AbortWithError(c, *errorBuilder.New(401, "authentication token is not valid", "there's no data"))
				return
			}

			// unmarshall payload from claims["dat"]
			var user models.User
			err := json.Unmarshal([]byte(data.(string)), &user)
			if err != nil {
				AbortWithError(c, *errorBuilder.New(500, "cannot parse payload"))
				return
			}

			// add user to gin context
			UpdateUserContext(c, user)
		}
	}
}
