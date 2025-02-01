package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Product struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Value     uint32 `json:"value"`
	UnityType string `json:"unityType"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar .env: ", err)
		return
	}

	fmt.Println(os.Getenv("DATABASE_URL"))

	dbNewPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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

	for results.Next() {
		var product Product
		var floatValue float64

		err3 := results.Scan(&product.Id, &product.Name, &product.UnityType, &floatValue)
		if err3 != nil {
			fmt.Println("Erro no scan", err3)
			continue
		}

		product.Value = uint32(floatValue * 100)

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO products (name, value, unity_type) VALUES ($1, $2, $3)`,
			product.Name, product.Value, product.UnityType)

		if err != nil {
			fmt.Println("Erro ao inserir no PostgreSQL: ", err)
			continue
		}

		fmt.Printf("Produto inserido com sucesso: %+v\n", product.Name)
	}
}
