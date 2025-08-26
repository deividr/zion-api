package postgres

import (
	"context"
	"encoding/json"
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

// // Main data query builder
// productsQuery := `COALESCE(
// 	(SELECT
// 		JSON_AGG(
// 			JSON_BUILD_OBJECT(
// 				'id', op.id,
// 				'orderId', op.order_id,
// 				'productId', op.product_id,
// 				'quantity', op.quantity,
// 				'unityType', op.unity_type,
// 				'price', op.price,
// 				'subProducts', (
// 					SELECT COALESCE(JSON_AGG(
// 						JSON_BUILD_OBJECT(
// 							'id', osp.id,
// 							'orderProductId', osp.order_product_id,
// 							'productId', osp.product_id
// 						)
// 					), '[]'::json)
// 					FROM order_sub_products osp
// 					WHERE osp.order_product_id = op.id
// 				)
// 			)
// 		)
// 	FROM order_products op
// 	WHERE op.order_id = o.id),
// 	'[]'::json
// ) AS products_json`

func (r *PgOrderRepository) FindAll(pagination domain.Pagination, filters domain.FindAllOrderFilters) ([]domain.Order, domain.Pagination, error) {
	offset := pagination.Limit * (pagination.Page - 1)

	baseBuilder := r.qb.
		Select().
		From("orders o").
		Where(squirrel.Eq{"o.is_deleted": false}).
		Where(squirrel.Expr("o.pickup_date BETWEEN ? AND ?", filters.PickupDateStart, filters.PickupDateEnd))

	if filters.Search != nil {
		baseBuilder = baseBuilder.
			Join("customers c ON c.id = o.customer_id").
			Where(squirrel.Or{
				squirrel.ILike{"c.name": fmt.Sprintf("%%%s%%", *filters.Search)},
				squirrel.ILike{"c.phone": fmt.Sprintf("%%%s%%", *filters.Search)},
			})
	}

	countBuilder := baseBuilder.Column("count(DISTINCT o.id)")

	totalCountQuery, totalCountArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error building total query: %w", err)
	}

	var totalCount int
	err = r.db.QueryRow(context.Background(), totalCountQuery, totalCountArgs...).Scan(&totalCount)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching total number of orders: %v", err)
	}

	customerQuery := `(
		SELECT
			JSON_BUILD_OBJECT(
				'id', c.id,
				'name', c.name,
				'phone', c.phone,
				'phone2', c.phone,
				'email', c.email,
				'createdAt', to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
			)
		FROM customers c
		WHERE c.id = o.customer_id
	) AS customer`

	queryBuilder := baseBuilder.
		Columns(
			"o.id",
			"o.order_number",
			"o.pickup_date",
			"o.created_at",
			"o.updated_at",
			"o.employee_id",
			"o.order_local",
			"o.observations",
			"o.is_picked_up",
		).
		Column(customerQuery).
		OrderBy("o.pickup_date DESC").
		Limit(uint64(pagination.Limit)).
		Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error building fetch orders query: %w", err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching orders: %v", err)
	}
	defer rows.Close()

	// Process rows
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		var customerJson []byte

		if err := rows.Scan(
			&order.Id,
			&order.Number,
			&order.PickupDate,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Employee,
			&order.OrderLocal,
			&order.Observations,
			&order.IsPickedUp,
			&customerJson,
		); err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("error scanning order data: %w", err)
		}

		if err := json.Unmarshal(customerJson, &order.Customer); err != nil {
			return nil, domain.Pagination{}, fmt.Errorf("error unmarshaling order customer: %w", err)
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error iterating order rows: %w", err)
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
		&order.Customer.Id,
		&order.Employee,
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
		Values(500, &newOrder.PickupDate, &newOrder.Customer.Id, &newOrder.Employee, &newOrder.OrderLocal, &newOrder.Observations, &newOrder.IsPickedUp).
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

	return &domain.Order{
		NewOrder: newOrder,
		Id:       id,
		Number:   "500", // TODO: get number from sequence
	}, nil
}
