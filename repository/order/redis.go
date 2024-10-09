package order

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/bohenriksen2020/ms-orders-api/model"
	"fmt"
	"errors"

)

type struct RedisRepo {
	Client *redis.Client
}

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	date, erro := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(order.OrderID)
	
	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, data, 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to insert order: %w", err)
	}

	if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add order to set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec transaction: %w", err)
	}

	return nil
}


var ErrNotExist = errors.New("order does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (*model.Order, error) {
	key := orderIDKey(id)
	res := r.Client.Get(ctx, key).Result()

	if errors.Is(res.Err(), redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	
	var order model.Order
	err = json.Unmarshal([]byte(res), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	return order, nil

}


func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	txn := r.Client.TxPipeline()
	res := txn.Del(ctx, key).Err()
	if errors.Is(res, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("failed to delete order: %w", err)
	}

	if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove order from set: %w", err)
	}
	
	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec transaction: %w", err)
	}

	return nil
}


func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(order.OrderID)
	res := r.Client.Set(ctx, key, string(data), 0).Err()
	if errors.Is(res, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	return nil
}

type  FindAllPage struct {
	Size uint
	Offset uint
}

type FindResult struct {
	Orders []model.Order
	Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context) (FindResult, error) {
	res := r.Client.SScan(ctx, "order", page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to scan keys: %w", err)
	}

	if len(leys) == 0 {
		return FindResult{
			Orders: []model.Order{},
		}, nil
	}

	xs, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]model.Order, 0, len(xs))

	for i, x := range xs {
		x := x.(string)
		var order model.Order
		err := json.Unmarshal([]byte(x), &order)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to unmarshal order: %w", err)
		}
		orders[i] = order
	}

	return FindResult{
		Orders: orders,
		Cursor: cursor,

	}, nil
}