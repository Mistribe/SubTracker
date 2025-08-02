.PHONY: install-deps build test run clean

install-deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/vektra/mockery/v3@v3.5.0
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

openapi:
	@swag fmt -d ./cmd/api/ &&\
	swag init -d ./cmd/api/ -o ./cmd/api/docs --v3.1 &&\
	cd ../scripts &&\
	./http-client-gen.sh

build:
	go build -v -o ./tmp/api ./cmd/api/

test:
	go test -v ./...

run:
	go run ./cmd/api/main.go

clean:
	go clean
	rm -f ./tmp/*


