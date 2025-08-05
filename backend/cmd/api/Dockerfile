FROM --platform=${BUILDPLATFORM} golang:1.24 AS builder
ARG APPUSER=appuser
#ARG GITHUB_TOKEN
ARG TARGETOS
ARG TARGETARCH
ENV USER=${APPUSER}
ENV UID=1001
WORKDIR /source
RUN mkdir app && \
#    apt-get update && apt-get install -y build-essential && \
    adduser --disabled-password --gecos "" --no-create-home --uid "${UID}" "${USER}"
#RUN git config --global "url.https://dummy:${GITHUB_TOKEN}@github.com/.insteadOf" "https://github.com/"
COPY go.* ./
RUN --mount=type=cache,target=/root/go/pkg go mod download
COPY ["cmd/", "cmd/"]
COPY ["pkg/", "pkg/"]
COPY ["internal/", "internal/"]
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -o ./app/main ./cmd/api/

FROM alpine:3.22
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