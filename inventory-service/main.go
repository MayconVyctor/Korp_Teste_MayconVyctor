package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/db"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/handlers"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/repository"
)

func main() {

	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	ProductRepository := repository.NewProductRepository(dbConnection)
	ProductHandler := handlers.NewProductHandler(ProductRepository)

	ProductRepository.CreateProduct()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.POST("/products", ProductHandler.CreateProduct)
	server.GET("/products/:code", ProductHandler.GetProduct)
	server.GET("/products", ProductHandler.GetAllProducts)
	server.Run(":8081")
}
