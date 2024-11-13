run:
	go run ./cmd/main.go

build:
	go build -o ./bin/main ./cmd/main.go

test:
	go test ./...

test-cover:
	go test -cover ./...

tidy:
	gofmt -w .
	go mod tidy