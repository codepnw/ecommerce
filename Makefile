postgres:
	docker run --name ecommerce_dev -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:15.4

migrateup:
	migrate -path pkg/database/migrations -database "postgres://postgres:123456@localhost:4444/ecommerce_dev?sslmode=disable" -verbose up
migratedown:
	migrate -path pkg/database/migrations -database "postgres://postgres:123456@localhost:4444/ecommerce_dev?sslmode=disable" -verbose down

.PHONY: postgres