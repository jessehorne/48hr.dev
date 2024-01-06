package mb

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", DataResponse(c, gin.H{
		"creds": GetFirebaseClientCredentials(),
	}))
}

func PostProfile(c *gin.Context) {
	var isAuthed bool
	authed, _ := c.Cookie("is_authed")
	if authed == "true" {
		isAuthed = true
	}
	
	if !isAuthed {
		c.Redirect(301, "/")
		return
	}

	id, err := c.Cookie("user_id")
	if err != nil {
		c.Redirect(301, "/")
		return
	}
	
	u := GetUserByID(id)
	
	desc := c.PostForm("desc")
	u.Description = desc
	u.Save()
	
	c.Redirect(301, "/profile")
	return
}
