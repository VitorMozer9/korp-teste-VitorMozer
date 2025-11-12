package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/repo/mem"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/usecase"
	httpTransport "github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/transport/http"
)

func main() {
	port := getEnv("PORT", "8081")

	// (Dependency Injection)
	// Repository -> UseCase -> Handler -> Router
	productRepo := mem.NewProductMemRepository()
	productService := usecase.NewProductService(productRepo)
	handler := httpTransport.NewHandler(productService)
	router := httpTransport.NewRouter(handler)

	// Servidor HTTP
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Inicia o servidor em uma goroutine
	go func() {
		log.Printf("Stock Service iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Desligamento normal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando Stock Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown: %v", err)
	}

	log.Println("Stock Service desligado com sucesso")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
