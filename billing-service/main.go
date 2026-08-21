package main

import (
	"log"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/db"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/handlers"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to billing database:", err)
	}
	defer dbConn.Close()

	invoiceRepo := repository.NewInvoiceRepository(dbConn)
	invoiceRepo.CreateInvoice()
	invoiceHandler := handlers.NewInvoiceHandler(invoiceRepo)

	server := gin.Default()
	server.Use(cors.Default())

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Billing Service is alive!"})
	})
	server.GET("/invoices", invoiceHandler.GetAllInvoices)
	server.POST("/invoices", invoiceHandler.CreateInvoice)
	server.PUT("/invoices/:id/print", invoiceHandler.PrintInvoice)
	server.GET("/invoices/:id/analysis", invoiceHandler.AnalyzeInvoiceHandler)

	log.Println("Billing Service running on port 8082...")
	server.Run(":8082")
}
