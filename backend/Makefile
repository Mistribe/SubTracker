.PHONY: install-deps build test run clean

install-deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/vektra/mockery/v3@v3.5.0
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/go-jet/jet/v2/cmd/jet@latest

openapi:
	@swag fmt -d ./cmd/api/ &&\
	swag init -d ./cmd/api/ -o ./cmd/api/docs --v3.1

clients: typescript-client dart-client

typescript-client:
	docker pull mcr.microsoft.com/openapi/kiota && \
    docker run --rm -v ../web/src/api:/app/output \
    -v ./cmd/api/docs/swagger.yaml:/app/openapi.yaml \
    --user $(id -u):$(id -g) \
    mcr.microsoft.com/openapi/kiota generate --language typescript --clean-output

dart-client:
	docker pull mcr.microsoft.com/openapi/kiota && \
    docker run --rm -v ../mobile/lib/api:/app/output \
    -v ./cmd/api/docs/swagger.yaml:/app/openapi.yaml \
    --user $(id -u):$(id -g) \
    mcr.microsoft.com/openapi/kiota generate --language dart --clean-output

run-deps:
	@docker compose down -v && docker compose up --build database database.migration

kinde:
	@rm -fr ./internal/infrastructure/kinde/gen
	@docker run --rm \
    -v ${PWD}:/local openapitools/openapi-generator-cli generate \
    -i /local/internal/infrastructure/kinde/kinde-management-api.yaml \
    -g go \
    --additional-properties=withGoMod=false \
    --global-property=apiDocs=false \
    --global-property=modelDocs=false \
    --global-property=apiTests=false \
    --global-property=modelTests=false \
    -o /local/internal/infrastructure/kinde/gen

sql-jet:
	@jet -dsn=postgresql://postgres:postgres@localhost:5432/app?sslmode=disable -schema=public -path=./internal/infrastructure/persistence/jet

prepare: install-deps openapi clients

build:
	go build -v -o ./tmp/api ./cmd/api/

test:
	go test -v ./...

run:
	go run ./cmd/api/main.go

clean:
	go clean
	rm -f ./tmp/*


