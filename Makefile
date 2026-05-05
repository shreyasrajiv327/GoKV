run:
	go run ./cmd/server/main.go

test:
	go test ./...

race:
	go test -race ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy