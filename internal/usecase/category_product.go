package usecase

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type CategoryProductUseCase struct {
	repo domain.CategoryProductRepository
}

func NewCategoryProductUseCase(repo domain.CategoryProductRepository) *CategoryProductUseCase {
	return &CategoryProductUseCase{repo: repo}
}

func (uc *CategoryProductUseCase) GetAll() ([]domain.CategoryProduct, error) {
	categories, err := uc.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %v", err)
	}
	return categories, nil
}

func (uc *CategoryProductUseCase) GetById(id string) (*domain.CategoryProduct, error) {
	category, err := uc.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categoria: %v", err)
	}
	return category, nil
}

func (uc *CategoryProductUseCase) Update(category domain.CategoryProduct) error {
	err := uc.repo.Update(category)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %v", err)
	}
	return nil
}

func (uc *CategoryProductUseCase) Delete(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %v", err)
	}
	return nil
}

func (uc *CategoryProductUseCase) Create(category domain.CategoryProduct) (*domain.CategoryProduct, error) {
	createdCategory, err := uc.repo.Create(category)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %v", err)
	}
	return createdCategory, nil
}
