package healthz

import "github.com/gin-gonic/gin"

func RunDomain(r *gin.Engine) {
	r.GET("/api/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
