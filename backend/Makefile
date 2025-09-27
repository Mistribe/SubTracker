.PHONY: install-deps build test run clean

install-deps:
	go install github.com/swaggo/swag/v2/cmd/swag@latest
	go install github.com/vektra/mockery/v3@v3.5.5
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/go-jet/jet/v2/cmd/jet@latest

openapi:
	swag fmt -d ./cmd/api/ -d ./internal/adapters/http/handlers/  &&\
	swag init -d ./cmd/api/ -pdl 3 -o ./internal/adapters/http/handlers/openapi --v3.1

clients: typescript-client dart-client

typescript-client:
	@echo "Clean old client api" &&\
		rm -fr "${PWD}/../web/src/api" &&\
    	mkdir -p "${PWD}/../web/src/api"
	@docker pull openapitools/openapi-generator-cli && \
    docker run --rm \
    -v ${PWD}/../web/src/api:/app/output \
    -v ${PWD}/internal/adapters/http/handlers/openapi/swagger.yaml:/app/openapi.yaml \
    --user $(id -u):$(id -g) \
    openapitools/openapi-generator-cli generate \
    -i /app/openapi.yaml \
    -g typescript-fetch \
    --global-property=apiDocs=false,modelDocs=false,apiTests=false,modelTests=false \
    -o /app/output

dart-client:
	@echo "Clean old dart client api" &&\
		rm -fr "${PWD}/../mobile/lib/http_client" &&\
		mkdir -p "${PWD}/../mobile/lib/http_client"
	@docker pull openapitools/openapi-generator-cli && \
    docker run --rm \
    -v ${PWD}/../mobile/lib/http_client_tmp:/app/output \
    -v ${PWD}/internal/adapters/http/handlers/openapi/swagger.yaml:/app/openapi.yaml \
    --user $(id -u):$(id -g) \
    openapitools/openapi-generator-cli generate \
    -i /app/openapi.yaml \
    -g dart \
    --global-property=apiDocs=false,modelDocs=false,apiTests=false,modelTests=false \
    -o /app/output && \
    cp -r "${PWD}/../mobile/lib/http_client_tmp/lib/"  "${PWD}/../mobile/lib/http_client" && \
    rm -fr "${PWD}/../mobile/lib/http_client_tmp"

run-deps:
	@docker compose down -v && docker compose up --build database database.migration

sql-jet:
	@jet -dsn=postgresql://postgres:postgres@localhost:5432/app?sslmode=disable -schema=public -path=./internal/adapters/persistence/db/jet

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


