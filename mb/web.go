package mb

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetIndex(c *gin.Context) {
	// get projects
	ps, err := StoreClient.Collection("posts").Documents(context.Background()).GetAll()

	var allProjects []*Project
	for _, p := range ps {
		var newP *Project
		err := p.DataTo(&newP)
		if err != nil {
			// TODO
		}

		if newP.Title != "Centrifuge" {
			allProjects = append(allProjects, newP)
		}
	}

	if err != nil {
		c.HTML(http.StatusOK, "index.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
			"discord": os.Getenv("DISCORD_URL"),
		}))
		return
	} else {
		c.HTML(http.StatusOK, "index.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
			"projects": allProjects,
			"discord": os.Getenv("DISCORD_URL"),
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
	c.SetCookie("is_authed", "false", -1, "/", "localhost", true, true)
	c.SetCookie("discord_user", "", -1, "/", "localhost", true, true)

	c.Redirect(http.StatusFound, "/")
}

func GetUserProjects(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		return
	}

	// get projects
	ps, err := StoreClient.Collection("posts").
		Where("UserID", "==", id).
		Documents(context.Background()).
		GetAll()

	var allProjects []*UserProject
	for _, p := range ps {
		var newP *Project
		err := p.DataTo(&newP)
		if err != nil {
			// TODO
		}
		
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
				Project: newP,
				NeedsBackend: needsBackend,
				NeedsFrontend: needsFrontend,
				NeedsInfra: needsInfra,
			})
		}
	}

	if err != nil {
		c.HTML(http.StatusOK, "projects.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
		}))
		return
	} else {
		c.HTML(http.StatusOK, "projects.html", DataResponse(c, gin.H{
			"creds": GetFirebaseClientCredentials(),
			"projects": allProjects,
		}))
		return
	}
}
