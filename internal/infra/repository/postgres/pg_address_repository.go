package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgAddressRepository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func NewPgAddressRepository(db *pgxpool.Pool) *PgAddressRepository {
	return &PgAddressRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PgAddressRepository) FindAll(pagination domain.Pagination) ([]domain.Address, domain.Pagination, error) {
	offset := pagination.Limit * (pagination.Page - 1)

	baseQuery := r.qb.Where(squirrel.Eq{"is_deleted": false})

	totalCountQuery, totalCountArgs, err := baseQuery.Select("count(*)").From("addresses").ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro on total query build: %v", err)
	}

	var totalCount int

	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error on search total addresses: %v", err)
	}

	query, args, err := baseQuery.
		Select(
			"id",
			"customer_id",
			"cep",
			"street",
			"number",
			"neighborhood",
			"city",
			"state",
			"aditional_details",
			"distance",
		).
		From("addresses").
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset)).
		ToSql()

	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error on query build: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error on search addresses: %v", err)
	}
	defer rows.Close()

	var addresses []domain.Address

	for rows.Next() {
		var address domain.Address
		err := rows.Scan(
			&address.Id,
			&address.CustomerId,
			&address.Cep,
			&address.Street,
			&address.Number,
			&address.Neighborhood,
			&address.City,
			&address.State,
			&address.AditionalDetails,
			&address.Distance,
		)
		if err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("error on read address informations: %v", err)
		}
		addresses = append(addresses, address)
	}

	// Update pagination with total count
	pagination.Total = totalCount

	return addresses, pagination, nil
}

func (r *PgAddressRepository) FindById(id string) (*domain.Address, error) {
	var address domain.Address
	err := r.db.QueryRow(context.Background(), `
		SELECT
			id,
			old_id,
			customer_id,
			cep,
			street,
			number,
			neighborhood,
			city,
			state,
			aditional_details,
			distance,
		FROM addresses 
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&address.Id,
		&address.OldId,
		&address.CustomerId,
		&address.Cep,
		&address.Street,
		&address.Number,
		&address.Neighborhood,
		&address.City,
		&address.State,
		&address.AditionalDetails,
		&address.Distance,
	)

	if err != nil {
		return nil, fmt.Errorf("address not found: %v", err)
	}

	return &address, nil
}

func (r *PgAddressRepository) FindBy(filters map[string]interface{}) ([]domain.Address, error) {
	baseQuery := r.qb.Select(
		"id",
		"old_id",
		"customer_id",
		"cep",
		"street",
		"number",
		"neighborhood",
		"city",
		"state",
		"aditional_details",
		"distance",
	).From("addresses").Where(squirrel.Eq{"is_deleted": false})

	for key, value := range filters {
		baseQuery = baseQuery.Where(squirrel.Eq{key: value})
	}

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building query: %w", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching addresses: %w", err)
	}
	defer rows.Close()

	var addresses []domain.Address

	for rows.Next() {
		var address domain.Address
		err := rows.Scan(
			&address.Id,
			&address.OldId,
			&address.CustomerId,
			&address.Cep,
			&address.Street,
			&address.Number,
			&address.Neighborhood,
			&address.City,
			&address.State,
			&address.AditionalDetails,
			&address.Distance,
		)
		if err != nil {
			return nil, fmt.Errorf("error reading address information: %w", err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *PgAddressRepository) Update(address domain.Address) error {
	updateBuilder, args, err := r.qb.
		Update("addresses").
		Set("cep", address.Cep).
		Set("street", address.Street).
		Set("Number", address.Number).
		Set("Neighborhood", address.Neighborhood).
		Set("City", address.City).
		Set("state", address.State).
		Set("aditional_details", address.AditionalDetails).
		Set("distance", address.Distance).
		Where(squirrel.Eq{"id": address.Id}).
		Where(squirrel.Eq{"is_deleted": false}).ToSql()

	if err != nil {
		return fmt.Errorf("erro ao construir query para atualizar o endereço: %v", err)
	}

	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("erro ao atualizar endereço: %v", err)
	}
	defer result.Close()

	return nil
}

func (r *PgAddressRepository) Delete(id string) error {
	result, err := r.db.Query(context.Background(), "UPDATE addresses SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("erro ao deletar endereço: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgAddressRepository) Create(newAddress domain.NewAddress) (*domain.Address, error) {
	insertBuilder, args, errQB := r.qb.Insert("addresses").
		Columns("customer_id", "cep", "street", "number", "neighborhood", "city", "state", "aditional_details", "distance").
		Values(
			&newAddress.CustomerId,
			&newAddress.Cep,
			&newAddress.Street,
			&newAddress.Number,
			&newAddress.Neighborhood,
			&newAddress.City,
			&newAddress.State,
			&newAddress.AditionalDetails,
			&newAddress.Distance).
		Suffix("RETURNING id").
		ToSql()

	if errQB != nil {
		return nil, fmt.Errorf("erro ao construir query para criar o endereço: %v", errQB)
	}

	var id string
	errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&id)

	if errQuery != nil {
		return nil, fmt.Errorf("erro ao criar endereço: %v", errQB)
	}

	createdAddress := &domain.Address{
		Id:               id,
		CustomerId:       newAddress.CustomerId,
		Cep:              newAddress.Cep,
		Street:           newAddress.Street,
		Number:           newAddress.Number,
		Neighborhood:     newAddress.Neighborhood,
		City:             newAddress.City,
		State:            newAddress.State,
		AditionalDetails: newAddress.AditionalDetails,
		Distance:         newAddress.Distance,
	}

	return createdAddress, nil
}
