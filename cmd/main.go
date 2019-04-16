package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-task/lib/http/controllers"
)

func main() {
	router := gin.New()
	router.Use(controllers.CORSMiddleware())
	router.Use(controllers.CheckJWT(os.Getenv("JWT_SECRET")))
	router.OPTIONS("/", func(c *gin.Context) {})

	v1 := router.Group("/v1")

	users := v1.Group("/users")
	users.POST("/", controllers.CreateUser)
	users.GET("/:userEmail", controllers.GetUser)
	users.GET("/", controllers.GetUsers)
	users.PUT("/:userEmail", controllers.UpdateUser)
	users.DELETE("/:userEmail", controllers.DeleteUser)

	tasks := v1.Group("/tasks")
	tasks.POST("/", controllers.CreateTask)
	tasks.GET("/:taskId", controllers.GetTask)
	tasks.GET("/", controllers.GetTasks)
	tasks.PUT("/:taskId", controllers.UpdateTask)
	tasks.DELETE("/:taskId", controllers.DeleteTask)

	http.Handle("/", router)
}
