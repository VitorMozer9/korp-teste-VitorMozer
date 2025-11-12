package http

import (
	"net/http"
	"time"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// cria e configura o roteador HTTP
func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) 
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS para permitir requisições do frontend 
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", handler.Health)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/invoices", func(r chi.Router) {
			r.Get("/", handler.GetAllInvoices)
			r.Post("/", handler.CreateInvoice)
			r.Get("/{id}", handler.GetInvoice)
			
			// Endpoint de impressão (fechamento) da nota fiscal
			r.Post("/{id}/print", handler.PrintInvoice)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"endpoint não encontrado"}`))
	})

	return r
}