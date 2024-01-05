package main

import (
	"net/http"
	"os"

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

	// Web Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"creds": mb.GetFirebaseClientCredentials(),
		})
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

	r.GET("/projects", func(c *gin.Context) {
		c.HTML(http.StatusOK, "projects.html", gin.H{
			"creds": mb.GetFirebaseClientCredentials(),
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
