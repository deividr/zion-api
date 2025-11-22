package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func main() {
	loc, _ := time.LoadLocation("America/Sao_Paulo")

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
			&order.Customer.Id,
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
				t := pickupDate.Time
				saoPauloTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
				order.PickupDate = saoPauloTime.UTC()
			}

			if createdAt.Valid {
				order.CreatedAt = createdAt.Time.In(loc).UTC()
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

	if err := results.Close(); err != nil {
		fmt.Println("Erro após iterar os resultados:", err)
		return
	}

	if err := dbOldPool.Close(); err != nil {
		fmt.Println("Erro ao fechar conexão com MySQL:", err)
		return
	}

	// Buscar produtos que contêm "ENTREGA" no nome e obter seus old_id
	deliveryProductsQuery := `SELECT old_id FROM products WHERE UPPER(name) LIKE '%ENTREGA%'`
	deliveryRows, err := dbNewPool.Query(context.Background(), deliveryProductsQuery)
	if err != nil {
		fmt.Println("Erro ao buscar produtos com ENTREGA:", err)
		return
	}

	deliveryProductIds := make(map[int]bool)
	for deliveryRows.Next() {
		var oldId int
		if err := deliveryRows.Scan(&oldId); err != nil {
			fmt.Println("Erro ao escanear old_id de produto com ENTREGA:", err)
			continue
		}
		deliveryProductIds[oldId] = true
	}
	deliveryRows.Close()

	// Criar um map para armazenar quais orders precisam de endereço
	ordersNeedingAddress := make(map[string]bool)

	// Verificar quais orders têm produtos com entrega
	for _, product := range order_products {
		oldProductId, err := strconv.Atoi(product.ProductId)
		if err != nil {
			continue
		}
		if deliveryProductIds[oldProductId] {
			ordersNeedingAddress[product.OrderId] = true
		}
	}

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

	processOrder := func(order domain.Order, needsAddress bool) {
		err := dbNewPool.QueryRow(context.Background(),
			`SELECT id FROM customers WHERE old_id = $1`,
			order.Customer.Id,
		).Scan(&order.Customer.Id)
		if err != nil {
			fmt.Printf("Erro ao obter Customer para pedido %s: %v\n", order.Number, err)
			atomic.AddInt64(&errorCount, 1)
			return
		}

		var addressId *string

		// Se a order precisa de endereço, buscar o endereço principal do customer
		if needsAddress {
			var addr string
			err := dbNewPool.QueryRow(context.Background(),
				`SELECT a.id
				 FROM addresses a
				 INNER JOIN address_customers ac ON ac.address_id = a.id
				 WHERE ac.customer_id = $1 AND ac.is_default = true
				 LIMIT 1`,
				order.Customer.Id,
			).Scan(&addr)
			if err != nil {
				fmt.Printf("Aviso: Endereço principal não encontrado para customer %s no pedido %s: %v\n", order.Customer.Id, order.Number, err)
				// Não retornamos erro aqui, apenas logamos e continuamos sem endereço
			} else {
				addressId = &addr
			}
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO orders (id, order_number, pickup_date, created_at, customer_id, employee_id, order_local, observations, is_picked_up, address_id)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			order.Id,
			order.Number,
			order.PickupDate,
			order.CreatedAt,
			order.Customer.Id,
			order.Employee,
			order.OrderLocal,
			order.Observations,
			order.IsPickedUp,
			addressId,
		)
		if err != nil {
			fmt.Printf("Erro ao inserir Order %s: %v\n", order.Number, err)
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
			defer func() { _ = bar.Add(1) }()

			// 1. Processa a ordem principal
			needsAddress := ordersNeedingAddress[o.Id]
			processOrder(o, needsAddress)

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

	if err := bar.Finish(); err != nil {
		fmt.Println("Erro ao finalizar a barra de progresso:", err)
	}

	fmt.Printf(" === ESTATÍSTICAS FINAIS === ")
	fmt.Printf("Total de pedidos únicos processados: %d\n", len(orders))
	fmt.Printf("Sucessos: %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Erros: %d\n", atomic.LoadInt64(&errorCount))
	if len(orders) > 0 {
		successRate := float64(atomic.LoadInt64(&successCount)) / float64(len(orders)) * 100
		fmt.Printf("Taxa de sucesso: %.2f%%\n", successRate)
	}
}
