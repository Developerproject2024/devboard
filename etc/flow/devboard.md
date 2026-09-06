# Devboard — Visión de alto nivel

## Descripción

Devboard es un backend en Go pensado para gestionar tareas y flujos de trabajo de un tablero colaborativo. La aplicación sigue una estructura modular orientada a separar dominio, servidor, handlers, validaciones y soporte transversal.

El proyecto se encuentra en una etapa inicial, pero su base está diseñada para crecer sin mezclar responsabilidades entre capas.

## Objetivo

Proporcionar una API REST con una estructura limpia para:

- gestionar tareas
- manejar estados de trabajo
- mantener validaciones del dominio
- incorporar futuras funcionalidades como persistencia, usuarios y notificaciones

## Arquitectura

```mermaid
flowchart LR
    Client[Cliente]
    API[Servidor HTTP]
    Router[Rutas]
    MW[Middlewares]
    Handler[Handlers]
    UseCase[Casos de uso]
    Domain[Dominio]
    Valid[Validaciones]
    Logger[Logger]
    Notify[Notificaciones]
    DB[(Base de datos futura)]

    Client --> API
    API --> Router
    Router --> MW
    MW --> Handler
    Handler --> UseCase
    UseCase --> Domain
    Handler --> Valid
    UseCase --> Notify
    API --> Logger
    UseCase --> DB
```

### Capas principales

- Entrada HTTP: servidor y enrutado
- Handlers: reciben y responden peticiones
- Casos de uso: coordinación de la lógica
- Dominio: entidades, reglas y estados
- Utilidades: logger, validación y extensiones futuras

## Secuencia general

```mermaid
sequenceDiagram
    participant C as Cliente
    participant S as Server
    participant M as Middlewares
    participant H as Handler
    participant U as Use Case
    participant D as Dominio

    C->>S: Petición HTTP
    S->>M: Ejecuta middleware
    M->>H: Invoca handler
    H->>U: Delegación de lógica
    U->>D: Aplica reglas del negocio
    D-->>U: Resultado
    U-->>H: Respuesta
    H-->>C: Payload HTTP
```

## Validaciones relevantes

El dominio ya incluye una validación clara para los estados de tarea:

- todo
- in_progress
- done

La validación se realiza en el tipo TaskStatus y evita aceptar valores fuera del flujo definido del negocio.

## Notas de diseño

- El proyecto mantiene una separación entre dominio e infraestructura.
- Los middlewares centralizan logging y manejo de errores.
- La API está preparada para evolucionar hacia CRUD real, persistencia y más casos de uso.
- La documentación Swagger y la estructura modular respaldan una base clara para crecer.

## Resumen

Devboard tiene una base tecnológica limpia y una arquitectura pensada para escalar. Actualmente cubre la capa fundamental de servidor, middleware y validación del dominio, y está listo para seguir incorporando más funcionalidades de negocio.
