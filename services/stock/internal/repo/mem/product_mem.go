package mem

import (
	"sync"
	"github.com/VitorMozer9/korp-teste-VitorMozer/services/stock/internal/domain"
)

// ProductMemRepository implementa ProductRepository em memória
// (Dependency Inversion Principle - DIP)
type ProductMemRepository struct {
	mu       sync.RWMutex
	products map[string]*domain.Product
	codes    map[string]string 
}

// NewProductMemRepository cria uma nova instância do repositório
func NewProductMemRepository() *ProductMemRepository {
	return &ProductMemRepository{
		products: make(map[string]*domain.Product),
		codes:    make(map[string]string),
	}
}

// Create adiciona um novo produto
func (r *ProductMemRepository) Create(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Verifica se o código já existe
	if _, exists := r.codes[product.Code]; exists {
		return domain.ErrDuplicateCode
	}

	r.products[product.ID] = product
	r.codes[product.Code] = product.ID
	return nil
}

// FindByID busca um produto por ID
func (r *ProductMemRepository) FindByID(id string) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, domain.ErrProductNotFound
	}
	return product, nil
}

// FindByCode busca um produto por código
func (r *ProductMemRepository) FindByCode(code string) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.codes[code]
	if !exists {
		return nil, domain.ErrProductNotFound
	}
	return r.products[id], nil
}

// FindAll retorna todos os produtos
func (r *ProductMemRepository) FindAll() ([]*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]*domain.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}
	return products, nil
}

// Update atualiza um produto existente
func (r *ProductMemRepository) Update(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.products[product.ID]
	if !exists {
		return domain.ErrProductNotFound
	}

	// Se o código mudou, atualiza o mapa de códigos
	if existing.Code != product.Code {
		// Verifica se o novo código já existe
		if _, codeExists := r.codes[product.Code]; codeExists {
			return domain.ErrDuplicateCode
		}
		delete(r.codes, existing.Code)
		r.codes[product.Code] = product.ID
	}

	r.products[product.ID] = product
	return nil
}

// Delete remove um produto
func (r *ProductMemRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.products[id]
	if !exists {
		return domain.ErrProductNotFound
	}

	delete(r.products, id)
	delete(r.codes, product.Code)
	return nil
}