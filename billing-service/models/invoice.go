package models

type Invoice struct {
	ID     int           `json:"id"`
	Status string        `json:"status"`
	Items  []InvoiceItem `json:"items"`
}
