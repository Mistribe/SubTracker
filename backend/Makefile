.PHONY: install-deps build test run clean

install-deps:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/vektra/mockery/v3@v3.5.0

build-swagger:
	@swag init -d ./cmd/api/ -o ./cmd/api/docs
build:
	go build -v ./cmd/api/main.go -o ./tmp/api

test:
	go test -v ./...

run:
	go run ./cmd/api/main.go

clean:
	go clean
	rm -f ./tmp/*


