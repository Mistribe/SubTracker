FROM golang:latest

# Build goose with only PostgreSQL enabled by excluding all other drivers
RUN git clone https://github.com/pressly/goose.git /tmp/goose && \
    cd /tmp/goose && \
    go build -tags='no_clickhouse no_libsql no_mssql no_mysql no_sqlite3 no_vertica no_ydb' -o /usr/local/bin/goose ./cmd/goose && \
    rm -rf /tmp/goose

WORKDIR /app/database
COPY ["*.sql", "."]

ENTRYPOINT [ "goose", "up" ]