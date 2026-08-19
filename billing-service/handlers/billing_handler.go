package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/clients"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/models"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/repository"
)

type InvoiceHandler struct {
	repo *repository.InvoiceRepository
}

func NewInvoiceHandler(repo *repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{repo: repo}
}

func (hand *InvoiceHandler) CreateInvoice(ctx *gin.Context) {
	var invoice models.Invoice

	if err := ctx.ShouldBindJSON(&invoice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Formato de payload inválido"})
		return
	}

	if len(invoice.Items) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The invoice must contain at least one item"})
		return
	}

	for _, item := range invoice.Items {
		err := clients.CheckProductExists(item.ProductCode)

		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   "Failed to validate product in inventory",
				"details": err.Error(),
			})
			return
		}
	}

	err := hand.repo.SaveInvoice(invoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error while saving the invoice"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "Invoice created successfully"})
}

func (h *InvoiceHandler) AnalyzeInvoiceHandler(ctx *gin.Context) {

	idParam := ctx.Param("id")

	invoiceID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID format. Must be an integer."})
		return
	}

	invoice, err := h.repo.GetInvoiceByID(invoiceID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	analysisText, err := clients.AnalyzeInvoice(invoice)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to analyze invoice with AI",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"invoice_id":  invoice.ID,
		"status":      invoice.Status,
		"ai_analysis": analysisText,
	})
}
