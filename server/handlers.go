package server

import (
	"cart-api/products"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type HTTPHandler struct {
	carts map[string]*products.Cart
	mtx   sync.RWMutex
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{
		carts: make(map[string]*products.Cart),
	}
}

func (h *HTTPHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	h.mtx.RLock()
	cart, exists := h.carts[userID]
	h.mtx.RUnlock()
	if !exists {
		cart = products.NewCartRef()
	}

	WriteJson(w, cart)
}

func (h *HTTPHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var item ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product, err := products.NewProduct(
		item.Title,
		item.Price,
		item.Quantity,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.mtx.Lock()
	defer h.mtx.Unlock()

	cart, exists := h.carts[userID]
	if !exists {
		cart = products.NewCartRef()
		h.carts[userID] = cart
	}

	cart.AddItem(product)

	WriteJson(w, cart)
}

func (h *HTTPHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	itemID := chi.URLParam(r, "item_id")

	h.mtx.Lock()
	defer h.mtx.Unlock()

	cart, exists := h.carts[userID]
	if !exists {
		http.Error(w, "Cart not found", http.StatusNotFound)
		return
	}

	err := cart.RemoveItem(itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	WriteJson(w, cart)
}

func (h *HTTPHandler) CleanCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	h.mtx.Lock()
	defer h.mtx.Unlock()

	cart, exists := h.carts[userID]
	if !exists {
		http.Error(w, "Cart not found", http.StatusBadRequest)
	}

	cart.CleanCart()

	WriteJson(w, cart)
}
