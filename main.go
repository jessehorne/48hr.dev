package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://localhost:8080", "https://48hr.dev"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// load html templates
	r.LoadHTMLGlob("templates/*")

	// service ./public/* to the "/" route
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir("public", false))))
	//r.GET("/css/style.css", func(c *gin.Context) {
	//	c.Header("Content-Type", "text/css")
	//	css, _ := os.ReadFile("./public/css/style.css")
	//	c.Writer.Write(css)
	//})

	r.POST("/projects", mb.DiscordAuthMiddleware, mb.PostProject)
	r.POST("/projects/:id", mb.DiscordAuthMiddleware, mb.UpdateProject)
	r.GET("/projects/:id/delete", mb.DiscordAuthMiddleware, mb.DeleteProject)
	r.GET("/projects/:id/start", mb.DiscordAuthMiddleware, mb.GetStart)
	r.GET("/projects/:id/apply/:which", mb.DiscordAuthMiddleware, mb.GetApply)
	r.GET("/projects/:id/approve/:applicantID/:applicantUsername", mb.DiscordAuthMiddleware, mb.GetApprove)
	r.GET("/projects/:id/deny/:applicantID", mb.DiscordAuthMiddleware, mb.GetDeny)
	r.GET("/projects/:id/remove/:memberID", mb.DiscordAuthMiddleware, mb.GetRemove)
	//r.GET("/projects/:id/deny/:applicantID", mb.DiscordAuthMiddleware, mb.GetDeny)
	//r.GET("/projects/:id/disable/:which", mb.DiscordAuthMiddleware, mb.GetDisable)

	r.GET("/profile", mb.DiscordAuthMiddleware, mb.GetProfile)
	r.POST("/profile", mb.DiscordAuthMiddleware, mb.PostProfile)

	// Web Routes
	r.GET("/logout", mb.GetLogout)
	r.GET("/", mb.GetIndex)
	r.GET("/login", mb.GetLogin)
	r.GET("/post", mb.GetPost)
	r.GET("/users/:id", mb.GetUserProfile)
	r.GET("/users/:id/projects", mb.GetUserProjects)
	r.GET("/auth/callback", mb.GetAuthCallback)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"data": "pong",
		})
	})

	port := os.Getenv("APP_PORT")
	key := os.Getenv("SSL_KEY")
	cert := os.Getenv("SSL_CERT")

	fmt.Println(port, key, cert)

	r.RunTLS(":"+port, cert, key) // listen and serve on 0.0.0.0:8080
}
