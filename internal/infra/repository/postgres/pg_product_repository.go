package postgres

import (
	"context"
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgProductRepository struct {
	db *pgxpool.Pool
}

func NewPgProductRepository(db *pgxpool.Pool) *PgProductRepository {
	return &PgProductRepository{
		db: db,
	}
}

func (r *PgProductRepository) FindAll() ([]domain.Product, error) {
	rows, err := r.db.Query(context.Background(), `
		SELECT id, name, value, unity_type 
		FROM products 
		WHERE is_deleted = false
	`)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %v", err)
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
			return nil, fmt.Errorf("erro ao ler produto: %v", err)
		}
		products = append(products, product)
	}

	return products, nil
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
