# utic.dev-server API

API para el Club de Programación de la UTIC. Los estudiantes pueden registrar sus repos de GitHub y la IA analiza automáticamente la info del proyecto.

## Base URL

```
https://tu-railway-url.up.railway.app
```

## Autenticación

Todos los endpoints protegidos requieren el header `X-API-Key`:

```bash
curl -H "X-API-Key: tu-api-key" https://tu-url/api/projects
```

## Endpoints

### Health Check (público)

```
GET /health
```

**Response:** `200 OK`

---

### Listar Proyectos (protegido)

```
GET /api/projects
```

**Response:**
```json
[
  {
    "ID": "c1f40078-e800-4eef-8f03-eb27f3241a20",
    "Name": "Golang Walking Skeleton",
    "Description": "Backend template con Clean Architecture y multi-tenancy",
    "StudentEmail": "hugo@email.com",
    "StudentName": "Hugo Britez",
    "TechStack": "[\"Go\", \"chi\", \"PostgreSQL\"]",
    "Categories": "[\"backend\", \"api\"]",
    "Stars": 12,
    "RepoUrl": "https://github.com/HugoBritez/golang-walking-skeleton",
    "CreatedAt": "2026-05-15T18:34:18Z"
  }
]
```

> `TechStack` y `Categories` son JSON strings. Hacé `JSON.parse()` en el front.

---

### Obtener Proyecto por ID (protegido)

```
GET /api/projects/{id}
```

**Response:** `200 OK` con el proyecto o `404` si no existe.

---

### Crear / Actualizar Proyecto (protegido)

```
POST /api/projects
Content-Type: application/json
X-API-Key: tu-api-key

{
  "repo_url": "https://github.com/usuario/mi-repo"
}
```

**Qué pasa:**
1. Se busca el repo en la DB por `repo_url`
2. Si ya existe → se actualiza con la info fresca
3. Si no existe → se crea nuevo
4. La IA analiza el README y metadata del repo para extraer: nombre, descripción, tech stack, categorías, estrellas, email y nombre del autor

**Response:** `201 Created` con el proyecto creado/actualizado.

**Error:** `400` si falta `repo_url`, `500` si falla la IA.

---

### Dar Star a un Proyecto (protegido)

```
POST /api/projects/{id}/star
```

**Response:** `200 OK` con el proyecto actualizado (stars incrementado en 1).

---

## Ejemplos con fetch (JavaScript)

```js
const API_URL = "https://tu-url.up.railway.app";
const API_KEY = "tu-api-key";

const headers = {
  "Content-Type": "application/json",
  "X-API-Key": API_KEY,
};

// Listar proyectos
const projects = await fetch(`${API_URL}/api/projects`, { headers })
  .then(r => r.json());

// Crear proyecto
const project = await fetch(`${API_URL}/api/projects`, {
  method: "POST",
  headers,
  body: JSON.stringify({ repo_url: "https://github.com/usuario/repo" }),
}).then(r => r.json());

// Dar star
await fetch(`${API_URL}/api/projects/${projectId}/star`, {
  method: "POST",
  headers,
});
```
