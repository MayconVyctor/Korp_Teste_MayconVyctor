package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/MayconVyctor/Korp_Teste_MayconVyctor/billing-service/models"
)

type InvoiceRepository struct {
	connection *sql.DB
}

func NewInvoiceRepository(db *sql.DB) *InvoiceRepository {
	return &InvoiceRepository{connection: db}
}

func (ir *InvoiceRepository) CreateInvoice() {
	createTableInvoiceQuery := "CREATE TABLE IF NOT EXISTS invoices (id SERIAL PRIMARY KEY, status VARCHAR(50) NOT NULL);"
	_, err := ir.connection.Exec(createTableInvoiceQuery)
	if err != nil {
		log.Fatal("Failed to create invoices table:", err)
	}

	createTableInvoiceItemQuery := "CREATE TABLE IF NOT EXISTS invoice_items (id SERIAL PRIMARY KEY, invoice_id INT REFERENCES invoices(id), product_code VARCHAR(50) NOT NULL, quantity INT NOT NULL);"
	_, err = ir.connection.Exec(createTableInvoiceItemQuery)
	if err != nil {
		log.Fatal("Failed to create invoice_items table:", err)
	}

	fmt.Println("Invoices and Invoice Items tables verified/created successfully!")
}

func (ir *InvoiceRepository) SaveInvoice(invoice models.Invoice) error {

	tx, err := ir.connection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer tx.Rollback()

	insertInvoiceQuery := "INSERT INTO invoices (status) VALUES ($1) RETURNING id"
	var invoiceID int
	err = tx.QueryRow(insertInvoiceQuery, invoice.Status).Scan(&invoiceID)
	if err != nil {
		return fmt.Errorf("failed to insert invoice into the database: %v", err)
	}

	for _, item := range invoice.Items {
		insertInvoiceItemQuery := "INSERT INTO invoice_items (invoice_id, product_code, quantity) VALUES ($1, $2, $3)"
		_, err := tx.Exec(insertInvoiceItemQuery, invoiceID, item.ProductCode, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to insert invoice item into the database: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	fmt.Println("Invoice and its items saved successfully!")

	return nil
}
