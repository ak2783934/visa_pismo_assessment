run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

docker:
	docker compose up --build

swagger:
	swag init -g cmd/server/main.go -o docs --parseInternal
