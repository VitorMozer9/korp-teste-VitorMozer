package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/client"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/repo/mem"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/usecase"
	httpTransport "github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/transport/http"
)

func main() {
	port := getEnv("PORT", "8082")
	stockServiceURL := getEnv("STOCK_SERVICE_URL", "http://localhost:8081")

	log.Printf("Configurando Billing Service...")
	log.Printf("   - Porta: %s", port)
	log.Printf("   - Stock Service URL: %s", stockServiceURL)

	// Inicialização das camadas (Dependency Injection)
	// Client -> Repository -> UseCase -> Handler -> Router
	stockClient := client.NewStockHTTPClient(stockServiceURL)
	invoiceRepo := mem.NewInvoiceMemRepository()
	invoiceService := usecase.NewInvoiceService(invoiceRepo, stockClient)
	handler := httpTransport.NewHandler(invoiceService)
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
		log.Printf("Billing Service iniciado na porta %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	//Desligamento normal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Desligando Billing Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown: %v", err)
	}

	log.Println("Billing Service desligado com sucesso")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}