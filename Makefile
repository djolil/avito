build:
	go build -o ./bin/avito-banner ./cmd/main.go

run:
	docker compose up -d avito-banner

migrate:
	docker compose up -d migrate

e2e:
	go test ./tests -timeout 30s

all: migrate run