package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/pedrocelso/go-rest-service/lib/config"
	"github.com/pedrocelso/go-rest-service/lib/db"
)

func main() {
	cfg := config.Init()
	db := db.ConnectToDB(cfg.Mysql)

	// router := gin.Default()
	// v1 := router.Group("/api/v1/todos")
	// {
	// 	v1.POST("/", CreateTodo)
	// 	v1.GET("/", FetchAllTodo)
	// 	v1.GET("/:id", FetchSingleTodo)
	// 	v1.PUT("/:id", UpdateTodo)
	// 	v1.DELETE("/:id", DeleteTodo)
	// }
	// router.Run()

	fmt.Println("\n\n================config.MustInit()====================")
	spew.Dump(cfg)
	spew.Dump(db)
	fmt.Println("====================================")
}
