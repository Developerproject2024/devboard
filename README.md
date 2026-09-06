# Devboard

Backend en Go para gestionar tareas y flujo de trabajo de un tablero.

## Requisitos

- Go 1.27+
- Git
- Docker (opcional)

## Comandos del Makefile

```bash
# Ejecutar la aplicación
make run

# Compilar binario
make build

# Formatear código
make fmt

# Verificar formato
make fmt-check

# Verificar dependencias
make verify

# Ejecutar análisis estático
make vet

# Ejecutar tests
make test

# Ejecutar lint
make lint

# Activar hooks de git
make install-hooks

# Aplicar migraciones
make migrate-up

# Revertir migración
make migrate-down

# Generar código
make generate

# Limpiar dependencias
make tidy

# Validación de pipeline CI
make ci

# Mostrar ayuda
make help

# Levantar servicios Docker
make docker-up

# Detener servicios Docker
make docker-down

# Ver logs Docker
make docker-logs

# Generar Swagger
make docs
```

## Ejecutar localmente

```bash
go mod download
make run
```

## Endpoints actuales

- http://localhost:8080/health
- http://localhost:8080/docs/

## Verificar salud

```bash
curl http://localhost:8080/health
```

Respuesta esperada:

```json
{"status":"Fabio Arango","version":"1.1.0"}
```

## Estructura básica

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── handler/
│   ├── logger/
│   ├── middlewares/
│   ├── server/
│   └── validator/
├── docs/
├── scripts/
├── go.mod
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── README.md
└── render.yaml
```
