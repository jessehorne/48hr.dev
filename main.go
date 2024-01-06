package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jessehorne/microblog/mb"
	"github.com/joho/godotenv"
)

func main() {
	// make environment variables accessible through os.Getenv(...)
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}


	// Initialize the state of the app that holds info like firebase/firestore/validator/etc
	mb.InitApp(false)

	r := gin.Default()

	// load html templates
	r.LoadHTMLGlob("templates/*")
	
	// service ./public/* to the "/" route
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("public", false))))

	// Auth'd API routes
	var api *gin.RouterGroup

	api = r.Group("/api")
	
	//api.Use(mb.FirebaseAuthMiddleware)
	api.Use(mb.DiscordAuthMiddleware)
	api.POST("/projects", mb.PostProject)
	api.GET("/projects/:id/apply/:which", mb.GetApply)
	
	r.GET("/profile", mb.DiscordAuthMiddleware, mb.GetProfile)
	r.POST("/profile", mb.DiscordAuthMiddleware, mb.PostProfile)

	// Web Routes
	r.GET("/logout", mb.GetLogout)
	r.GET("/", mb.GetIndex)
	r.GET("/login", mb.GetLogin)
	r.GET("/post", mb.GetPost)
	r.GET("/users/:id/projects", mb.GetUserProjects)
	r.GET("/auth/callback", mb.GetAuthCallback)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})

	r.RunTLS("localhost:8080", "./ssl/server-cert.pem", "./ssl/server-key.pem") // listen and serve on 0.0.0.0:8080
}
