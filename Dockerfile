# Build stage
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine3.17 AS builder
WORKDIR /bin
COPY . .
RUN go mod download
ARG API_VERSION="v0.0.0+unknown"
RUN go build -ldflags "-s -w -X 'main.version=$API_VERSION'" -o mono ./cmd/mono

# Run stage
FROM golang:${GO_VERSION}-alpine3.17 AS build-release-stage
WORKDIR /bin
COPY --from=builder /bin/.envrc .
COPY --from=builder /bin/mono .

EXPOSE 8080
EXPOSE 8082

#TODO: create .envrc via github action and inject variables using github environment secrets
#TODO: add go run flags depending on the environment
ENTRYPOINT ["/bin/sh", "-c", "source .envrc && /bin/mono -cors-trusted-origins=\"$CORS_TRUSTED_ORIGINS\" -db-dsn=\"$DATABASE_URL\""]
