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

	baseQuery := r.qb

	totalCountQuery, totalCountArgs, err := baseQuery.Select("count(*)").From("addresses").ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error on total query build: %v", err)
	}

	var totalCount int

	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error on search total addresses: %v", err)
	}

	query, args, err := baseQuery.
		Select(
			"id",
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
			cep,
			street,
			number,
			neighborhood,
			city,
			state,
			aditional_details,
			distance
		FROM addresses
		WHERE id = $1
	`, id).Scan(
		&address.Id,
		&address.OldId,
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
		"cep",
		"street",
		"number",
		"neighborhood",
		"city",
		"state",
		"aditional_details",
		"distance",
	).From("addresses")

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

func (r *PgAddressRepository) FindByCustomerId(customerId string) ([]domain.Address, error) {
	query := `
		SELECT
			a.id,
			a.old_id,
			a.cep,
			a.street,
			a.number,
			a.neighborhood,
			a.city,
			a.state,
			a.aditional_details,
			a.distance,
			ac.is_default
		FROM addresses a
		INNER JOIN address_customers ac ON a.id = ac.address_id
		WHERE ac.customer_id = $1
		ORDER BY ac.is_default DESC
	`

	rows, err := r.db.Query(context.Background(), query, customerId)
	if err != nil {
		return nil, fmt.Errorf("error searching addresses by customer id: %w", err)
	}
	defer rows.Close()

	var addresses []domain.Address

	for rows.Next() {
		var address domain.Address
		err := rows.Scan(
			&address.Id,
			&address.OldId,
			&address.Cep,
			&address.Street,
			&address.Number,
			&address.Neighborhood,
			&address.City,
			&address.State,
			&address.AditionalDetails,
			&address.Distance,
			&address.IsDefault,
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
		Where(squirrel.Eq{"id": address.Id}).ToSql()

	if err != nil {
		return fmt.Errorf("error building query to update address: %v", err)
	}

	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("error updating address: %v", err)
	}
	defer result.Close()

	return nil
}

func (r *PgAddressRepository) UpdateDefaultAddress(customerId string, addressId string) error {
	// First, remove default flag from all customer addresses
	_, err := r.db.Exec(context.Background(),
		"UPDATE address_customers SET is_default = false WHERE customer_id = $1",
		customerId,
	)
	if err != nil {
		return fmt.Errorf("error removing previous default address: %v", err)
	}

	// Set the new default address
	_, err = r.db.Exec(context.Background(),
		"UPDATE address_customers SET is_default = true WHERE customer_id = $1 AND address_id = $2",
		customerId,
		addressId,
	)
	if err != nil {
		return fmt.Errorf("error setting new default address: %v", err)
	}

	return nil
}

func (r *PgAddressRepository) Delete(customerId string, addressId string) error {
	// Remove only the relationship in address_customers table, keeping the address in addresses table
	result, err := r.db.Exec(context.Background(),
		"DELETE FROM address_customers WHERE address_id = $1 AND customer_id = $2",
		addressId, customerId)
	if err != nil {
		return fmt.Errorf("error deleting address relationship: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("address relationship not found")
	}

	return nil
}

func (r *PgAddressRepository) Create(customerId string, newAddress domain.NewAddress) (*domain.Address, error) {
	var addressId string

	// Check if an address with the same CEP and number already exists
	checkQuery := `
		SELECT id
		FROM addresses
		WHERE cep = $1 AND number = $2
		LIMIT 1
	`

	err := r.db.QueryRow(context.Background(), checkQuery, newAddress.Cep, newAddress.Number).Scan(&addressId)

	// If no existing address found, create a new one
	if err != nil {
		insertBuilder, args, errQB := r.qb.Insert("addresses").
			Columns("cep", "street", "number", "neighborhood", "city", "state", "aditional_details", "distance").
			Values(
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
			return nil, fmt.Errorf("error building query to create address: %v", errQB)
		}

		errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&addressId)
		if errQuery != nil {
			return nil, fmt.Errorf("error creating address: %v", errQuery)
		}
	}

	// Check if the relationship already exists
	var existingRelationship bool
	checkRelationshipQuery := `
		SELECT EXISTS(
			SELECT 1 FROM address_customers
			WHERE address_id = $1 AND customer_id = $2
		)
	`
	err = r.db.QueryRow(context.Background(), checkRelationshipQuery, addressId, customerId).Scan(&existingRelationship)
	if err != nil {
		return nil, fmt.Errorf("error checking existing relationship: %v", err)
	}

	if existingRelationship {
		return nil, fmt.Errorf("this address is already associated with the customer")
	}

	// If the address is marked as default, remove default flag from other addresses
	if newAddress.IsDefault != nil && *newAddress.IsDefault {
		_, err := r.db.Exec(context.Background(),
			"UPDATE address_customers SET is_default = false WHERE customer_id = $1",
			customerId,
		)
		if err != nil {
			return nil, fmt.Errorf("error removing previous default address: %v", err)
		}
	}

	// Create the relationship between address and customer
	_, err = r.db.Exec(context.Background(),
		"INSERT INTO address_customers (address_id, customer_id, is_default) VALUES ($1, $2, $3)",
		addressId,
		customerId,
		*newAddress.IsDefault,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating relationship between address and customer: %v", err)
	}

	// Fetch the created/existing address to return
	createdAddress, err := r.FindById(addressId)
	if err != nil {
		return nil, fmt.Errorf("error fetching created address: %v", err)
	}

	createdAddress.IsDefault = newAddress.IsDefault

	return createdAddress, nil
}
