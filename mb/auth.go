package mb

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func FirebaseAuthMiddleware(c *gin.Context) {
	var idToken string
	splittedToken := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splittedToken) < 2 {
		idToken = c.Query("token")
	}

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
