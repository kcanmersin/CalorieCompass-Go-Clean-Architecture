.PHONY: build run clean test migrate-up migrate-down swag

build:
	go build -o bin/app cmd/app/main.go

run:
	go run cmd/app/main.go

clean:
	rm -rf bin/

test:
	go test -v ./...

migrate-up:
	migrate -path migrations -database "$(POSTGRES_URL)" up

migrate-down:
	migrate -path migrations -database "$(POSTGRES_URL)" down

swag:
	swag init -g cmd/app/main.go

dev:
	air

.PHONY: build run clean test migrate-up migrate-down swag dev