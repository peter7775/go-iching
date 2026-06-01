# I Ching app template (Go + Fiber v2)

Základ projektu pro I-ťing / Yi Jing aplikaci s **Fiber v2** backendem, web klientem a desktop kostrou.

## Co obsahuje
- Fiber v2 HTTP API
- repository pro SQLite
- repository pro PostgreSQL
- fallback in-memory repository
- web MVP frontend
- Docker Compose pro Postgres
- náhodné generování hexagramu přes `crypto/rand`
- interpretaci nad Wilhelm-compatible dataset modelem
- Wails stub pro desktop klient

## Režimy spuštění
### 1) SQLite
```bash
export APP_STORAGE=sqlite
export APP_SQLITE_PATH=./iching.db
make run
```

### 2) PostgreSQL přes Docker Compose
```bash
make docker-up
export APP_STORAGE=postgres
export APP_PG_DSN='postgres://postgres:postgres@localhost:5432/iching?sslmode=disable'
make run
```

Pak otevři `http://localhost:8080`.
