package mb

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func DiscordAuthMiddleware(c *gin.Context) {
	isAuthed, err := c.Cookie("is_authed")

	if isAuthed != "true" || err != nil {
		c.JSON(http.StatusUnauthorized, []byte("{}"))
		c.Abort()
		return
	}
}

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

func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GetAuthCallback(c *gin.Context) {
	token, err := OAuthConf.Exchange(context.Background(), c.Query("code"))

	// look down
	/* The TOKEN HERE !!! it has access token and refresh token which apparently should be stored here...*/
	// Look up

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "something went really bad with auth",
			"error": err.Error(),
		})
		return
	}

	// Step 4: Use the access token, here we use it to get the logged in user's info.
	res, err := OAuthConf.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		c.JSON(500, gin.H{
			"msg":    "something went wrong with getting your profile",
			"error":  err.Error(),
			"status": res.StatusCode,
		})
		return
	}

	defer res.Body.Close()

	d, _ := io.ReadAll(res.Body)
	//fmt.Println(string(d))

	// set cookie to be used in further auth requests
	authTimeFor := 86400 * 7 // 7 days
	c.SetCookie("is_authed", "true", authTimeFor, "/", os.Getenv("APP_DOMAIN"), true, true)

	// create/update user in firestore
	var discUser DiscordUser
	err = json.Unmarshal(d, &discUser)
	if err != nil {
		// TODO
		fmt.Println(err)
		return
	}
	foundUser := CreateOrUpdateUser(discUser)

	fmt.Println("Creating: ", foundUser.ID)

	c.SetCookie("user_id", foundUser.ID, authTimeFor, "/", os.Getenv("APP_DOMAIN"), true, true)

	c.Redirect(http.StatusFound, "/profile?a=1")
}
