package model

import "github.com/google/uuid"

type Order struct {
	OrderID 	uint64 `json:"order_id"`
	CustomerID 	uuid.UUID `json:"customer_id"`
	LineItems 	[]LineItem `json:"line_items"`
	Created 	*time.Time `json:"created"`
	SnippedAt 	*time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`

} 

type LineItem struct {
	ItemId uuid.UUID
	Quantity uint
	Price uint

}