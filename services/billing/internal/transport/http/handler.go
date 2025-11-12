package http

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/domain"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/usecase"
)

type Handler struct {
	invoiceService *usecase.InvoiceService
}

// NewHandler cria um novo handler
func NewHandler(invoiceService *usecase.InvoiceService) *Handler {
	return &Handler{
		invoiceService: invoiceService,
	}
}

type CreateInvoiceRequest struct {
	Items []InvoiceItemRequest `json:"items"`
}

// InvoiceItemRequest representa um item no payload
type InvoiceItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// PrintResponse representa a resposta da impressão
type PrintResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Invoice *domain.Invoice `json:"invoice,omitempty"`
}

// cria uma nova nota fiscal
func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload inválido", err.Error())
		return
	}

	// Converte os itens do request para domain
	items := make([]domain.InvoiceItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.InvoiceItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// Cria a nota fiscal
	invoice, err := h.invoiceService.CreateInvoice(items)
	if err != nil {
		switch err {
		case domain.ErrInvoiceNoItems:
			respondError(w, http.StatusBadRequest, "Nota fiscal deve ter ao menos um item", err.Error())
		case domain.ErrInvalidQuantity:
			respondError(w, http.StatusBadRequest, "Quantidade inválida", err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Erro ao criar nota fiscal", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, invoice)
}

// GetInvoice busca uma nota fiscal por ID
func (h *Handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	invoice, err := h.invoiceService.GetInvoice(id)
	if err != nil {
		if err == domain.ErrInvoiceNotFound {
			respondError(w, http.StatusNotFound, "Nota fiscal não encontrada", err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Erro ao buscar nota fiscal", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, invoice)
}

// GetAllInvoices lista todas as notas fiscais
func (h *Handler) GetAllInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.invoiceService.GetAllInvoices()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao listar notas fiscais", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, invoices)
}

// PrintInvoice "imprime" uma nota fiscal (fecha e atualiza estoque)
func (h *Handler) PrintInvoice(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// Processa a impressão
	invoice, err := h.invoiceService.PrintInvoice(id)
	if err != nil {
		switch err {
		case domain.ErrInvoiceNotFound:
			respondError(w, http.StatusNotFound, "Nota fiscal não encontrada", err.Error())
		case domain.ErrCannotPrintOpenInvoice:
			respondError(w, http.StatusBadRequest, "Nota fiscal já está fechada", err.Error())
		default:
			// Este é o cenário de falha do microsserviço
			// Retorna um erro detalhado para o frontend
			respondError(w, http.StatusServiceUnavailable, 
				"Falha na comunicação com o serviço de estoque", 
				"Não foi possível atualizar o estoque. Tente novamente.")
		}
		return
	}

	respondJSON(w, http.StatusOK, PrintResponse{
		Success: true,
		Message: "Nota fiscal impressa com sucesso",
		Invoice: invoice,
	})
}

// Health endpoint para healthcheck
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "billing",
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, error string, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   error,
		Message: message,
	})
}