# I Ching Fiber App

A Go application for I Ching divination with a Fiber-based HTTP API, a browser frontend, and persistent reading storage. The current project supports SQLite, PostgreSQL, and in-memory repositories, and stores each reading together with the original question, line data, interpretation payload, and optional reflection fields. [file:248][file:249][file:258]

## Overview

The application is structured as a small web app. A Go backend handles reading generation and persistence, while the frontend is served as static files from `web/static`. The API currently exposes endpoints for listing readings, fetching a reading by ID, creating a reading, and creating a randomized reading. [file:249]

Each stored reading includes:
- question text,
- cast method,
- language,
- raw line values,
- primary and relating hexagram numbers,
- serialized interpretation,
- optional reflection rating, note, and timestamp,
- creation timestamp. [file:258]

## Architecture

Main components in the current project:

- `cmd/api` — application entrypoint and startup logic. [file:248]
- `internal/httpfiber` — Fiber app configuration and HTTP routes. [file:249]
- `internal/service` — application and reading logic. [file:248]
- `internal/storage/sqlite` — SQLite persistence for saved readings. [file:248][file:258]
- `internal/storage/postgres` — PostgreSQL repository option. [file:248]
- `internal/storage/memory` — in-memory fallback repository. [file:248]
- `web/static` — browser UI served by the backend. [file:249]

## Storage backends

The repository is selected by application configuration. If `cfg.Storage` is set to `sqlite`, the app opens `cfg.SQLitePath` and uses the SQLite repository. If it is set to `postgres`, the app connects through `cfg.PostgresDSN`. Any other value falls back to in-memory storage. [file:248]

### SQLite

The SQLite repository creates a `readings` table if it does not exist and stores both structured fields and serialized JSON payloads. The current schema includes `question`, `method`, `language`, `lines_json`, `primary_number`, `relating_number`, `changing_lines_json`, `interpretation_json`, reflection fields, and `created_at`. [file:258]

Reads are ordered by `created_at DESC`, which means the latest saved readings appear first in history. The repository also supports loading a single reading by ID and updating reflection data with `SaveReflection(...)`. [file:258]

### PostgreSQL

PostgreSQL is available as an alternative repository and is activated when the storage mode is `postgres`. The application opens the connection using `sql.Open("pgx", cfg.PostgresDSN)`. [file:248]

### In-memory mode

In-memory storage is only a fallback. It is useful for tests or temporary runs, but readings are not persisted across process restarts. [file:248]

## API endpoints

The current HTTP layer exposes these routes:

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Returns a simple health response. [file:249] |
| GET | `/api/readings` | Returns the saved reading history. [file:249] |
| GET | `/api/readings/:id` | Returns one saved reading by ID. [file:249] |
| POST | `/api/readings` | Creates a reading from the request payload. [file:249] |
| POST | `/api/readings/random` | Creates a randomized reading using coin casting mode. [file:249] |

The frontend is currently served by `app.Static("/", "./web/static")`, so static files must exist in that directory unless the app is later changed to embedded assets. [file:249]

## Run locally

Start the application with:

```bash
go run ./cmd/api
```

The current startup code loads configuration, initializes the selected repository, creates the reading service, builds the Fiber app, and starts listening on `cfg.Addr`. It also installs graceful shutdown handling for `SIGINT` and `SIGTERM`. [file:248]

## Configuration

The current application behavior depends on `config.Load()`, which supplies values such as storage mode, SQLite path, PostgreSQL DSN, and listen address. The startup logic clearly branches by `cfg.Storage`, so that setting is the key switch between persistent and non-persistent modes. [file:248]

Typical local values are:

```env
STORAGE=sqlite
SQLITE_PATH=iching.db
ADDR=:8080
```

If you want PostgreSQL instead, use a valid `POSTGRES_DSN` and set `STORAGE=postgres`. [file:248]

## Windows build

A Windows binary can be built with standard Go cross-compilation using `GOOS=windows` and `GOARCH=amd64`. This is the normal Go approach for producing a `.exe` from the API entrypoint. [web:268]

Example:

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o iching.exe ./cmd/api
```

With the current static-file setup, distributing only the `.exe` is not enough unless you also embed `web/static` into the Go binary. Otherwise the frontend files must still be available on disk. [file:249][web:281]

## Browser behavior

In the currently known `main.go`, the application starts the server and logs the listening address, but it does not automatically launch the browser. Browser auto-open would need to be added explicitly in startup code. [file:248]

## Notes

- The current SQLite repository stores line data and interpretation data as JSON strings in the database. [file:258]
- Reading history is returned in reverse chronological order. [file:258]
- Reflection persistence is present at repository level through `SaveReflection(...)`, but the currently known HTTP layer does not yet expose dedicated reflection endpoints. [file:249][file:258]
- The current repository scans `language` as a plain string, so nullable historical rows may require a migration or a repository fix if older data contains `NULL` in that column. [file:258]

## Suggested next steps

- Add embedded static assets for a true single-binary desktop-style distribution. [web:281]
- Add HTTP endpoints for reflection and relevance statistics to match the richer frontend behavior. [file:249]
- Harden SQLite scanning for nullable fields such as `language`. [file:258]
