package order

import (
	"context"
	"github.com/bohenriksen2020/ms-orders-api/model"
)

// Repo is the interface that defines the methods a repository should implement
type Repo interface {
	Insert(ctx context.Context, order model.Order) error
	FindByID(ctx context.Context, id uint64) (model.Order, error)
	Update(ctx context.Context, order model.Order) error
	DeleteByID(ctx context.Context, id uint64) error
	FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
}
