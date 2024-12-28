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

	// Construindo a query com Squirrel
	query, args, err := r.qb.
		Select("id", "name", "value", "unity_type").
		From("products").
		Where(squirrel.Eq{"is_deleted": false}).
		Where(squirrel.ILike{"name": "%" + filters.Name + "%"}).
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao construir query: %v", err)
	}

	fmt.Println(query)
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
		)
		if err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("erro ao ler produto: %v", err)
		}
		products = append(products, product)
	}

	return products, pagination, nil
}

func (r *PgProductRepository) FindById(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, value, unity_type 
		FROM products 
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&product.Id,
		&product.Name,
		&product.Value,
		&product.UnityType,
	)

	if err != nil {
		return nil, fmt.Errorf("produto não encontrado: %v", err)
	}

	return &product, nil
}
