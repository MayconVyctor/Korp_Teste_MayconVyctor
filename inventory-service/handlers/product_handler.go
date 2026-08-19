package handlers

import (
	"net/http"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/models"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/repository"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var newProduct models.Product

	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := ph.repo.SaveProduct(newProduct)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save product to the database"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	code := ctx.Param("code")
	product, err := ph.repo.GetProductByCode(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (ph *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := ph.repo.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products from the database"})
		return
	}
	ctx.JSON(http.StatusOK, products)
}
