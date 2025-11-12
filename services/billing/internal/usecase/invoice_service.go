package usecase

import (
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/billing/internal/domain"
)

// InvoiceService contém a lógica de negócio de notas fiscais
type InvoiceService struct {
	repo        domain.InvoiceRepository
	stockClient domain.StockClient
}

// NewInvoiceService cria uma nova instância do serviço
func NewInvoiceService(repo domain.InvoiceRepository, stockClient domain.StockClient) *InvoiceService {
	return &InvoiceService{
		repo:        repo,
		stockClient: stockClient,
	}
}

// CreateInvoice cria uma nova nota fiscal
func (s *InvoiceService) CreateInvoice(items []domain.InvoiceItem) (*domain.Invoice, error) {
	// Valida se há itens
	if len(items) == 0 {
		return nil, domain.ErrInvoiceNoItems
	}

	// Valida e enriquece os itens com informações do produto
	enrichedItems := make([]domain.InvoiceItem, 0, len(items))
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, domain.ErrInvalidQuantity
		}

		// Busca informações do produto no Stock Service
		product, err := s.stockClient.GetProduct(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar produto %s: %w", item.ProductID, err)
		}

		// Verifica disponibilidade
		if product.Balance < item.Quantity {
			return nil, fmt.Errorf("produto %s com estoque insuficiente (disponível: %d, solicitado: %d)", 
				product.Code, product.Balance, item.Quantity)
		}

		// Enriquece o item com informações do produto
		enrichedItems = append(enrichedItems, domain.InvoiceItem{
			ProductID:   item.ProductID,
			ProductCode: product.Code,
			Description: product.Description,
			Quantity:    item.Quantity,
		})
	}

	// Gera o próximo número sequencial
	number, err := s.repo.GetNextNumber()
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar número da nota: %w", err)
	}

	// Cria a nota fiscal
	invoice := &domain.Invoice{
		ID:        uuid.New().String(),
		Number:    number,
		Status:    domain.StatusOpen,
		Items:     enrichedItems,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Valida a nota
	if err := invoice.Validate(); err != nil {
		return nil, err
	}

	// Persiste a nota
	if err := s.repo.Create(invoice); err != nil {
		return nil, fmt.Errorf("erro ao criar nota fiscal: %w", err)
	}

	return invoice, nil
}

// GetInvoice busca uma nota fiscal por ID
func (s *InvoiceService) GetInvoice(id string) (*domain.Invoice, error) {
	return s.repo.FindByID(id)
}

// GetAllInvoices retorna todas as notas fiscais
func (s *InvoiceService) GetAllInvoices() ([]*domain.Invoice, error) {
	return s.repo.FindAll()
}

// PrintInvoice "imprime" a nota fiscal (fecha e atualiza estoque)
// Esta é a operação mais crítica do sistema
func (s *InvoiceService) PrintInvoice(id string) (*domain.Invoice, error) {
	// Busca a nota fiscal
	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Verifica se a nota pode ser impressa
	if !invoice.CanBePrinted() {
		return nil, domain.ErrCannotPrintOpenInvoice
	}

	// Reserva os produtos no Stock Service
	// IMPORTANTE: Esta operação deve ser idempotente em produção
	err = s.stockClient.ReserveProducts(invoice.Items)
	if err != nil {
		return nil, fmt.Errorf("erro ao reservar produtos: %w", err)
	}

	// Fecha a nota fiscal
	if err := invoice.Close(); err != nil {
		return nil, err
	}

	// Atualiza no repositório
	if err := s.repo.Update(invoice); err != nil {
		return nil, fmt.Errorf("erro ao atualizar nota fiscal: %w", err)
	}

	return invoice, nil
}

// ValidateInvoiceItems valida se os itens podem ser adicionados à nota
func (s *InvoiceService) ValidateInvoiceItems(items []domain.InvoiceItem) error {
	for _, item := range items {
		if item.Quantity <= 0 {
			return domain.ErrInvalidQuantity
		}

		// Verifica se o produto existe
		_, err := s.stockClient.GetProduct(item.ProductID)
		if err != nil {
			return fmt.Errorf("produto %s não encontrado", item.ProductID)
		}

		// Verifica disponibilidade
		available, err := s.stockClient.CheckAvailability(item.ProductID, item.Quantity)
		if err != nil {
			return fmt.Errorf("erro ao verificar disponibilidade: %w", err)
		}
		if !available {
			return fmt.Errorf("produto %s sem estoque suficiente", item.ProductID)
		}
	}
	return nil
}