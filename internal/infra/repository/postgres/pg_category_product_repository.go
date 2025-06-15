package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCategoryProductRepository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func NewPgCategoryProductRepository(db *pgxpool.Pool) *PgCategoryProductRepository {
	return &PgCategoryProductRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PgCategoryProductRepository) FindAll() ([]domain.CategoryProduct, error) {
	totalCountQuery, totalCountArgs, err := r.qb.Select("count(*)").From("category_products").ToSql()
	if err != nil {
		return nil, fmt.Errorf("erro ao construir query de total: %v", err)
	}

	var totalCount int
	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar total de categorias: %v", err)
	}

	query, args, err := r.qb.Select("id", "name", "description").From("category_products").ToSql()
	if err != nil {
		return nil, fmt.Errorf("erro ao construir query: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %v", err)
	}
	defer rows.Close()

	var categories []domain.CategoryProduct

	for rows.Next() {
		var category domain.CategoryProduct
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler categoria: %v", err)
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *PgCategoryProductRepository) FindById(id string) (*domain.CategoryProduct, error) {
	var category domain.CategoryProduct
	err := r.db.QueryRow(context.Background(), `SELECT id, name, description FROM category_products WHERE id = $1`, id).Scan(
		&category.Id,
		&category.Name,
		&category.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("categoria não encontrada: %v", err)
	}
	return &category, nil
}

func (r *PgCategoryProductRepository) Update(category domain.CategoryProduct) error {
	updateBuilder, args, err := r.qb.Update("category_products").Set("name", category.Name).Set("description", category.Description).Where(squirrel.Eq{"id": category.Id}).ToSql()
	if err != nil {
		return fmt.Errorf("erro ao construir query para atualizar a categoria: %v", err)
	}
	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar categoria: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgCategoryProductRepository) Delete(id string) error {
	result, err := r.db.Query(context.Background(), "DELETE FROM category_products WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar categoria: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgCategoryProductRepository) Create(category domain.CategoryProduct) (*domain.CategoryProduct, error) {
	insertBuilder, args, errQB := r.qb.Insert("category_products").Columns("name", "description").Values(&category.Name, &category.Description).Suffix("RETURNING id").ToSql()
	if errQB != nil {
		return nil, fmt.Errorf("erro ao construir query para criar a categoria: %v", errQB)
	}
	var id string
	errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&id)
	if errQuery != nil {
		return nil, fmt.Errorf("erro ao criar categoria: %v", errQuery)
	}
	createdCategory := &domain.CategoryProduct{
		Id:          id,
		Name:        category.Name,
		Description: category.Description,
	}
	return createdCategory, nil
}
