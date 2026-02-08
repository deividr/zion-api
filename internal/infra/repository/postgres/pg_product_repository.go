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

func (r *PgProductRepository) FindAll(filters domain.FindAllProductFilters) ([]domain.Product, error) {
	query, args, err := r.qb.
		Select("id", "name", "value", "unity_type", "category_id", "image_url", "is_variable_price").
		From("products").
		Where(squirrel.Eq{"is_deleted": false}).
		Where(squirrel.ILike{"name": "%" + filters.Name + "%"}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("erro ao construir query: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
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
			&product.CategoryId,
			&product.ImageUrl,
			&product.IsVariablePrice,
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
		SELECT id, name, value, unity_type, category_id, image_url, is_variable_price
		FROM products
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&product.Id,
		&product.Name,
		&product.Value,
		&product.UnityType,
		&product.CategoryId,
		&product.ImageUrl,
		&product.IsVariablePrice,
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
		Set("image_url", product.ImageUrl).
		Set("is_variable_price", product.IsVariablePrice).
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
		Columns("name", "value", "unity_type", "category_id", "image_url", "is_variable_price").
		Values(&newProduct.Name, &newProduct.Value, &newProduct.UnityType, &newProduct.CategoryId, &newProduct.ImageUrl, &newProduct.IsVariablePrice).
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
		Id:              id,
		Name:            newProduct.Name,
		Value:           newProduct.Value,
		UnityType:       newProduct.UnityType,
		CategoryId:      newProduct.CategoryId,
		ImageUrl:        newProduct.ImageUrl,
		IsVariablePrice: newProduct.IsVariablePrice,
	}

	return createdProduct, nil
}
