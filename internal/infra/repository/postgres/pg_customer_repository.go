package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgCustomerRepository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func NewPgCustomerRepository(db *pgxpool.Pool) *PgCustomerRepository {
	return &PgCustomerRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PgCustomerRepository) FindAll(pagination domain.Pagination, filters domain.FindAllCustomerFilters) ([]domain.Customer, domain.Pagination, error) {
	offset := pagination.Limit * (pagination.Page - 1)

	baseQuery := r.qb.
		Where(squirrel.Eq{"is_deleted": false})

	var filterConditions []squirrel.Sqlizer

	if filters.Name != "" {
		filterConditions = append(filterConditions, squirrel.ILike{"name": "%" + filters.Name + "%"})
	}

	if filters.Phone != "" {
		filterConditions = append(filterConditions, squirrel.Or{
			squirrel.ILike{"phone": "%" + filters.Phone + "%"},
			squirrel.ILike{"phone2": "%" + filters.Phone + "%"},
		})
	}

	if filters.Email != "" {
		filterConditions = append(filterConditions, squirrel.ILike{"email": "%" + filters.Email + "%"})
	}

	if len(filterConditions) > 0 {
		if len(filterConditions) == 1 {
			baseQuery = baseQuery.Where(filterConditions[0])
		} else {
			orConditions := squirrel.Or{}
			for _, cond := range filterConditions {
				orConditions = append(orConditions, cond)
			}
			baseQuery = baseQuery.Where(orConditions)
		}
	}

	totalCountQuery, totalCountArgs, err := baseQuery.Select("count(*)").From("customers").ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao construir query de total: %v", err)
	}

	var totalCount int

	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao buscar total de clientes: %v", err)
	}

	query, args, err := baseQuery.
		Select("id", "name", "phone", "phone2", "email").
		From("customers").
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao construir query: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao buscar clientes: %v", err)
	}
	defer rows.Close()

	var customers []domain.Customer

	for rows.Next() {
		var customer domain.Customer
		err := rows.Scan(
			&customer.Id,
			&customer.Name,
			&customer.Phone,
			&customer.Phone2,
			&customer.Email,
		)
		if err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("erro ao ler cliente: %v", err)
		}
		customers = append(customers, customer)
	}

	pagination.Total = totalCount

	return customers, pagination, nil
}

func (r *PgCustomerRepository) FindById(id string) (*domain.Customer, error) {
	var customer domain.Customer
	err := r.db.QueryRow(context.Background(), `
		SELECT id, name, phone, phone2, email
		FROM customers
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&customer.Id,
		&customer.Name,
		&customer.Phone,
		&customer.Phone2,
		&customer.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("cliente não encontrado: %v", err)
	}

	return &customer, nil
}

func (r *PgCustomerRepository) Update(customer domain.Customer) error {
	updateBuilder, args, err := r.qb.
		Update("customers").Set("name", customer.Name).
		Set("phone", customer.Phone).
		Set("phone2", customer.Phone2).
		Set("email", customer.Email).
		Where(squirrel.Eq{"id": customer.Id}).
		Where(squirrel.Eq{"is_deleted": false}).
		ToSql()
	if err != nil {
		return fmt.Errorf("erro ao construir query para atualizar o cliente: %v", err)
	}

	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %v", err)
	}
	defer result.Close()

	return nil
}

func (r *PgCustomerRepository) Delete(id string) error {
	result, err := r.db.Query(context.Background(), "UPDATE customers SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar cliente: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgCustomerRepository) Create(newCustomer domain.NewCustomer) (*domain.Customer, error) {
	if newCustomer.Phone2 != nil && *newCustomer.Phone2 == "" {
		newCustomer.Phone2 = nil
	}

	insertBuilder, args, errQB := r.qb.Insert("customers").
		Columns("name", "phone", "phone2", "email").
		Values(&newCustomer.Name, &newCustomer.Phone, &newCustomer.Phone2, &newCustomer.Email).
		Suffix("RETURNING id").
		ToSql()

	if errQB != nil {
		return nil, fmt.Errorf("erro ao construir query para criar o cliente: %v", errQB)
	}

	var id string
	errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&id)

	if errQuery != nil {
		return nil, fmt.Errorf("erro ao criar cliente: %v", errQuery)
	}

	createdCustomer := &domain.Customer{
		Id:     id,
		Name:   newCustomer.Name,
		Phone:  newCustomer.Phone,
		Phone2: newCustomer.Phone2,
		Email:  newCustomer.Email,
	}

	return createdCustomer, nil
}
