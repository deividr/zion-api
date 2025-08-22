package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar .env: ", err)
		return
	}

	dbNewPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL_LOAD"))
	if err != nil {
		fmt.Println("Erro ao conectar no PostgreSQL:", err)
		return
	}
	defer dbNewPool.Close()

	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbname)
	dbOldPool, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Erro na conexão com o banco MySql", err)
		return
	}
	defer dbOldPool.Close()

	var totalCount int
	err = dbOldPool.QueryRow("SELECT COUNT(*) FROM pedido").Scan(&totalCount)
	if err != nil {
		fmt.Println("Erro ao contar registros:", err)
		return
	}

	fmt.Printf("Total de pedidos a serem processados: %d", totalCount)

	results, err := dbOldPool.Query(`
		SELECT p.cd_pedido,
		       p.nr_pedido,
		       p.dt_retirada,
		       p.dt_pedido,
		       p.cd_cliente,
		       p.nr_geladeira,
		       p.ds_observacao,
		       p.st_retirado,
		       ip.cd_produto AS "item_cd_produto",
		       ip.cd_molho AS "item_cd_molho",
           CAST(
		       CASE
					WHEN p2.st_unidade = 'UN' THEN
						CASE
							WHEN MOD(ip.vl_quantidade, 1) != 0 THEN CEIL(ip.vl_quantidade)
							ELSE ip.vl_quantidade
						END
					WHEN ip.vl_quantidade < 10 THEN
						ROUND(ip.vl_quantidade * 1000)
					ELSE
						ip.vl_quantidade
				END AS SIGNED
			) AS "item_vl_quantidade"
	  FROM pedido p
	  INNER JOIN item_pedido ip ON ip.cd_pedido = p.cd_pedido
	  INNER JOIN produto p2 ON p2.cd_produto = ip.cd_produto
	;`)
	if err != nil {
		fmt.Println("Erro na query", err)
		return
	}

	var orders []domain.Order
	var order_products []domain.OrderProduct
	var order_sub_products []domain.OrderSubProduct

	ordersMap := make(map[string]string)

	for results.Next() {
		var order domain.Order
		var order_product domain.OrderProduct
		var pickupDate, createdAt sql.NullTime
		var orderSubProductId *string

		err := results.Scan(
			&order.Id,
			&order.Number,
			&pickupDate,
			&createdAt,
			&order.CustomerId,
			&order.OrderLocal,
			&order.Observations,
			&order.IsPickedUp,
			&order_product.ProductId,
			&orderSubProductId,
			&order_product.Quantity,
		)
		if err != nil {
			fmt.Println("Erro no scan do resultado geral", err)
			continue
		}

		newOrderId, exists := ordersMap[order.Id]
		if !exists {
			newOrderId = uuid.NewString()
			ordersMap[order.Id] = newOrderId
			order.Id = newOrderId

			if pickupDate.Valid {
				order.PickupDate = pickupDate.Time
			}
			if createdAt.Valid {
				order.CreatedAt = createdAt.Time
			}
			orders = append(orders, order)
		}

		order_product.Id = uuid.NewString()
		order_product.OrderId = newOrderId

		if orderSubProductId != nil && *orderSubProductId != "" {
			var order_sub_product domain.OrderSubProduct
			order_sub_product.Id = uuid.NewString()
			order_sub_product.OrderProductId = order_product.Id
			order_sub_product.ProductId = *orderSubProductId
			order_sub_products = append(order_sub_products, order_sub_product)
		}

		order_products = append(order_products, order_product)
	}

	results.Close()
	dbOldPool.Close()

	// Agrupar produtos e subprodutos por seus respectivos IDs de pedido/produto
	orderProductsMap := make(map[string][]domain.OrderProduct)
	for _, p := range order_products {
		orderProductsMap[p.OrderId] = append(orderProductsMap[p.OrderId], p)
	}

	orderSubProductsMap := make(map[string][]domain.OrderSubProduct)
	for _, sp := range order_sub_products {
		orderSubProductsMap[sp.OrderProductId] = append(orderSubProductsMap[sp.OrderProductId], sp)
	}

	bar := progressbar.NewOptions(len(orders),
		progressbar.OptionSetDescription("Processando pedidos"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	var successCount int64
	var errorCount int64
	wg := sync.WaitGroup{}

	processOrder := func(order domain.Order) {
		err := dbNewPool.QueryRow(context.Background(),
			`SELECT id FROM customers WHERE old_id = $1`,
			order.CustomerId,
		).Scan(&order.CustomerId)
		if err != nil {
			fmt.Printf("Erro ao obter Customer para pedido %s: %v", order.Number, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO orders (id, order_number, pickup_date, created_at, customer_id, employee_id, order_local, observations, is_picked_up)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			order.Id,
			order.Number,
			order.PickupDate,
			order.CreatedAt,
			order.CustomerId,
			order.EmployeeId,
			order.OrderLocal,
			order.Observations,
			order.IsPickedUp,
		)
		if err != nil {
			fmt.Printf("Erro ao inserir Order %s: %v", order.Number, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		atomic.AddInt64(&successCount, 1)
	}

	processOrderProduct := func(order_product domain.OrderProduct) {
		oldProductId, err := strconv.Atoi(order_product.ProductId)
		if err != nil {
			fmt.Printf("Erro ao converter ProductId %s para int: %v", order_product.ProductId, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		err = dbNewPool.QueryRow(context.Background(),
			`SELECT id, unity_type, value FROM products WHERE old_id = $1`,
			oldProductId,
		).Scan(&order_product.ProductId, &order_product.UnityType, &order_product.Price)
		if err != nil {
			fmt.Printf("Erro ao selecionar produto %s para order_product: %v", order_product.ProductId, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO order_products (id, order_id, product_id, quantity, unity_type, price) VALUES ($1, $2, $3, $4, $5, $6)`,
			order_product.Id, order_product.OrderId, order_product.ProductId, order_product.Quantity, order_product.UnityType, order_product.Price,
		)
		if err != nil {
			fmt.Printf("Erro ao inserir order_product %s: %v", order_product.Id, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}
	}

	processSubProduct := func(order_sub_product domain.OrderSubProduct) {
		oldProductId, err := strconv.Atoi(order_sub_product.ProductId)
		if err != nil {
			fmt.Printf("Erro ao converter ProductId %s para int: %v", order_sub_product.ProductId, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		err = dbNewPool.QueryRow(context.Background(),
			`SELECT id FROM products WHERE old_id = $1`,
			oldProductId,
		).Scan(&order_sub_product.ProductId)
		if err != nil {
			fmt.Printf("Erro ao selecionar produto %s para order_sub_product: %v", order_sub_product.ProductId, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO order_sub_products (id, order_product_id, product_id) VALUES ($1, $2, $3)`,
			order_sub_product.Id, order_sub_product.OrderProductId, order_sub_product.ProductId)
		if err != nil {
			fmt.Printf("Erro ao inserir order_sub_product %s: %v", order_sub_product.Id, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}
	}

	// Loop único para processar cada pedido e seus itens de forma atômica
	for _, order := range orders {
		wg.Add(1)
		go func(o domain.Order) {
			defer wg.Done()
			defer bar.Add(1)

			// 1. Processa a ordem principal
			processOrder(o)

			// 2. Processa os produtos associados
			if orderProduct, ok := orderProductsMap[o.Id]; ok {
				for _, p := range orderProduct {
					processOrderProduct(p)

					// 3. Processa os sub-produtos associados
					if subProducts, ok := orderSubProductsMap[p.Id]; ok {
						for _, sp := range subProducts {
							processSubProduct(sp)
						}
					}
				}
			}
		}(order)
	}

	wg.Wait()
	bar.Finish()

	fmt.Printf(" === ESTATÍSTICAS FINAIS === ")
	fmt.Printf("Total de pedidos únicos processados: %d", len(orders))
	fmt.Printf("Sucessos: %d", atomic.LoadInt64(&successCount))
	fmt.Printf("Erros: %d", atomic.LoadInt64(&errorCount))
	if len(orders) > 0 {
		successRate := float64(atomic.LoadInt64(&successCount)) / float64(len(orders)) * 100
		fmt.Printf("Taxa de sucesso: %.2f%%", successRate)
	}
}
