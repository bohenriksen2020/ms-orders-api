package order

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/bohenriksen2020/ms-orders-api/model"
	"fmt"
	"errors"

)

type RedisRepo struct {
	Client *redis.Client
}

// NewRedisRepo creates a new RedisRepo
func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{Client: client}
}


func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	key := orderIDKey(order.OrderID)
	
	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
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


func (r *RedisRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	key := orderIDKey(id)
	value, err := r.Client.Get(ctx, key).Result()
	fmt.Println("value: ", value)

	if errors.Is(err, redis.Nil) {
		return model.Order{}, ErrNotExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("failed to get order: %w", err)
	}
	
	var order model.Order
	err = json.Unmarshal([]byte(value), &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("failed to unmarshal order: %w", err)
	}

	fmt.Println("order: ", order)
	return order, nil

}


func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
	key := orderIDKey(id)
	txn := r.Client.TxPipeline()
	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
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


func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	// Find all order keys in Redis
	keys, err := r.Client.Keys(ctx, "order:*").Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to retrieve keys: %w", err)
	}

	var orders []model.Order

	// Loop through the keys and fetch each order
	for _, key := range keys {
		// Get the order data from Redis
		orderData, err := r.Client.Get(ctx, key).Result()

		// Check if the order data is nil or empty
		if err == redis.Nil {
			continue // Key does not exist, skip this one
		} else if err != nil {
			return FindResult{}, fmt.Errorf("failed to get order data for key %s: %w", key, err)
		}

		// Assuming you have a method to unmarshal the Redis data into an order struct
		var ord model.Order
		err = json.Unmarshal([]byte(orderData), &ord)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to unmarshal order data for key %s: %w", key, err)
		}

		// Append the order to the result set
		orders = append(orders, ord)
	}

	// Return the list of orders and the current cursor (pagination)
	return FindResult{
		Orders: orders,
		Cursor: page.Offset + uint64(len(orders)),
	}, nil
}

// Use the redis.Client's Ping method
func (r *RedisRepo) Ping(ctx context.Context) *redis.StatusCmd {
	return r.Client.Ping(ctx)
}

// Close closes the Redis connection
func (r *RedisRepo) Close() error {
	return r.Client.Close()
}