package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Customer struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Phone2    *string   `json:"phone2"`
	Email     *string   `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	OldId     string    `json:"oldId"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erro ao carregar .env: ", err)
		return
	}

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

	results, err2 := dbOldPool.Query("SELECT cd_cliente, nm_cliente, nr_telefone1, nr_telefone2, ds_email, dt_criacao FROM cliente;")

	if err2 != nil {
		fmt.Println("Erro na query", err2)
		return
	}

	for results.Next() {
		var customer Customer
		var createdAtRaw []uint8 // variável temporária para receber o dado bruto

		err3 := results.Scan(
			&customer.OldId,
			&customer.Name,
			&customer.Phone,
			&customer.Phone2,
			&customer.Email,
			&createdAtRaw,
		)
		if err3 != nil {
			fmt.Println("Erro no scan", err3)
			continue
		}

		// Converter []uint8 para time.Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtRaw))
		if err != nil {
			fmt.Println("Erro ao converter data:", err)
			continue
		}
		customer.CreatedAt = createdAt

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO customers (name, phone, phone2, email, old_id, created_at) 
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			customer.Name,
			customer.Phone,
			customer.Phone2,
			customer.Email,
			customer.OldId,
			customer.CreatedAt,
		)

		if err != nil {
			fmt.Println("Erro ao inserir cliente no PostgreSQL: ", err)
			continue
		}

		fmt.Printf("Cliente inserido com sucesso: %+v\n", customer.Name)
	}
}
