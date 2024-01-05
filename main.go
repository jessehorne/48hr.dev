package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/jessehorne/microblog/mb"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	
	debug := true
	if os.Getenv("APP_DEBUG") == "true" {
		debug = true
	} else {
		debug = false
	}
	mb.InitApp(debug)

	r := gin.Default()

	// HTML Templates
	r.LoadHTMLGlob("templates/*")
	
	// Public files
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("public", false))))

	//// Auth'd API routes
	var api *gin.RouterGroup

	api = r.Group("/api")
	
	api.Use(mb.FirebaseAuthMiddleware)
	api.POST("/projects", mb.PostProject)
	api.GET("/projects/:id/apply/:which", func(c *gin.Context) {
		// get user id from request
		token, exists := c.Get("token")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "not authenticated",
			})
			return
		}

		t := token.(*auth.Token)
		if t == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "invalid token",
			})
			return
		}

		userID := t.UID

		u, err := mb.AuthClient.GetUser(context.Background(), userID)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		which := c.Param("which")
		id := c.Param("id")
		
		newApplicant := mb.Applicant{
			ID:          userID,
			DisplayName: u.DisplayName,
			Which: which,
		}
		
		posts, err := mb.StoreClient.Collection("posts").Where("ProjectID", "==", id).Documents(context.Background()).GetAll()
		if err != nil {
			fmt.Println(err)
			return
		}
		
		for _, p := range posts {
			var newP *mb.Project
			p.DataTo(&newP)
			newP.Applicants = append(newP.Applicants, newApplicant)
			
			po := mb.StoreClient.Collection("posts").Doc(p.Ref.ID)
			_, err := po.Update(context.Background(), []firestore.Update{
				{Path: "Applicants", Value: newP.Applicants},
			})
			if err != nil {
				fmt.Println(err)
			}
		}
		
		c.Redirect(http.StatusFound, "/")
	})

	// Web Routes
	r.GET("/", func(c *gin.Context) {
		// get projects
		ps, err := mb.StoreClient.Collection("posts").Documents(context.Background()).GetAll()
		
		var allProjects []*mb.Project
		for _, p := range ps {
			var newP *mb.Project
			err := p.DataTo(&newP)
			if err != nil {
				fmt.Println(err)
			}
			
			if newP.Title != "Centrifuge" {
				allProjects = append(allProjects, newP)	
			}
		}
		
		if err != nil {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"creds": mb.GetFirebaseClientCredentials(),
			})
			return
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"creds": mb.GetFirebaseClientCredentials(),
				"projects": allProjects,
			})
			return
		}
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"creds": mb.GetFirebaseClientCredentials(),
		})
	})

	r.GET("/post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", gin.H{
			"creds": mb.GetFirebaseClientCredentials(),
		})
	})
	
	r.GET("/users/:id/projects", func(c *gin.Context) {
		id := c.Param("id")
		
		if id == "" {
			return
		}
		
		// get projects
		ps, err := mb.StoreClient.Collection("posts").
			Where("UserID", "==", id).
			Documents(context.Background()).
			GetAll()

		var allProjects []*mb.Project
		for _, p := range ps {
			var newP *mb.Project
			err := p.DataTo(&newP)
			if err != nil {
				fmt.Println(err)
			}

			if newP.Title != "Centrifuge" {
				allProjects = append(allProjects, newP)
			}
		}

		if err != nil {
			c.HTML(http.StatusOK, "projects.html", gin.H{
				"creds": mb.GetFirebaseClientCredentials(),
			})
			return
		} else {
			c.HTML(http.StatusOK, "projects.html", gin.H{
				"creds": mb.GetFirebaseClientCredentials(),
				"projects": allProjects,
			})
			return
		}
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
