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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func (h *HTTPHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var item products.Product
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
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

	for k, v := range cart.Items {
		if v.UUID == item.UUID {
			tmp := cart.Items[k]
			tmp.Quantity += item.Quantity
			cart.Items[k] = tmp
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(cart)
			return
		}
	}

	cart.Items[item.UUID] = item
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
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

	for _, v := range cart.Items {
		if v.UUID == itemID {
			cart.RemoveItem(itemID)
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(cart)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}
