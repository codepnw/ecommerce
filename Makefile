postgres:
	docker run --name ecommerce_dev -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:15.4


.PHONY: postgres