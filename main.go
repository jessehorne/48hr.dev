package main

import (
	"net/http"

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
	mb.InitApp(debug)

	r := gin.Default()

	// HTML Templates
	r.LoadHTMLGlob("templates/*")
	
	// Public files
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("public", false))))

	//// Auth'd API routes
	//var api *gin.RouterGroup
	//
	//if debug {
	//	api = r.Group("/api")
	//} else {
	//	api = r.Group("/api", mb.FirebaseAuthMiddleware)
	//}

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

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
