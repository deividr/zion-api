-include .env

create_migration:
	migrate create -ext=sql -dir=internal/infra/database/migrations $(name)

migrate_up:
	migrate -path=internal/infra/database/migrations -database "${DATABASE_URL_MIGRATE}" -verbose up

migrate_down:
	migrate -path=internal/infra/database/migrations -database "${DATABASE_URL_MIGRATE}" -verbose down

load_products:
	go run cmd/scripts/load/products/main.go

load_customers:
	go run cmd/scripts/load/customers/main.go

load_addresses:
	go run cmd/scripts/load/address/main.go

.PHONY: create_migration migrate_up migrate_down load_products load_customers load_addresses
