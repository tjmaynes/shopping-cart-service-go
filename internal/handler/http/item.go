package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"

	jsonHandler "github.com/tjmaynes/shopping-cart-service-go/internal/handler/json"
	cart "github.com/tjmaynes/shopping-cart-service-go/internal/pkg/item"
)

// NewItemHandler ..
func NewItemHandler(service cart.Service) *ItemHandler {
	return &ItemHandler{Service: service}
}

// ItemHandler ..
type ItemHandler struct {
	Service cart.Service
}

// GetItems ..
func (c *ItemHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.ParseInt(r.URL.Query().Get("pageSize"), 10, 64)
	if err != nil {
		pageSize = 10
	}

	data, err := c.Service.GetItems(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonHandler.CreateResponse(w, http.StatusOK, map[string][]cart.Item{"data": data})
}

// GetItemByID ..
func (c *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id, errorCode := getID(r.URL.Path)
	if errorCode >= 400 {
		http.Error(w, http.StatusText(errorCode), errorCode)
	}

	data, err := c.Service.GetItemByID(r.Context(), uuid.MustParse(*id))
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	jsonHandler.CreateResponse(w, http.StatusOK, map[string]cart.Item{"data": data})
}

// AddItem ..
func (c *ItemHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	itemName := r.Form.Get("name")
	itemManufacturer := r.Form.Get("manufacturer")
	itemPrice, errorCode := getItemPrice(r.Form.Get("price"))
	if errorCode >= 400 {
		http.Error(w, http.StatusText(errorCode), errorCode)
		return
	}

	item := cart.ItemDTO{Name: itemName, Price: itemPrice, Manufacturer: itemManufacturer}

	data, err := c.Service.AddItem(r.Context(), &item)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	jsonHandler.CreateResponse(w, http.StatusCreated, map[string]cart.Item{"data": data})
}

// UpdateItem ..
func (c *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	rawId, errorCode := getID(r.URL.Path)
	if errorCode >= 400 {
		http.Error(w, http.StatusText(errorCode), errorCode)
	}

	decoder := json.NewDecoder(r.Body)
	type RawItemRequest struct {
		Name         string `json:"name"`
		Price        string `json:"price"`
		Manufacturer string `json:"manufacturer"`
	}
	var rawItemRequest RawItemRequest
	err := decoder.Decode(&rawItemRequest)
	if err != nil {
		panic(err)
	}

	price, errorCode := getItemPrice(rawItemRequest.Price)
	if errorCode >= 400 {
		http.Error(w, http.StatusText(errorCode), errorCode)
		return
	}

	id := uuid.MustParse(*rawId)

	result, err := c.Service.GetItemByID(r.Context(), id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if id != result.ID {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
		return
	}

	item := cart.Item{
		ID:           id,
		Name:         rawItemRequest.Name,
		Price:        price,
		Manufacturer: rawItemRequest.Manufacturer,
	}

	result, serviceError := c.Service.UpdateItem(r.Context(), &item)
	if serviceError != nil {
		switch serviceError.StatusCode() {
		case cart.InvalidItem:
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	jsonHandler.CreateResponse(w, http.StatusOK, map[string]cart.Item{"data": result})
}

// RemoveItem ..
func (c *ItemHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	id, errorCode := getID(r.URL.Path)
	if errorCode >= 400 {
		http.Error(w, http.StatusText(errorCode), errorCode)
	}

	_, serviceError := c.Service.RemoveItem(r.Context(), uuid.MustParse(*id))
	if serviceError != nil {
		http.Error(w, serviceError.Message(), 500)
		return
	}

	jsonHandler.CreateResponse(w, http.StatusOK, http.StatusText(200))
}

func getID(urlPath string) (*string, int) {
	params := strings.Split(urlPath, "/")
	if len(params) < 2 {
		return nil, http.StatusBadRequest
	}

	return &params[2], 0
}

func getItemPrice(rawPrice string) (cart.Decimal, int) {
	result, err := strconv.ParseInt(rawPrice, 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest
	}
	return cart.Decimal(result), 0
}
