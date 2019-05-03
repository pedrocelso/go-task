package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-task/lib/http/controllers"
	"google.golang.org/appengine"
)

func main() {
	router := gin.New()
	router.Use(controllers.CORSMiddleware())
	router.OPTIONS("/", func(c *gin.Context) {})

	v1 := router.Group("/v1")

	public := v1.Group("/public")
	public.POST("/signin", controllers.AuthenticateUser)

	users := v1.Group("/users")
	users.Use(controllers.CheckJWT(os.Getenv("JWT_SECRET")))
	users.POST("/", controllers.CreateUser)
	users.GET("/:userEmail", controllers.GetUser)
	users.GET("/", controllers.GetUsers)
	users.PUT("/:userEmail", controllers.UpdateUser)
	users.DELETE("/:userEmail", controllers.DeleteUser)

	tasks := v1.Group("/tasks")
	tasks.Use(controllers.CheckJWT(os.Getenv("JWT_SECRET")))
	tasks.POST("/", controllers.CreateTask)
	tasks.GET("/:taskId", controllers.GetTask)
	tasks.GET("/", controllers.GetTasks)
	tasks.PUT("/:taskId", controllers.UpdateTask)
	tasks.DELETE("/:taskId", controllers.DeleteTask)

	router.Run()
	appengine.Main()
}
