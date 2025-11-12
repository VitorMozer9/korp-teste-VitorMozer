package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/domain"
)

// StockHTTPClient implementa StockClient para comunicação HTTP com Stock Service
type StockHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewStockHTTPClient(baseURL string) *StockHTTPClient {
	return &StockHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ReservationRequest representa a estrutura de reserva do Stock Service
type ReservationRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// ReservationResponse representa a resposta do Stock Service
type ReservationResponse struct {
	Success      bool   `json:"success"`
	ProductID    string `json:"product_id"`
	NewBalance   int    `json:"new_balance"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// ReserveProducts reserva múltiplos produtos no Stock Service
func (c *StockHTTPClient) ReserveProducts(items []domain.InvoiceItem) error {
	// Converte InvoiceItems para ReservationRequests
	requests := make([]ReservationRequest, len(items))
	for i, item := range items {
		requests[i] = ReservationRequest{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Serializa o payload
	payload, err := json.Marshal(requests)
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Faz a requisição HTTP
	url := fmt.Sprintf("%s/api/products/reserve", c.baseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao comunicar com Stock Service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Stock Service retornou erro: status %d", resp.StatusCode)
	}

	var responses []ReservationResponse
	if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
		return fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Verifica se todas as reservas foram bem-sucedidas
	for _, res := range responses {
		if !res.Success {
			return fmt.Errorf("falha ao reservar produto %s: %s", res.ProductID, res.ErrorMessage)
		}
	}

	return nil
}

//verifica disponibilidade de um produto
func (c *StockHTTPClient) CheckAvailability(productID string, quantity int) (bool, error) {
	// Busca informações do produto
	product, err := c.GetProduct(productID)
	if err != nil {
		return false, err
	}

	// Verifica se há saldo suficiente
	return product.Balance >= quantity, nil
}

// busca informações de um produto
func (c *StockHTTPClient) GetProduct(productID string) (*domain.ProductInfo, error) {
	url := fmt.Sprintf("%s/api/products/%s", c.baseURL, productID)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao comunicar com Stock Service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("produto não encontrado")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Stock Service retornou erro: status %d", resp.StatusCode)
	}

	var product domain.ProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &product, nil
}