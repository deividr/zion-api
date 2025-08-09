package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgOrderRepository struct {
	db *pgxpool.Pool
	qb squirrel.StatementBuilderType
}

func NewPgOrderRepository(db *pgxpool.Pool) *PgOrderRepository {
	return &PgOrderRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PgOrderRepository) FindAll(pagination domain.Pagination, filters domain.FindAllOrderFilters) ([]domain.Order, domain.Pagination, error) {
	offset := pagination.Limit * (pagination.Page - 1)

	baseQuery := r.qb.
		Where(squirrel.Eq{"is_deleted": false})

	totalCountQuery, totalCountArgs, err := baseQuery.
		Select("count(*)").
		From("orders o").
		ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error building total query: %w", err)
	}

	var totalCount int

	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching total number of orders: %v", err)
	}

	query, args, err := baseQuery.
		Select("id", "order_number", "pickup_date", "created_at", "updated_at", "customer_id", "employee_id", "order_local", "observations", "is_picked_up").
		From("orders").
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error building query: %v", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching orders: %v", err)
	}
	defer rows.Close()

	var orders []domain.Order

	for rows.Next() {
		var order domain.Order
		err := rows.Scan(
			&order.Id,
			&order.Number,
			&order.PickupDate,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.CustomerId,
			&order.EmployeeId,
			&order.OrderLocal,
			&order.Observations,
			&order.IsPickedUp,
		)
		if err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("error reading customer: %v", err)
		}

		orders = append(orders, order)
	}

	pagination.Total = totalCount

	return orders, pagination, nil
}

func (r *PgOrderRepository) FindById(id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.QueryRow(context.Background(), `
		SELECT id, number, pickup_date, created_at, updated_at, customer_id, employee_id, order_local, observations, is_picked_up
		FROM orders
		WHERE id = $1 AND is_deleted = false
	`, id).Scan(
		&order.Id,
		&order.Number,
		&order.PickupDate,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.CustomerId,
		&order.EmployeeId,
		&order.OrderLocal,
		&order.Observations,
		&order.IsPickedUp,
	)

	if err != nil {
		return nil, fmt.Errorf("order not found: %v", err)
	}

	return &order, nil
}

func (r *PgOrderRepository) Update(order domain.Order) error {
	updateBuilder, args, err := r.qb.
		Update("orders").Set("number", order.Number).
		Set("pickup_date", order.PickupDate).
		Set("order_local", order.OrderLocal).
		Set("observations", order.Observations).
		Set("is_picked_up", order.IsPickedUp).
		Where(squirrel.Eq{"id": order.Id}).
		Where(squirrel.Eq{"is_deleted": false}).ToSql()

	if err != nil {
		return fmt.Errorf("error building query to update customer: %v", err)
	}

	result, err := r.db.Query(context.Background(), updateBuilder, args...)
	if err != nil {
		return fmt.Errorf("error updating customer: %v", err)
	}
	defer result.Close()

	return nil
}

func (r *PgOrderRepository) Delete(id string) error {
	result, err := r.db.Query(context.Background(), "UPDATE orders SET is_deleted = true WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting order: %v", err)
	}
	defer result.Close()
	return nil
}

func (r *PgOrderRepository) Create(newOrder domain.NewOrder) (*domain.Order, error) {
	insertBuilder, args, errQB := r.qb.Insert("orders").
		Columns("number", "pickup_date", "customer_id", "employee_id", "order_local", "observations", "is_picked_up").
		Values(&newOrder.Number, &newOrder.PickupDate, &newOrder.CustomerId, &newOrder.EmployeeId, &newOrder.OrderLocal, &newOrder.Observations, &newOrder.IsPickedUp).
		Suffix("RETURNING id").
		ToSql()

	if errQB != nil {
		return nil, fmt.Errorf("error building query to create order: %v", errQB)
	}

	var id string
	errQuery := r.db.QueryRow(context.Background(), insertBuilder, args...).Scan(&id)

	if errQuery != nil {
		return nil, fmt.Errorf("error creating order: %v", errQB)
	}

	createdOrder := &domain.Order{
		Id:           id,
		Number:       newOrder.Number,
		PickupDate:   newOrder.PickupDate,
		CustomerId:   newOrder.CustomerId,
		EmployeeId:   newOrder.EmployeeId,
		OrderLocal:   newOrder.OrderLocal,
		Observations: newOrder.Observations,
		IsPickedUp:   newOrder.IsPickedUp,
	}

	return createdOrder, nil
}
