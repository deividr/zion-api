package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgProductRepository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func NewPgProductRepository(db *pgxpool.Pool) *PgProductRepository {
	return &PgProductRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PgProductRepository) FindAll(pagination domain.Pagination, filters domain.FindAllProductFilters) ([]domain.Product, domain.Pagination, error) {
	offset := pagination.Limit * (pagination.Page - 1)

	baseQuery := r.qb.
		Where(squirrel.Eq{"is_deleted": false}).
		Where(squirrel.ILike{"name": "%" + filters.Name + "%"})

	totalCountQuery, totalCountArgs, err := baseQuery.Select("count(*)").From("products").ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao construir query de total: %v", err)
	}

	var totalCount int

	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao buscar total de produtos: %v", err)
	}

	query, args, err := baseQuery.
		Select("id", "name", "value", "unity_type", "category_id").
		From("products").
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao construir query: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao buscar produtos: %v", err)
	}
	defer rows.Close()

	var products []domain.Product

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Value,
			&product.UnityType,
			&product.CategoryId,
		)
		if err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("erro ao ler produto: %v", err)
		}
		products = append(products, product)
	}

	// Update pagination with total count
	pagination.Total = totalCount

	return products, pagination, nil
}

func (r *PgProductRepository) FindById(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, value, unity_type, category_id
		FROM products 
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&product.Id,
		&product.Name,
		&product.Value,
		&product.UnityType,
		&product.CategoryId,
	)
	if err != nil {
		return nil, fmt.Errorf("produto não encontrado: %v", err)
	}

	return &product, nil
}

func (r *PgProductRepository) Update(product domain.Product) error {
	updateBuilder, args, err := r.qb.
		Update("products").Set("name", product.Name).
		Set("value", product.Value).
		Set("unity_type", product.UnityType).
		Set("category_id", product.CategoryId).
		Where(squirrel.Eq{"id": product.Id}).
		Where(squirrel.Eq{"is_deleted": false}).ToSql()
	if err != nil {
		return fmt.Errorf("erro ao construir query para atualizar o produto: %v", err)
	}

	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %v", err)
	}
	defer result.Close()

	return nil
}

func (r *PgProductRepository) Delete(id string) error {
	result, err := r.db.Query(context.Background(), "UPDATE products SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgProductRepository) Create(newProduct domain.NewProduct) (*domain.Product, error) {
	insertBuilder, args, errQB := r.qb.Insert("products").
		Columns("name", "value", "unity_type", "category_id").
		Values(&newProduct.Name, &newProduct.Value, &newProduct.UnityType, &newProduct.CategoryId).
		Suffix("RETURNING id").
		ToSql()

	if errQB != nil {
		return nil, fmt.Errorf("erro ao construir query para criar o produto: %v", errQB)
	}

	var id string
	errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&id)

	if errQuery != nil {
		return nil, fmt.Errorf("erro ao criar produto: %v", errQB)
	}

	createdProduct := &domain.Product{
		Id:         id,
		Name:       newProduct.Name,
		Value:      newProduct.Value,
		UnityType:  newProduct.UnityType,
		CategoryId: newProduct.CategoryId,
	}

	return createdProduct, nil
}
