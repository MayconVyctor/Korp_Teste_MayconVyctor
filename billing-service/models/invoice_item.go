package models

type InvoiceItem struct {
	ID                 int    `json:"id"`
	InvoiceID          int    `json:"invoice_id"`
	ProductCode        string `json:"product_code"`
	ProductDescription string `json:"product_description"`
	Quantity           int    `json:"quantity"`
}
