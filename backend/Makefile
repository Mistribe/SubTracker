.PHONY: install-deps build test run clean

install-deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/vektra/mockery/v3@v3.5.3
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/go-jet/jet/v2/cmd/jet@latest

openapi:
	swag fmt -d ./cmd/api/ -d ./internal/adapters/http/handlers/  &&\
	swag init -d ./cmd/api/ -pdl 3 -o ./internal/adapters/http/handlers/openapi --v3.1

clients: typescript-client dart-client

typescript-client:
	@echo "Clean old client api" &&\
		rm -fr "${PWD}/../web/src/api" &&\
    	mkdir -p "${PWD}/../web/src/api" &&\
    	cp "${PWD}/internal/adapters/http/handlers/openapi/swagger.yaml" "${PWD}/../web/src/api"
	@docker pull mcr.microsoft.com/openapi/kiota && \
    docker run --rm -v ../web/src/api:/app/output \
    -v ./internal/adapters/http/handlers/openapi/swagger.yaml:/app/openapi.yaml \
    --user $(id -u):$(id -g) \
    mcr.microsoft.com/openapi/kiota generate --language typescript --clean-output

dart-client:
	docker pull mcr.microsoft.com/openapi/kiota && \
    docker run --rm -v ../mobile/lib/api:/app/output \
    -v ./internal/adapters/http/handlers/openapi/swagger.yaml:/app/openapi.yaml \
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


