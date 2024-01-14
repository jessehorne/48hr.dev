package mb

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/xeonx/timeago"
)

func GetIndex(c *gin.Context) {
	// get projects
	ps, err := StoreClient.Collection("posts").OrderBy("CreatedAt", firestore.Desc).
		Documents(context.Background()).GetAll()

	var allProjects []*Project
	for _, p := range ps {
		var newP *Project
		err := p.DataTo(&newP)
		newP.EnglishTime = timeago.English.Format(newP.CreatedAt)

		if err != nil {
			// TODO
		}

		if len(newP.Tags) > 0 {
			splitted := strings.Split(newP.Tags, ",")
			for _, s := range splitted {
				newP.FormattedTags = append(newP.FormattedTags, s)
			}
		}

		if newP.Title != "Centrifuge" {
			allProjects = append(allProjects, newP)
		}
	}

	if err != nil {
		c.HTML(http.StatusOK, "index.html", DataResponse(c, gin.H{
			"creds":   GetFirebaseClientCredentials(),
			"discord": os.Getenv("DISCORD_URL"),
		}))
		return
	} else {
		c.HTML(http.StatusOK, "index.html", DataResponse(c, gin.H{
			"creds":    GetFirebaseClientCredentials(),
			"projects": allProjects,
			"discord":  os.Getenv("DISCORD_URL"),
		}))
		return
	}
}

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", DataResponse(c, gin.H{
		"creds": GetFirebaseClientCredentials(),
	}))
}

func GetPost(c *gin.Context) {
	c.HTML(http.StatusOK, "post.html", DataResponse(c, gin.H{
		"creds": GetFirebaseClientCredentials(),
	}))
}

func GetLogout(c *gin.Context) {
	c.SetCookie("is_authed", "false", -1, "/", os.Getenv("APP_DOMAIN"), true, true)
	c.SetCookie("discord_user", "", -1, "/", os.Getenv("APP_DOMAIN"), true, true)

	c.Redirect(http.StatusFound, "/")
}

func GetUserProjects(c *gin.Context) {
	userID, _ := c.Cookie("user_id")
	id := c.Param("id")

	fmt.Println(userID)
	fmt.Println(id)
	if userID != id {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	if id == "" {
		return
	}

	// get projects
	ps, err := StoreClient.Collection("posts").
		Documents(context.Background()).
		GetAll()

	var allProjects []*UserProject
	for _, p := range ps {
		var newP *Project
		err := p.DataTo(&newP)
		if err != nil {
			// TODO
		}

		for _, ms := range newP.Members {
			if ms.ID == id {
				var needsBackend, needsFrontend, needsInfra bool
				for _, lf := range newP.LookingFor {
					if lf == "Backend" {
						needsBackend = true
					} else if lf == "Frontend" {
						needsFrontend = true
					} else if lf == "Infra" {
						needsInfra = true
					}
				}

				if newP.Title != "Centrifuge" {
					allProjects = append(allProjects, &UserProject{
						Project:            newP,
						NeedsBackend:       needsBackend,
						NeedsFrontend:      needsFrontend,
						NeedsInfra:         needsInfra,
						EnglishStartTime:   timeago.English.Format(newP.StartedAt),
						EnglishCreatedTime: timeago.English.Format(newP.CreatedAt),
					})
				}
			}
		}
	}

	if err != nil {
		c.HTML(http.StatusOK, "projects.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
		}))
		return
	} else {
		c.HTML(http.StatusOK, "projects.html", DataResponse(c, gin.H{
			"creds":    GetFirebaseClientCredentials(),
			"projects": allProjects,
		}))
		return
	}
}
