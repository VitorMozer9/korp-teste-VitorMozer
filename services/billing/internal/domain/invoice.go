package domain

import (
	"errors"
	"time"
)

// InvoiceStatus representa o status de uma nota fiscal
type InvoiceStatus string

const (
	StatusOpen   InvoiceStatus = "ABERTA"
	StatusClosed InvoiceStatus = "FECHADA"
)

// Invoice representa uma nota fiscal
type Invoice struct {
	ID        string          `json:"id"`
	Number    int             `json:"number"`    // Numeração sequencial
	Status    InvoiceStatus   `json:"status"`    // ABERTA ou FECHADA
	Items     []InvoiceItem   `json:"items"`     // Produtos da nota
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	ClosedAt  *time.Time      `json:"closed_at,omitempty"` // Data de fechamento
}

// InvoiceItem representa um item (produto) na nota fiscal
type InvoiceItem struct {
	ProductID   string `json:"product_id"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// Erros de domínio
var (
	ErrInvoiceNotFound       = errors.New("nota fiscal não encontrada")
	ErrInvalidInvoice        = errors.New("nota fiscal inválida")
	ErrInvoiceAlreadyClosed  = errors.New("nota fiscal já está fechada")
	ErrInvoiceNoItems        = errors.New("nota fiscal deve ter ao menos um item")
	ErrInvalidQuantity       = errors.New("quantidade inválida")
	ErrCannotPrintOpenInvoice = errors.New("não é possível imprimir nota em status diferente de ABERTA")
)

// Validate valida os dados da nota fiscal
func (i *Invoice) Validate() error {
	if i.Number <= 0 {
		return ErrInvalidInvoice
	}
	if len(i.Items) == 0 {
		return ErrInvoiceNoItems
	}
	for _, item := range i.Items {
		if item.ProductID == "" || item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
	}
	return nil
}

// CanBePrinted verifica se a nota pode ser impressa (fechada)
func (i *Invoice) CanBePrinted() bool {
	return i.Status == StatusOpen
}

// Close fecha a nota fiscal (equivalente a "imprimir")
func (i *Invoice) Close() error {
	if !i.CanBePrinted() {
		return ErrCannotPrintOpenInvoice
	}
	
	now := time.Now()
	i.Status = StatusClosed
	i.ClosedAt = &now
	i.UpdatedAt = now
	
	return nil
}

// IsOpen verifica se a nota está aberta
func (i *Invoice) IsOpen() bool {
	return i.Status == StatusOpen
}

// IsClosed verifica se a nota está fechada
func (i *Invoice) IsClosed() bool {
	return i.Status == StatusClosed
}

// InvoiceRepository define o contrato para persistência de notas fiscais
type InvoiceRepository interface {
	Create(invoice *Invoice) error
	FindByID(id string) (*Invoice, error)
	FindAll() ([]*Invoice, error)
	Update(invoice *Invoice) error
	GetNextNumber() (int, error) // Retorna o próximo número sequencial
}

// StockClient define o contrato para comunicação com o Stock Service
// (Interface Segregation Principle - ISP)
type StockClient interface {
	ReserveProducts(items []InvoiceItem) error
	CheckAvailability(productID string, quantity int) (bool, error)
	GetProduct(productID string) (*ProductInfo, error)
}

// ProductInfo representa informações básicas de um produto
type ProductInfo struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}