FROM --platform=${BUILDPLATFORM} golang:1.25 AS builder
ARG APPUSER=appuser
ARG TARGETOS
ARG TARGETARCH
ENV USER=${APPUSER}
ENV UID=1001
WORKDIR /source
RUN mkdir app && \
    adduser --disabled-password --gecos "" --no-create-home --uid "${UID}" "${USER}"
COPY go.* ./
RUN --mount=type=cache,target=/root/go/pkg go mod download

COPY ["cmd/", "cmd/"]
COPY ["pkg/", "pkg/"]
COPY ["internal/", "internal/"]

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -extldflags '-static'" -trimpath -o ./app/main ./cmd/api/

FROM scratch
ARG APPUSER=appuser
ENV UID=1001
EXPOSE 8080

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder --chown=${APPUSER}:${APPUSER} /source/app /app
WORKDIR /app
USER ${APPUSER}:${APPUSER}

ENTRYPOINT ["./main"]