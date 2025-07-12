package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/deividr/zion-api/internal/domain"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func convertDistanceToInt8(distanceStr string) (int, error) {
	cleanStr := strings.TrimSpace(strings.ReplaceAll(distanceStr, "km", ""))
	distance, err := strconv.ParseFloat(cleanStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error ao converter distância: %v", err)
	}

	return int(math.Ceil(distance)), nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("erro ao carregar .env: %v", err)
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

	var oldAddresses []domain.Address

	for results.Next() {
		var address domain.Address
		var distanceStr *string

		err := results.Scan(
			&address.OldId,
			&address.CustomerId,
			&address.Cep,
			&address.Street,
			&address.Number,
			&address.Neighborhood,
			&address.City,
			&address.State,
			&address.AditionalDetails,
			&distanceStr,
		)
		if err != nil {
			continue
		}

		if distanceStr != nil {
			result, err := convertDistanceToInt8(*distanceStr)
			if err != nil {
				continue
			}
			address.Distance = &result
		} else {
			zero := 0
			address.Distance = &zero
		}

		oldAddresses = append(oldAddresses, address)
	}

	results.Close()
	dbOldPool.Close()

	bar := progressbar.NewOptions(len(oldAddresses),
		progressbar.OptionSetDescription("Processando clientes"),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)

	for _, address := range oldAddresses {
		err4 := dbNewPool.QueryRow(context.Background(), `SELECT id FROM customers WHERE old_id = $1`, address.CustomerId).Scan(&address.CustomerId)
		if err4 != nil {
			fmt.Println("Erro no select do customer", err4)
			break
		}

		_, err = dbNewPool.Exec(context.Background(),
			`INSERT INTO addresses (old_id, customer_id, cep, street, number, neighborhood, city, state, aditional_details, distance, is_default) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			address.OldId,
			address.CustomerId,
			address.Cep,
			address.Street,
			address.Number,
			address.Neighborhood,
			address.City,
			address.State,
			address.AditionalDetails,
			address.Distance,
			true,
		)
		if err != nil {
			fmt.Println("Erro ao inserir endereço no PostgreSQL: ", err)
			break
		}

		bar.Add(1)
	}
}
