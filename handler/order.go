package handler

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/bohenriksen2020/ms-orders-api/model"
	"github.com/bohenriksen2020/ms-orders-api/repository/order"
)

type Order struct {
	Repo order.Repo // *order.RedisRepo
}

func (h *Order) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerID uuid.UUID        `json:"customer_id"`
		LineItems  []model.LineItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode request body: ", err)
		if _, err = w.Write([]byte("failed to decode request body: " + err.Error() + "\n")); err != nil {
			fmt.Println("failed to write response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	orderID := make([]byte, 8)
	if _, err := rand.Read(orderID); err != nil {
		// handle error
		fmt.Println("failed to generate order id: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ord := model.Order{
		OrderID:    binary.BigEndian.Uint64(orderID),
		CustomerID: body.CustomerID,
		LineItems:  body.LineItems,
		Created:    &now,
	}

	err := h.Repo.Insert(r.Context(), ord)
	if err != nil {
		fmt.Println("failed to insert order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(ord)
	if err != nil {
		fmt.Println("failed to marshal order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(res); err != nil {
		fmt.Println("failed to write response: ", err)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Order) List(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.Repo.FindAll(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		fmt.Println("failed to find all orders: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println("len(res): ", len(res.Orders))

	var response struct {
		Items []model.Order `json:"items"`
		Next  uint64        `json:"next,omitempty"`
	}
	response.Items = res.Orders
	response.Next = res.Cursor
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed to marshal orders: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(data); err != nil {
		fmt.Println("failed to write response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		fmt.Println("failed to parse order id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("idParam: ", orderID)

	theOrder, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		fmt.Println("order does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// unmarshal order and return it
	res, err := json.Marshal(theOrder)
	if err != nil {
		fmt.Println("failed to marshal order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(res); err != nil {
		fmt.Println("failed to write response: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed to decode requestbody: ", err)
		if _, err = w.Write([]byte("failed to decode request body: " + err.Error() + "\n")); err != nil {
			fmt.Println("failed to write response: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		fmt.Println("failed to parse order id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	theOrder, err := h.Repo.FindByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		fmt.Println("order does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	const completedStatus = "completed"
	const shippedStatus = "shipped"

	now := time.Now().UTC()
	switch body.Status {
	case shippedStatus:
		if theOrder.ShippedAt != nil {
			fmt.Println("order already shipped")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.ShippedAt = &now
	case completedStatus:
		if theOrder.CompletedAt != nil || theOrder.ShippedAt == nil {
			fmt.Println("order already completed")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		theOrder.CompletedAt = &now
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.Update(r.Context(), theOrder)
	if err != nil {
		fmt.Println("failed to update order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(theOrder); err != nil {
		fmt.Println("failed to encode order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {

	idParm := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	orderID, err := strconv.ParseUint(idParm, base, bitSize)
	if err != nil {
		fmt.Println("failed to parse order id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteByID(r.Context(), orderID)
	if errors.Is(err, order.ErrNotExist) {
		fmt.Println("order does not exist")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to delete order: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
