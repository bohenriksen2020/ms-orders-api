package handler
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"github.com/google/uuid"
	
	"github.com/bohenriksen2020/ms-orders-api/model"
	"github.com/bohenriksen2020/ms-orders-api/repository/order"
	"strconv"

)


type Order struct{
	Repo *order.RedisRepo
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID 		`json:"customer_id"`
		LineItems []model.LineItem 	`json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode request body: ", err)
		w.Write([]byte("failed to decode request body: " + err.Error() + "\n"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	order := model.Order{
		OrderID: rand.Uint64(), // Don't do this in production
		CustomerID: body.CustomerID,
		LineItems: body.LineItems,
		Created: &now,
	}

	err := o.Repo.Insert(r.Context(), order)
	if err != nil {
		fmt.Println("failed to insert order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	if err != nil {
		fmt.Println("failed to marshal order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		fmt.Println("failed to parse cursor: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	const size = 50
	res, err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size: size,

	})
	if err != nil {
		fmt.Println("failed to find all orders: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var response struct {
		Items []model.Order `json:"items"`
		Next uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal orders: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)	
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get order by Id")
	w.WriteHeader(http.StatusOK)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update order by ID")
	w.WriteHeader(http.StatusOK)
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete order by ID")
	w.WriteHeader(http.StatusOK)
}

