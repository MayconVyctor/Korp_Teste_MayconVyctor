package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/database"
	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/models"
)

func CreateProduct(ctx *gin.Context) {
	var newProduct models.Product

	if err := ctx.ShouldBindJSON(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload format"})
		return
	}

	insertQuery := `INSERT INTO products (code, description, balance) VALUES ($1, $2, $3)`

	_, err := database.DB.Exec(insertQuery, newProduct.Code, newProduct.Description, newProduct.Balance)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product into the database"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func GetProduct(ctx *gin.Context) {
	code := ctx.Param("code")

	var product models.Product

	selectQuery := `SELECT code, description, balance FROM products WHERE code = $1`

	err := database.DB.QueryRow(selectQuery, code).Scan(&product.Code, &product.Description, &product.Balance)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func GetAllProducts(ctx *gin.Context) {
	var products []models.Product

	selectQuery := `SELECT code, description, balance FROM products`

	rows, err := database.DB.Query(selectQuery)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products from the database"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Code, &product.Description, &product.Balance); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product data"})
			return
		}
		products = append(products, product)
	}

	ctx.JSON(http.StatusOK, products)
}
