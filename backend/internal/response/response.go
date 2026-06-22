package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, status int, message string, payload gin.H) {
	out := gin.H{"message": message}
	for k, v := range payload {
		out[k] = v
	}
	c.JSON(status, out)
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"message": message})
}
