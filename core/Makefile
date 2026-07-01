.PHONY: run worker test generate migrate-up migrate-down

run:
	go run cmd/api/main.go

worker:
	go run cmd/worker/main.go

test:
	go test -v ./...

generate:
	sqlc generate

migrate-up:
	migrate -path migrations -database "${DATABASE_URL}" -verbose up

migrate-down:
	migrate -path migrations -database "${DATABASE_URL}" -verbose down
