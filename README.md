# utic.dev-server

API para el Club de Programación de la UTIC. Los estudiantes pueden registrar sus repos de GitHub y la IA analiza automáticamente la información del proyecto.

## Tech Stack

- **Go** con [chi](https://github.com/go-chi/chi) router
- **SQLite** con [sqlc](https://sqlc.dev/) para queries type-safe
- **Groq** (llama-3.3-70b) para análisis de repos con IA
- Arquitectura limpia (domain/application/infrastructure)

## Estructura

```
├── cmd/                          # Entry points
├── config/                       # Configuración
├── db/
│   ├── schema/                   # Migraciones SQL
│   └── queries/                  # Queries para sqlc
├── internal/
│   ├── domain/
│   │   ├── entities/             # Entidades del dominio
│   │   ├── repositories/         # Interfaces de repositorio
│   │   └── services/             # Interfaces de servicios (IA)
│   ├── application/
│   │   └── projects/             # Casos de uso
│   └── infrastructure/
│       ├── database/             # Conexión SQLite
│       ├── db/                   # Código generado por sqlc
│       ├── http/                 # Handlers HTTP
│       ├── middleware/           # Middleware (API Key)
│       ├── repository/           # Implementación de repos
│       └── services/             # Cliente Groq
├── main.go
├── sqlc.yaml
└── Dockerfile
```

## Requisitos

- Go 1.25+
- [sqlc](https://sqlc.dev/docs/install/)
- Una API key de [Groq](https://console.groq.com/keys)

## Setup Local

```bash
# 1. Clonar y entrar
git clone https://github.com/HugoBritez/utic_server.git
cd utic.dev-server

# 2. Instalar dependencias
go mod download

# 3. Crear .env
cp .env.example .env

# 4. Editar .env con tus credenciales
# API_KEY=tu-api-key
# AI_API_KEY=tu-groq-api-key

# 5. Generar código de sqlc
sqlc generate

# 6. Correr
go run main.go
```

## Variables de Entorno

| Variable | Descripción | Default |
|----------|-------------|---------|
| `API_KEY` | Key para proteger endpoints | - |
| `AI_API_KEY` | API key de Groq | - |
| `AI_MODEL` | Modelo de Groq | `llama-3.3-70b-versatile` |
| `DB_PATH` | Ruta de la base de datos SQLite | `./data/app.db` |
| `PORT` | Puerto del servidor | `3000` |

## API

Ver [API.md](API.md) para documentación completa.

### Endpoints principales

| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/projects` | Crear/actualizar proyecto (IA analiza el repo) |
| `GET` | `/api/projects` | Listar proyectos |
| `GET` | `/api/projects/{id}` | Obtener proyecto por ID |
| `POST` | `/api/projects/{id}/star` | Dar star a un proyecto |

Todos los endpoints de `/api/*` requieren el header `X-API-Key`.

### Ejemplo

```bash
curl -X POST http://localhost:3000/api/projects \
  -H "X-API-Key: tu-api-key" \
  -H "Content-Type: application/json" \
  -d '{"repo_url": "https://github.com/usuario/mi-repo"}'
```

## sqlc

Para regenerar el código Go después de cambiar queries:

```bash
sqlc generate
```

## Docker

```bash
docker build -t utic-dev-server .
docker run -p 3000:3000 --env-file .env utic-dev-server
```

## Deploy en Railway

1. Conectar el repo a Railway
2. Agregar las variables de entorno en el dashboard
3. Deploy automático al pushear a `main`

## Licencia

MIT
