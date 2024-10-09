package order

import "github.com/bohenriksen2020/ms-orders-api/model"

// FindAllPage defines pagination options for retrieving orders
type FindAllPage struct {
	Size   uint64
	Offset uint64
}

// FindResult holds the result of a paginated query, with the orders and a cursor for the next page
type FindResult struct {
	Orders []model.Order
	Cursor uint64
}
