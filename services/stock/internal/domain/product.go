package domain

import (
	"errors"
	"time"
)

// produto no sistema
type Product struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Balance     int       `json:"balance"` // Saldo em estoque
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Erros de domínio
var (
	ErrProductNotFound      = errors.New("produto não encontrado")
	ErrInsufficientBalance  = errors.New("saldo insuficiente")
	ErrInvalidProduct       = errors.New("produto inválido")
	ErrDuplicateCode        = errors.New("código de produto já existe")
	ErrInvalidQuantity      = errors.New("quantidade inválida")
)

// Validate -> valida os dados do produto
func (p *Product) Validate() error {
	if p.Code == "" {
		return ErrInvalidProduct
	}
	if p.Description == "" {
		return ErrInvalidProduct
	}
	if p.Balance < 0 {
		return ErrInvalidProduct
	}
	return nil
}

// CanReserve verifica se há saldo suficiente para reserva
func (p *Product) CanReserve(quantity int) bool {
	return p.Balance >= quantity && quantity > 0
}

// Reserve reduz o saldo do produto (usado ao fechar nota fiscal)
func (p *Product) Reserve(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if !p.CanReserve(quantity) {
		return ErrInsufficientBalance
	}
	p.Balance -= quantity
	p.UpdatedAt = time.Now()
	return nil
}

// ProductRepository define o contrato para persistência de produtos
// (Interface Segregation Principle - ISP)
type ProductRepository interface {
	Create(product *Product) error
	FindByID(id string) (*Product, error)
	FindByCode(code string) (*Product, error)
	FindAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id string) error
}

// ReservationRequest representa uma solicitação de reserva de estoque
type ReservationRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// ReservationResponse representa o resultado de uma reserva
type ReservationResponse struct {
	Success      bool   `json:"success"`
	ProductID    string `json:"product_id"`
	NewBalance   int    `json:"new_balance"`
	ErrorMessage string `json:"error_message,omitempty"`
}