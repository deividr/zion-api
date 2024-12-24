package usecase

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type ProductUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) *ProductUseCase {
	return &ProductUseCase{repo: repo}
}

func (uc *ProductUseCase) GetAll() ([]domain.Product, error) {
	products, err := uc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %v", err)
	}
	return products, nil
}

func (uc *ProductUseCase) GetById(id string) (*domain.Product, error) {
	product, err := uc.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produto: %v", err)
	}
	return product, nil
}
