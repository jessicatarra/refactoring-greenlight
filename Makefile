include .envrc

.PHONY: run/mono
run/mono:
	go run ./cmd/mono -db-dsn=${DATABASE_URL} -cors-trusted-origins=${CORS_TRUSTED_ORIGINS} -jwt-secret=${JWT_SECRET} -smtp-host=${SMTP_HOST} -smtp-password=${SMTP_PASSWORD} -smtp-username=${SMTP_USERNAME}

.PHONY: run/mono/help
run/mono/help:
	go run ./cmd/mono/ -help

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -tags auth ./... -v


.PHONY: ci/cd/audit
ci/cd/audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...

.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...' go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

.PHONY: generate/mono/docs
generate/mono/docs:
	@echo 'Remove docs...'
	rm -rf docs
	@echo 'Generate updated docs folder'
	swag init -d cmd/mono,ms/auth/internal/ --parseDependency --parseInternal --parseDepth 2


.PHONY: generate/auth/mocks
generate/auth/mocks:
	@echo 'Remove mocks...'
	rm -rf ms/auth/internal/domain/mocks
	@echo 'Generate updated mocks...'
	mockery --all --output=ms/auth/internal/domain/mocks --dir=ms/auth/internal/domain
	mockery --all --output=ms/auth/internal/service/mocks --dir=ms/auth/internal/service


.PHONY: run/test/coverage
run/test/coverage:
	@echo 'Run test of auth tags and generate coverage.out file'
	go test ./... -tags auth -coverprofile=coverage.out
	@echo 'Generate coverage.html file'
	 go tool cover -html=coverage.out -o coverage.html

# ====================================================================================
# # BUILD
# ==================================================================================== #
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-X main.version=${git_description}'

.PHONY: build/mono
build/mono:
	@echo 'Building cmd/mono...'
	go build -ldflags=${linker_flags} -o=./bin/mono ./cmd/mono
