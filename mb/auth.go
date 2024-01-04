package mb

import (
	"context"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(c *gin.Context) {
	idToken := c.GetHeader("Authorization")

	token, err := AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(401, gin.H{
			"msg":  "invalid token",
			"data": err.Error(),
		})
		c.Abort()
		return
	}

	c.Set("token", token)
	c.Next()
}
