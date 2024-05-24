build:
	@go build -o "beli-mang-be" cmd/main.go

migrate-up:
	migrate -database "postgres://postgres:admin@localhost:5432/beli-mang?sslmode=disable" -path internal/db/migrations up

migrate-down:
	migrate -database "postgres://postgres:admin@localhost:5432/beli-mang?sslmode=disable" -path internal/db/migrations down

run:
	./beli-mang-be

docker-build:
	docker build --no-cache -t beli-mang-be .

docker-run:
	docker compose -f "docker-compose.yaml" up -d --build