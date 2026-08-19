package main

import (
	"log"

	"github.com/gin-gonic/gin"

	// Lembre-se de ajustar o caminho base de acordo com o seu go.mod
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/db"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/handlers"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/repository"
)

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Could not connect to billing database:", err)
	}
	defer dbConn.Close()

	invoiceRepo := repository.NewInvoiceRepository(dbConn)
	invoiceRepo.CreateInvoice()
	invoiceHandler := handlers.NewInvoiceHandler(invoiceRepo)

	server := gin.Default()

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Billing Service is alive!"})
	})

	server.POST("/invoices", invoiceHandler.CreateInvoice)

	log.Println("Billing Service running on port 8082...")
	server.Run(":8082")
}
