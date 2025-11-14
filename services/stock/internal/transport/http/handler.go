package http

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/domain"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/usecase"
)

// Handler gerencia as requisições HTTP do serviço de estoque
type Handler struct {
	productService *usecase.ProductService
}

// NewHandler cria um novo handler
func NewHandler(productService *usecase.ProductService) *Handler {
	return &Handler{
		productService: productService,
	}
}

// CreateProductRequest representa o payload de criação de produto
type CreateProductRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}

// UpdateProductRequest representa o payload de atualização
type UpdateProductRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateProduct cria um novo produto
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload inválido", err.Error())
		return
	}

	product, err := h.productService.CreateProduct(req.Code, req.Description, req.Balance)
	if err != nil {
		switch err {
		case domain.ErrInvalidProduct:
			respondError(w, http.StatusBadRequest, "Dados do produto inválidos", err.Error())
		case domain.ErrDuplicateCode:
			respondError(w, http.StatusConflict, "Código de produto já existe", err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Erro ao criar produto", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusCreated, product)
}

// GetProduct busca um produto por ID
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	product, err := h.productService.GetProduct(id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			respondError(w, http.StatusNotFound, "Produto não encontrado", err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Erro ao buscar produto", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// GetAllProducts lista todos os produtos
func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao listar produtos", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, products)
}

// UpdateProduct atualiza um produto
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Payload inválido", err.Error())
		return
	}

	product, err := h.productService.UpdateProduct(id, req.Code, req.Description, req.Balance)
	if err != nil {
		switch err {
		case domain.ErrProductNotFound:
			respondError(w, http.StatusNotFound, "Produto não encontrado", err.Error())
		case domain.ErrInvalidProduct:
			respondError(w, http.StatusBadRequest, "Dados do produto inválidos", err.Error())
		case domain.ErrDuplicateCode:
			respondError(w, http.StatusConflict, "Código de produto já existe", err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "Erro ao atualizar produto", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, product)
}

// DeleteProduct deleta um produto
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	err := h.productService.DeleteProduct(id)
	if err != nil {
		if err == domain.ErrProductNotFound {
			respondError(w, http.StatusNotFound, "Produto não encontrado", err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "Erro ao deletar produto", err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Produto deletado com sucesso"})
}

// ReserveStock reserva estoque (chamado pelo serviço de Billing)
func (h *Handler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	var requests []domain.ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
		respondError(w, http.StatusBadRequest, "Payload inválido", err.Error())
		return
	}

	responses, err := h.productService.ReserveMultipleProducts(requests)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao reservar estoque", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, responses)
}

// Health endpoint para healthcheck
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"service": "stock",
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