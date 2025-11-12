package mem

import (
	"sync"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/domain"
)

// InvoiceMemRepository implementa InvoiceRepository em memória
type InvoiceMemRepository struct {
	mu           sync.RWMutex
	invoices     map[string]*domain.Invoice
	lastNumber   int // Controla a numeração sequencial
}

// NewInvoiceMemRepository cria uma nova instância do repositório
func NewInvoiceMemRepository() *InvoiceMemRepository {
	return &InvoiceMemRepository{
		invoices:   make(map[string]*domain.Invoice),
		lastNumber: 0,
	}
}

// Create adiciona uma nova nota fiscal
func (r *InvoiceMemRepository) Create(invoice *domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.invoices[invoice.ID] = invoice
	return nil
}

// FindByID busca uma nota fiscal por ID
func (r *InvoiceMemRepository) FindByID(id string) (*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invoice, exists := r.invoices[id]
	if !exists {
		return nil, domain.ErrInvoiceNotFound
	}
	return invoice, nil
}

// FindAll retorna todas as notas fiscais
func (r *InvoiceMemRepository) FindAll() ([]*domain.Invoice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	invoices := make([]*domain.Invoice, 0, len(r.invoices))
	for _, invoice := range r.invoices {
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

// Update atualiza uma nota fiscal existente
func (r *InvoiceMemRepository) Update(invoice *domain.Invoice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.invoices[invoice.ID]; !exists {
		return domain.ErrInvoiceNotFound
	}

	r.invoices[invoice.ID] = invoice
	return nil
}

// GetNextNumber retorna o próximo número sequencial
func (r *InvoiceMemRepository) GetNextNumber() (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastNumber++
	return r.lastNumber, nil
}