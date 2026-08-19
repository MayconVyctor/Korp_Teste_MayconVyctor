package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/inventory-service/models"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{connection: db}
}

func (pr *ProductRepository) CreateProduct() {
	createTableProductQuery := "CREATE TABLE IF NOT EXISTS products (code VARCHAR(50) PRIMARY KEY, description VARCHAR(255) NOT NULL, balance INT NOT NULL);"
	_, err := pr.connection.Exec(createTableProductQuery)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Products table verified/created successfully!")
}

func (pr *ProductRepository) SaveProduct(product models.Product) error {
	insertQuery := "INSERT INTO products (code, description, balance) VALUES ($1, $2, $3)"

	_, err := pr.connection.Exec(insertQuery, product.Code, product.Description, product.Balance)
	if err != nil {
		return fmt.Errorf("failed to insert product into the database: %v", err)
	}

	return nil
}

func (pr *ProductRepository) GetProductByCode(code string) (models.Product, error) {
	var product models.Product

	selectQuery := "SELECT code, description, balance FROM products WHERE code = $1"
	err := pr.connection.QueryRow(selectQuery, code).Scan(&product.Code, &product.Description, &product.Balance)
	if err != nil {
		return product, fmt.Errorf("product not found: %v", err)
	}

	return product, err
}

func (pr *ProductRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product

	Query := "SELECT code, description, balance FROM products"
	rows, err := pr.connection.Query(Query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products from the database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Code, &product.Description, &product.Balance); err != nil {
			return nil, fmt.Errorf("failed to scan product row: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %v", err)
	}

	return products, nil
}
