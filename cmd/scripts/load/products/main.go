package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

type Product struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Value      uint32 `json:"value"`
	UnityType  string `json:"unityType"`
	OldId      int    `json:"oldId"`
	CategoryId string `json:"categoryId"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar .env: ", err)
		return
	}

	fmt.Println(os.Getenv("DATABASE_URL_LOAD"))

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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	dbOldPool, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Erro na conexão com o banco MySql", err)
		return
	}
	defer dbOldPool.Close()

	results, err2 := dbOldPool.Query("select cd_produto, nm_produto, st_unidade, vl_produto from produto;")
	if err2 != nil {
		fmt.Println("Erro na query", err2)
		return
	}

	var products []Product

	for results.Next() {
		var product Product
		var floatValue float64

		err := results.Scan(&product.OldId, &product.Name, &product.UnityType, &floatValue)
		if err != nil {
			fmt.Println("Erro no scan", err)
			continue
		}

		product.Value = uint32(floatValue * 100)
		products = append(products, product)
	}

	bar := progressbar.NewOptions(len(products),
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

	for _, product := range products {
		wg.Add(1)
		go func(product Product) {
			defer wg.Done()
			_, err = dbNewPool.Exec(context.Background(),
				`INSERT INTO products (old_id, name, value, unity_type) VALUES ($1, $2, $3, $4)`,
				product.OldId, product.Name, product.Value, product.UnityType)
			if err != nil {
				atomic.AddInt64(&errorCount, 1)
				fmt.Println("Erro ao inserir no PostgreSQL: ", err)
				return
			}

			atomic.AddInt64(&successCount, 1)
		}(product)
	}

	dbNewPool.Exec(context.Background(), `UPDATE products SET category_id = (SELECT id FROM category_products WHERE name = 'Bebidas') WHERE unity_type = 'LT';`)
	dbNewPool.Exec(context.Background(), `UPDATE products SET category_id = (SELECT id FROM category_products WHERE name = 'Saladas') WHERE name ILIKE '{%maionese%,%\salpicão%}';`)
	dbNewPool.Exec(context.Background(), `UPDATE products SET category_id = (SELECT id FROM category_products WHERE name = 'Massas') WHERE unity_type = 'KG' ;`)
	dbNewPool.Exec(context.Background(), `UPDATE products SET category_id = (SELECT id FROM category_products WHERE name = 'Bebidas') WHERE name ~* '(coca|fanta|guarana|tubaina|suco|água|vinho)';`)
	dbNewPool.Exec(context.Background(), `UPDATE products SET category_id = (SELECT id FROM category_products WHERE name = 'Diversos') WHERE category_id IS NULL;`)
	dbNewPool.Exec(context.Background(), `ALTER TABLE products ALTER COLUMN category_id SET NOT NULL;`)

	bar.Finish()

	fmt.Printf("\n\n=== ESTATÍSTICAS FINAIS ===\n")
	fmt.Printf("Total processados: %d\n", len(products))
	fmt.Printf("Sucessos: %d\n", atomic.LoadInt64(&successCount))
	fmt.Printf("Erros: %d\n", atomic.LoadInt64(&errorCount))
	fmt.Printf("Taxa de sucesso: %.2f%%\n", float64(atomic.LoadInt64(&successCount))/float64(len(products))*100)
}
