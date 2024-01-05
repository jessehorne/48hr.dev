package mb

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(c *gin.Context) {
	idToken := c.GetHeader("Authorization")
	idToken = strings.Split(idToken, " ")[1]

	token, err := AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(401, gin.H{
			"msg":  "invalid token from middleware",
			"data": err.Error(),
		})
		c.Abort()
		return
	}

	c.Set("token", token)
	c.Set("userID", token.UID)
	c.Next()
}
