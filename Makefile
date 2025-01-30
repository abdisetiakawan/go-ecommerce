.PHONY: run
run:
	@echo "Starting Go Ecommerce"
	go run cmd/web/main.go

.PHONY: migrate
migrate:
	@echo "Starting Migrate Ecommerce"
	go run db/db.go