package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/deividr/zion-api/internal/domain"
	"github.com/jackc/pgx/v5"
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
	var customerJSON, productsJSON string
	var addressJSON *string

	err := r.db.QueryRow(context.Background(), `
		SELECT o.id,
			   o.order_number,
			   o.pickup_date,
			   o.created_at,
			   o.updated_at,
			   o.employee_id,
			   o.order_local,
			   o.observations,
			   o.is_picked_up,
			   CASE
				   WHEN a.id IS NULL THEN NULL
				   ELSE JSON_BUILD_OBJECT(
					   'id', a.id,
					   'cep', a.cep,
					   'street', a.street,
					   'number', a.number,
					   'neighborhood', a.neighborhood,
					   'city', a.city,
					   'state', a.state,
					   'aditionalDetails', a.aditional_details,
					   'distance', a.distance
				   )
			   END AS address,
			   JSON_BUILD_OBJECT(
				   'id', c.id,
				   'name', c.name,
				   'phone', c.phone,
				   'phone2', c.phone2,
				   'email', c.email
			   ) AS customer,
			   COALESCE((
				   SELECT JSON_AGG(
					   JSON_BUILD_OBJECT(
						   'id', op.id,
						   'orderId', op.order_id,
						   'productId', op.product_id,
						   'quantity', op.quantity,
						   'unityType', op.unity_type,
						   'price', op.price,
						   'name', p.name,
						   'subProducts', COALESCE((
							   SELECT JSON_AGG(
								   JSON_BUILD_OBJECT(
									   'id', osp.id,
									   'orderProductId', osp.order_product_id,
									   'productId', osp.product_id,
									   'name', p.name
								   )
							   )
							   FROM order_sub_products osp
							   JOIN products p ON p.id = osp.product_id
							   WHERE osp.order_product_id = op.id
						   ), '[]'::json)
					   )
				   )
				   FROM order_products op
				   JOIN products p ON p.id = op.product_id
				   WHERE op.order_id = o.id
			   ), '[]'::json) AS products
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		LEFT JOIN addresses a ON a.id = o.address_id
		WHERE o.id = $1 AND o.is_deleted = false
	`, id).Scan(
		&order.Id,
		&order.Number,
		&order.PickupDate,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Employee,
		&order.OrderLocal,
		&order.Observations,
		&order.IsPickedUp,
		&addressJSON,
		&customerJSON,
		&productsJSON,
	)
	if err != nil {
		fmt.Println("Erro no scan do resultado", err)
		return nil, fmt.Errorf("order not found: %v", err)
	}

	if addressJSON != nil {
		if err := json.Unmarshal([]byte(*addressJSON), &order.Address); err != nil {
			return nil, fmt.Errorf("error parsing address JSON: %v", err)
		}
	}

	if err := json.Unmarshal([]byte(customerJSON), &order.Customer); err != nil {
		return nil, fmt.Errorf("error parsing customer JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(productsJSON), &order.Products); err != nil {
		return nil, fmt.Errorf("error parsing products JSON: %v", err)
	}

	return &order, nil
}

func (r *PgOrderRepository) insertOrderProducts(tx pgx.Tx, orderID string, products []domain.OrderProduct) error {
	if len(products) == 0 {
		return nil
	}

	// Insert new products and get their IDs
	productsInsertBuilder := r.qb.Insert("order_products").
		Columns("order_id", "product_id", "quantity", "unity_type", "price")

	for _, p := range products {
		productsInsertBuilder = productsInsertBuilder.Values(orderID, p.ProductId, p.Quantity, p.UnityType, p.Price)
	}

	sql, args, err := productsInsertBuilder.Suffix("RETURNING id").ToSql()
	if err != nil {
		return fmt.Errorf("error building insert products query: %w", err)
	}

	rows, err := tx.Query(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("error inserting products: %w", err)
	}
	defer rows.Close()

	var insertedProductIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error scanning inserted product id: %w", err)
		}
		insertedProductIDs = append(insertedProductIDs, id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error after iterating inserted product ids: %w", err)
	}

	if len(insertedProductIDs) != len(products) {
		return fmt.Errorf("mismatch count of inserted products, expected %d but got %d", len(products), len(insertedProductIDs))
	}

	// Bulk insert sub-products
	subProductRows := [][]any{}
	for i, p := range products {
		if len(p.SubProducts) > 0 {
			orderProductID := insertedProductIDs[i]
			for _, sp := range p.SubProducts {
				subProductRows = append(subProductRows, []any{orderProductID, sp.ProductId})
			}
		}
	}

	if len(subProductRows) > 0 {
		_, err = tx.CopyFrom(
			context.Background(),
			pgx.Identifier{"order_sub_products"},
			[]string{"order_product_id", "product_id"},
			pgx.CopyFromRows(subProductRows),
		)
		if err != nil {
			return fmt.Errorf("error bulk inserting sub-products: %w", err)
		}
	}

	return nil
}

func (r *PgOrderRepository) Update(order domain.Order) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	// Update order details
	updateBuilder, args, err := r.qb.
		Update("orders").
		Set("pickup_date", order.PickupDate).
		Set("order_local", order.OrderLocal).
		Set("observations", order.Observations).
		Set("is_picked_up", order.IsPickedUp).
		Where(squirrel.Eq{"id": order.Id}).
		Where(squirrel.Eq{"is_deleted": false}).ToSql()
	if err != nil {
		return fmt.Errorf("error building query to update order: %w", err)
	}

	if _, err := tx.Exec(context.Background(), updateBuilder, args...); err != nil {
		return fmt.Errorf("error updating order: %w", err)
	}

	// Delete old products
	if _, err := tx.Exec(context.Background(), "DELETE FROM order_products WHERE order_id = $1", order.Id); err != nil {
		return fmt.Errorf("error deleting old order products: %w", err)
	}

	if err := r.insertOrderProducts(tx, order.Id, order.Products); err != nil {
		return err
	}

	return tx.Commit(context.Background())
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
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	insertBuilder, args, errQB := r.qb.Insert("orders").
		Columns("number", "pickup_date", "customer_id", "employee_id", "order_local", "observations", "is_picked_up").
		Values(500, &newOrder.PickupDate, &newOrder.Customer.Id, &newOrder.Employee, &newOrder.OrderLocal, &newOrder.Observations, &newOrder.IsPickedUp).
		Suffix("RETURNING id").
		ToSql()

	if errQB != nil {
		return nil, fmt.Errorf("error building query to create order: %w", errQB)
	}

	var orderID string
	if err := tx.QueryRow(context.Background(), insertBuilder, args...).Scan(&orderID); err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	// Assuming newOrder has a 'Products' field.
	if err := r.insertOrderProducts(tx, orderID, newOrder.Products); err != nil {
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return &domain.Order{
		NewOrder: newOrder,
		Id:       orderID,
		Number:   "500", // TODO: get number from sequence
	}, nil
}
