package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/deividr/zion-api/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func convertDistanceToInt8(distanceStr string) (int8, error) {
	cleanStr := strings.TrimSpace(strings.ReplaceAll(distanceStr, "km", ""))

	distance, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		return 0, fmt.Errorf("Error ao converter distância")
	}

	return int8(distance * 10), nil
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

	results, err2 := dbOldPool.Query("SELECT cd_endereco, cd_cliente, cd_cep, ds_logradouro, nr_logradouro, ds_bairro, ds_cidade, ds_uf, ds_complemento, ds_distancia FROM endereco WHERE ds_logradouro != '' AND cd_cep != '';")

	if err2 != nil {
		fmt.Println("Erro na query", err2)
		return
	}

	for results.Next() {
		var address domain.Address

		err3 := results.Scan(
			&address.OldId,
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
		if err3 != nil {
			fmt.Println("Erro no scan", err3)
			continue
		}

		distance := int8(0)

		if address.Distance != nil {
			result, err := convertDistanceToInt8(*address.Distance)
			if err != nil {
				continue
			}
			distance = result
		}

		err4 := dbNewPool.QueryRow(context.Background(), `SELECT id FROM customers WHERE old_id = $1`, address.Id).Scan(&address.Id)
		if err4 != nil {
			fmt.Println("Erro no select do customer", err4)
			continue
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO addresses (old_id, customer_id, cep, street, adr_number, neighborhood, city, adr_state, aditional_details, distance) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			address.OldId,
			address.Id,
			address.Cep,
			address.Street,
			address.Number,
			address.Neighborhood,
			address.City,
			address.State,
			address.AditionalDetails,
			distance,
		)

		if err != nil {
			fmt.Println("Erro ao inserir endereço no PostgreSQL: ", err)
			continue
		}

		fmt.Printf("Endereço inserido com sucesso: %+v\n", address.Street)
	}
}
