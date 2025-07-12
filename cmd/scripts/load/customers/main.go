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
	"github.com/schollz/progressbar/v3"
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

	// Contar total de registros primeiro
	var totalCount int
	err = dbOldPool.QueryRow("SELECT COUNT(*) FROM cliente").Scan(&totalCount)
	if err != nil {
		fmt.Println("Erro ao contar registros:", err)
		return
	}

	fmt.Printf("Total de clientes a serem processados: %d\n", totalCount)

	results, err := dbOldPool.Query("SELECT cd_cliente, nm_cliente, nr_telefone1, nr_telefone2, ds_email, dt_criacao FROM cliente;")
	if err != nil {
		fmt.Println("Erro na query", err)
		return
	}

	var customers []Customer

	for results.Next() {
		var customer Customer
		var createdAtRaw []uint8

		err3 := results.Scan(
			&customer.OldId,
			&customer.Name,
			&customer.Phone,
			&customer.Phone2,
			&customer.Email,
			&createdAtRaw,
		)
		if err3 != nil {
			continue
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtRaw))
		if err != nil {
			continue
		}

		customer.CreatedAt = createdAt
		customers = append(customers, customer)
	}

	results.Close()
	dbOldPool.Close()

	bar := progressbar.NewOptions(len(customers),
		progressbar.OptionSetDescription("Processando clientes"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		// progressbar.OptionSetTheme(progressbar.Theme{
		// 	Saucer:        "[green]=[reset]",
		// 	SaucerHead:    "[green]>[reset]",
		// 	SaucerPadding: " ",
		// 	BarStart:      "[",
		// 	BarEnd:        "]",
		// }),
	)

	successCount := 0
	errorCount := 0

	for _, customer := range customers {
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
			fmt.Printf("Erro ao inserir cliente no PostgreSQL: %v\n", err)
			break
		} else {
			successCount++
		}

		bar.Add(1)
	}

	bar.Finish()

	fmt.Printf("\n\n=== ESTATÍSTICAS FINAIS ===\n")
	fmt.Printf("Total processados: %d\n", totalCount)
	fmt.Printf("Sucessos: %d\n", successCount)
	fmt.Printf("Erros: %d\n", errorCount)
	fmt.Printf("Taxa de sucesso: %.2f%%\n", float64(successCount)/float64(totalCount)*100)
}
