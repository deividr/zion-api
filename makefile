-include .env

create_migration:
	migrate create -ext=sql -dir=internal/infra/database/migrations $(name)

migrate_up:
	migrate -path=internal/infra/database/migrations -database "${DATABASE_URL}" -verbose up

migrate_down:
	migrate -path=internal/infra/database/migrations -database "${DATABASE_URL}" -verbose down

.PHONY: create_migration migrate_up migrate_down
