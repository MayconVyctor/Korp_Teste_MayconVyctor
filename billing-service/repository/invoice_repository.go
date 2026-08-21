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

func (ir *InvoiceRepository) GetInvoiceByID(id int) (models.Invoice, error) {
	var invoice models.Invoice
	invoice.ID = id
	queryInvoice := "SELECT status FROM invoices WHERE id = $1"
	err := ir.connection.QueryRow(queryInvoice, id).Scan(&invoice.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return invoice, fmt.Errorf("invoice with ID %d not found", id)
		}
		return invoice, fmt.Errorf("failed to retrieve invoice header: %v", err)
	}

	queryItems := "SELECT id, product_code, quantity FROM invoice_items WHERE invoice_id = $1"
	rows, err := ir.connection.Query(queryItems, id)
	if err != nil {
		return invoice, fmt.Errorf("failed to retrieve invoice items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.InvoiceItem

		if err := rows.Scan(&item.ID, &item.ProductCode, &item.Quantity); err != nil {
			return invoice, fmt.Errorf("failed to scan item data: %v", err)
		}

		item.InvoiceID = id
		invoice.Items = append(invoice.Items, item)
	}
	if err := rows.Err(); err != nil {
		return invoice, fmt.Errorf("error during iteration of invoice items: %v", err)
	}

	return invoice, nil
}
func (ir *InvoiceRepository) GetAllInvoices() ([]models.Invoice, error) {
	var invoices []models.Invoice

	queryInvoices := "SELECT id, status FROM invoices"
	rows, err := ir.connection.Query(queryInvoices)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve invoices: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invoice models.Invoice
		if err := rows.Scan(&invoice.ID, &invoice.Status); err != nil {
			return nil, fmt.Errorf("failed to scan invoice data: %v", err)
		}

		queryItems := "SELECT id, product_code, quantity FROM invoice_items WHERE invoice_id = $1"
		rows, err := ir.connection.Query(queryItems, invoice.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve invoice items: %v", err)
		}

		for rows.Next() {
			var item models.InvoiceItem
			if err := rows.Scan(&item.ID, &item.ProductCode, &item.Quantity); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan item data: %v", err)
			}
			item.InvoiceID = invoice.ID
			invoice.Items = append(invoice.Items, item)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error occurred during row iteration: %v", err)
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration of invoices: %v", err)
	}

	return invoices, nil
}

func (ir *InvoiceRepository) UpdateInvoiceStatus(id int, status string) error {
	updateQuery := "UPDATE invoices SET status = $1 WHERE id = $2"
	result, err := ir.connection.Exec(updateQuery, status, id)
	if err != nil {
		return fmt.Errorf("failed to update invoice status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no invoice found with ID %d", id)
	}

	return nil
}
