build:
	@go build -o bin/Back-end-spotychafa-go cmd/main.go

test:
	@go test -v ./...
	
run: build
	@./bin/Back-end-spotychafa-go

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down