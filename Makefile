# Variables
APP_NAME=devboard
MAIN_PATH=./cmd/api
MIGRATE_PATH=./migrations
DB_URL=postgresql://postgres:password@localhost:5432/devboard?sslmode=disable
COVERAGE_PROFILE=coverage.out
COVERAGE_MIN=90.0

.PHONY: run build fmt fmt-check verify vet test lint ci install-hooks migrate-up migrate-down generate tidy help docker-up docker-down docker-logs

## run: correr la aplicación
run:
	go run $(MAIN_PATH)/main.go

## build: compilar el binario
build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

## fmt: formatear el código Go
fmt:
	gofmt -w .

## fmt-check: comprobar el formato del código Go
fmt-check:
	test -z "$(gofmt -l .)"

## verify: verificar la integridad de los módulos
verify:
	go mod verify

## vet: analizar errores sospechosos en el código
vet:
	go vet ./...

## test: correr todos los tests con race detector y barra de progreso
test:
	@echo "🧪 Ejecutando tests..."
	@echo ""
	@./scripts/run_tests.sh
	@echo ""
	@coverage=$$(go tool cover -func=coverage.out | awk '/^total:/ { sub(/%/, "", $$NF); print $$NF }'); \
	test -n "$$coverage" || { echo "❌ Coverage could not be measured"; exit 1; }; \
	awk -v coverage="$$coverage" -v minimum="$(COVERAGE_MIN)" 'BEGIN { if (coverage < minimum) exit 1 }' || { echo "❌ Coverage $${coverage}% is below minimum $(COVERAGE_MIN)%"; exit 1; }; \
	echo "✅ Coverage: $${coverage}% (minimum $(COVERAGE_MIN)%)"

## lint: analizar el código con golangci-lint
lint:
	golangci-lint run ./...

## install-hooks: activar los hooks de Git del repositorio
install-hooks:
	git config core.hooksPath .githooks
	@echo "Git hooks enabled from .githooks"

## migrate-up: aplicar todas las migraciones pendientes
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: revertir la última migración
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## generate: correr go generate en todo el proyecto
generate:
	go generate ./...

## tidy: limpiar y verificar dependencias
tidy:
	go mod tidy
	go mod verify

## ci: ejecutar las comprobaciones y tests del pipeline
ci: fmt-check verify vet test

## help: mostrar este menú
help:
	@grep -E '^##' Makefile | sed 's/## //'

## docker-up: levantar los servicios de desarrollo
docker-up:
	docker compose up -d

## docker-down: detener los servicios de desarrollo
docker-down:
	docker compose down

## docker-logs: ver logs de todos los servicios
docker-logs:
	docker compose logs -f

## docs: generar documentación OpenAPI con swaggo
docs: 
	swag init -g cmd/api/main.go -o docs