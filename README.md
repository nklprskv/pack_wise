# pack_wise

`pack_wise` is a full-stack pack optimization app:

- `backend/`: Go + Huma API + SQLC
- `frontend/`: React + Vite + Tailwind dashboard
- `db`: PostgreSQL

The service follows the rules:

1. Ship only whole packs.
2. Ship the smallest possible quantity that still fulfills the order.
3. If multiple combinations ship the same quantity, use the fewest packs.

## Live Demo

- Application: `https://packwise.up.railway.app/`
- API base: `https://api-packwise.up.railway.app/`
- Swagger UI: `https://api-packwise.up.railway.app/swagger`

## Stack

- Go
- Huma v2
- PostgreSQL
- SQLC
- golang-migrate
- React
- Vite
- Tailwind CSS
- lucide-react
- Docker Compose

## Project Structure

```text
.
├── backend/        # Go API, migrations, store, tests
├── frontend/       # React dashboard
├── docker-compose.yml
├── .env.example
└── Makefile
```

## Quick Start

Prerequisites:

- Docker
- Docker Compose

Create local env file:

```sh
cp .env.example .env
```

Start the full stack:

```sh
docker compose up --build
```

Or:

```sh
make up
```

## Services

- Frontend: `http://localhost:5173`
- Backend API: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/swagger`
- PostgreSQL: `localhost:5432`

## Docker Compose

`docker-compose.yml` starts 3 services:

1. `db` using PostgreSQL 17
2. `backend` using the code in `backend/`
3. `frontend` using the code in `frontend/`

Runtime wiring:

- `backend` depends on healthy `db`
- `frontend` depends on `backend`
- `frontend` uses `API_URL=http://backend:8080`
- browser traffic reaches the frontend on port `5173`

## Frontend

The frontend is a two-column dashboard:

- left: pack sizes list, add size, delete size
- right: calculation form and result table

Included UX:

- current pack configuration display
- `Challenge Case` preset for `23, 31, 53` and `500000`
- default reset button
- error state handling
- responsive layout for desktop and mobile

### Local Frontend Development

```sh
cd frontend
npm install
npm run dev
```

Vite proxy forwards API calls to the backend. The proxy target is configured in `frontend/vite.config.ts`.

## Backend

The backend exposes:

- `GET /health`
- `GET /api/v1/packs`
- `POST /api/v1/packs`
- `DELETE /api/v1/packs/{size}`
- `POST /api/v1/calculate`

### Local Backend Development

Start PostgreSQL first, then run:

```sh
cd backend
DATABASE_URL=postgres://packwise:packwise@localhost:5432/pack_wise?sslmode=disable go run .
```

The backend:

- runs migrations automatically on startup
- loads pack sizes from PostgreSQL
- serves Swagger UI generated from the Huma OpenAPI spec
- enables CORS for configured frontend origins

### API Docs

- Swagger UI: `/swagger`
- OpenAPI JSON: `/openapi.json`
- OpenAPI YAML: `/openapi.yaml`

## API Examples

### Health

```sh
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "ok"
}
```

### Get Pack Sizes

```sh
curl http://localhost:8080/api/v1/packs
```

Response:

```json
[250, 500, 1000, 2000, 5000]
```

### Update Pack Sizes

```sh
curl -X POST http://localhost:8080/api/v1/packs \
  -H 'Content-Type: application/json' \
  -d '{"sizes":[23,31,53]}'
```

### Calculate Packs

```sh
curl -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"items":501}'
```

Response:

```json
{
  "requestedItems": 501,
  "totalItems": 750,
  "packs": [
    { "size": 500, "quantity": 1 },
    { "size": 250, "quantity": 1 }
  ]
}
```

### Delete Pack Size

```sh
curl -X DELETE http://localhost:8080/api/v1/packs/31
```

## Challenge Case

Assignment edge case:

- pack sizes: `23, 31, 53`
- items: `500000`

Expected result:

```json
{
  "requestedItems": 500000,
  "totalItems": 500000,
  "packs": [
    { "size": 53, "quantity": 9429 },
    { "size": 31, "quantity": 7 },
    { "size": 23, "quantity": 2 }
  ]
}
```

You can reproduce it either from the UI via `Challenge Case` or through the API.

## Algorithm

The calculation uses dynamic programming instead of a greedy strategy.

That matters because pack sizes are configurable at runtime. Sets like `23, 31, 53` are not reliably solved by greedy selection, while DP guarantees:

1. minimal shipped quantity
2. minimal pack count for that quantity

## Testing

Backend tests:

```sh
cd backend
go test ./...
```

Frontend production build:

```sh
cd frontend
npm run build
```

## Useful Commands

```sh
make up
make down
make logs
make db-up
make migrate-up
make run
```
