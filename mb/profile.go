package mb

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xeonx/timeago"
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

func GetUserProfile(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		return
	}

	// get user from db
	us, err := StoreClient.Collection("users").Where("ID", "==", id).Documents(context.Background()).GetAll()

	var user *User
	for _, u := range us {
		err := u.DataTo(&user)
		if err != nil {
			// TODO
		}
		break
	}
	user.EnglishCreatedAt = timeago.English.Format(user.CreatedAt)

	// get projects
	ps, err := StoreClient.Collection("posts").
		Documents(context.Background()).
		GetAll()
	if err != nil {
		// TODO
	}

	var allProjects []*Project
	for _, p := range ps {
		var newP *Project
		err := p.DataTo(&newP)
		if err != nil {
			// TODO
		}
		newP.EnglishCreatedTime = timeago.English.Format(newP.CreatedAt)
		newP.EnglishStartedTime = timeago.English.Format(newP.StartedAt)

		if len(newP.Tags) > 0 {
			splitted := strings.Split(newP.Tags, ",")
			for _, s := range splitted {
				newP.FormattedTags = append(newP.FormattedTags, s)
			}
		}

		if newP.Title != "Centrifuge" {
			if newP.UserID == user.ID {
				// only add if its not the first doc and it belongs to the user
				allProjects = append(allProjects, newP)
			}
		}
	}

	if err != nil {
		c.HTML(http.StatusOK, "user.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
		}))
		return
	} else {
		c.HTML(http.StatusOK, "user.html", DataResponse(c, gin.H{
			"creds":    GetFirebaseClientCredentials(),
			"projects": allProjects,
			"user":     user,
		}))
		return
	}
}
