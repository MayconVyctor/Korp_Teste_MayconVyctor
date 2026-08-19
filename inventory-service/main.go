package main

import (
	"github.com/gin-gonic/gin"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/database"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/handlers"
)

func main() {

	database.InitDB()
	database.CreateTable()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/products", handlers.CreateProduct)
	r.Run(":8081")
}
