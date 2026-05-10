package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	handlers *HTTPHandler
}

func simpleUserIDGenerator() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func NewServer(handlers *HTTPHandler) *Server {
	return &Server{
		handlers: handlers,
	}
}

func (s *Server) StartServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string

			if cookie, err := r.Cookie("user_id"); err == nil {
				userID = cookie.Value
			}

			if userID == "" {
				userID = simpleUserIDGenerator()
				http.SetCookie(w, &http.Cookie{
					Name:   "user_id",
					Value:  userID,
					Path:   "/",
					MaxAge: 3600,
				})
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Get("/carts", s.handlers.GetCartHandler)
	r.Post("/carts/items", s.handlers.AddProductHandler)
	r.Delete("/carts/items/{item_id}", s.handlers.DeleteProductHandler)

	http.ListenAndServe(":3000", r)
}
