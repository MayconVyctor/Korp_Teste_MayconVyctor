package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/database"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/handlers"
)

func main() {

	database.InitDB()
	database.CreateTable()

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.POST("/products", handlers.CreateProduct)
	server.GET("/products/:code", handlers.GetProduct)
	server.GET("/products", handlers.GetAllProducts)
	server.Run(":8081")
}
