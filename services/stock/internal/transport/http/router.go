package http

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter cria e configura o roteador HTTP
func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // Recupera de panics
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS para permitir requisições do frontend Angular
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
		r.Route("/products", func(r chi.Router) {
			r.Get("/", handler.GetAllProducts)
			r.Post("/", handler.CreateProduct)
			r.Get("/{id}", handler.GetProduct)
			r.Put("/{id}", handler.UpdateProduct)
			r.Delete("/{id}", handler.DeleteProduct)

			// Endpoint para reserva de estoque (chamado pelo Billing Service)
			r.Post("/reserve", handler.ReserveStock)
		})
	})

	// 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"endpoint não encontrado"}`))
	})

	return r
}

// LoggingMiddleware personalizado (exemplo de middleware customizado)
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}
