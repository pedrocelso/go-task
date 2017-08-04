package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pedrocelso/go-rest-service/lib/controllers"
)

func init() {
	router := gin.New()

	v1 := router.Group("/api/v1/users")
	v1.POST("/", controllers.CreateUser)
	v1.GET("/:userEmail", controllers.GetUser)
	v1.GET("/", controllers.GetUsers)
	v1.PUT("/:userEmail", controllers.UpdateUser)
	v1.DELETE("/:userEmail", controllers.DeleteUser)

	http.Handle("/", router)
}
