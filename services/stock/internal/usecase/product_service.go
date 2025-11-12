package usecase

import (
	"time"
	"github.com/google/uuid"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/domain"
)

// ProductService contém a lógica de negócio de produtos
// (Single Responsibility Principle - SRP: apenas lógica de negócio)
type ProductService struct {
	repo domain.ProductRepository
}

// NewProductService cria uma nova instância do serviço
// (Dependency Injection via construtor)
func NewProductService(repo domain.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

// CreateProduct cria um novo produto
func (s *ProductService) CreateProduct(code, description string, balance int) (*domain.Product, error) {
	product := &domain.Product{
		ID:          uuid.New().String(),
		Code:        code,
		Description: description,
		Balance:     balance,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Valida o produto
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// Persiste o produto
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct busca um produto por ID
func (s *ProductService) GetProduct(id string) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

// GetAllProducts retorna todos os produtos
func (s *ProductService) GetAllProducts() ([]*domain.Product, error) {
	return s.repo.FindAll()
}

// UpdateProduct atualiza um produto
func (s *ProductService) UpdateProduct(id, code, description string, balance int) (*domain.Product, error) {
	// Busca o produto existente
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Atualiza os campos
	product.Code = code
	product.Description = description
	product.Balance = balance
	product.UpdatedAt = time.Now()

	// Valida
	if err := product.Validate(); err != nil {
		return nil, err
	}

	// Persiste
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// ReserveStock reserva uma quantidade de estoque de um produto
// Esta função será chamada pelo serviço de Billing
func (s *ProductService) ReserveStock(productID string, quantity int) error {
	// Busca o produto
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}

	// Tenta reservar
	if err := product.Reserve(quantity); err != nil {
		return err
	}

	// Atualiza no repositório
	return s.repo.Update(product)
}

// ReserveMultipleProducts reserva múltiplos produtos (transação)
// Implementação simplificada - em produção, considerar padrão Saga
func (s *ProductService) ReserveMultipleProducts(requests []domain.ReservationRequest) ([]domain.ReservationResponse, error) {
	responses := make([]domain.ReservationResponse, 0, len(requests))
	
	// Em uma implementação real, isso seria uma transação
	// Por hora, vamos processar sequencialmente
	for _, req := range requests {
		product, err := s.repo.FindByID(req.ProductID)
		if err != nil {
			responses = append(responses, domain.ReservationResponse{
				Success:      false,
				ProductID:    req.ProductID,
				ErrorMessage: err.Error(),
			})
			continue
		}

		err = product.Reserve(req.Quantity)
		if err != nil {
			responses = append(responses, domain.ReservationResponse{
				Success:      false,
				ProductID:    req.ProductID,
				NewBalance:   product.Balance,
				ErrorMessage: err.Error(),
			})
			continue
		}

		// Atualiza no repositório
		if err := s.repo.Update(product); err != nil {
			responses = append(responses, domain.ReservationResponse{
				Success:      false,
				ProductID:    req.ProductID,
				ErrorMessage: err.Error(),
			})
			continue
		}

		responses = append(responses, domain.ReservationResponse{
			Success:    true,
			ProductID:  req.ProductID,
			NewBalance: product.Balance,
		})
	}

	return responses, nil
}

// CheckAvailability verifica se há estoque disponível (sem reservar)
func (s *ProductService) CheckAvailability(productID string, quantity int) (bool, error) {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return false, err
	}
	return product.CanReserve(quantity), nil
}