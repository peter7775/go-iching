# I Ching app template (Go + Fiber)

Základ projektu pro I-ťing / Yi Jing aplikaci s Fiber backendem, web klientem a desktop kostrou.

## Co obsahuje
- Fiber HTTP API
- repository pro SQLite
- repository pro PostgreSQL
- fallback in-memory repository
- web MVP frontend
- Docker Compose pro Postgres
- importní příkaz pro dataset hexagramů
- Wails stub pro desktop klient

## Režimy spuštění
### 1) SQLite (nejjednodušší lokálně)
```bash
export APP_STORAGE=sqlite
export APP_SQLITE_PATH=./iching.db
go run ./cmd/api
```
Pak otevři http://localhost:8080

### 2) PostgreSQL přes Docker Compose
```bash
docker compose up -d
export APP_STORAGE=postgres
export APP_PG_DSN='postgres://postgres:postgres@localhost:5432/iching?sslmode=disable'
go run ./cmd/api
```

Pak otevři `http://localhost:8080`.

## Poznámka
Mapování všech 64 hexagramů ještě není plně implementované; importer a datové modely jsou připravené, aby šlo snadno doplnit kompletní dataset.
