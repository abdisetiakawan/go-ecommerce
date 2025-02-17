.PHONY: run
run:
	@echo "Starting Go Ecommerce"
	go run cmd/web/main.go

.PHONY: migrate-gorm
migrate:
	@echo "Starting Migrate Ecommerce"
	go run db/db.go

.PHONY: migrate-up
migrate-up:
	@echo "Starting Migrate Up Ecommerce"
	migrate -database "mysql://root:password@tcp(localhost:3306)/mydb" -path db/migrations/sql up
