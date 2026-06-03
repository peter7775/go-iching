# I Ching App

A Go application for I Ching divination with a Fiber-based HTTP API, a browser frontend, and persistent reading storage. The current project supports SQLite, PostgreSQL, and in-memory repositories.


<img width="1196" height="996" alt="image" src="https://github.com/user-attachments/assets/b229fa5d-9edf-4e09-be43-7c80e8adc5f6" />


## Buy me a coffee, if you like and use this app

[![Buy me a coffee](https://img.shields.io/badge/Buy%20me%20a%20coffee-%23FF813F?style=flat&logo=buy-me-a-coffee&logoColor=white)](https://www.buymeacoffee.com/petrms)

## Quick Start

### Installation Options

#### Option 1: Pre-built Binary (Recommended)

Download the latest release from [GitHub Releases](https://github.com/peter7775/go-iching/releases) and extract it for your platform.

**Linux/macOS**:
```bash
# Download the binary (replace vX.Y.Z with the latest version)
wget https://github.com/peter7775/go-iching/releases/download/vX.Y.Z/iching-api-linux-amd64-vX.Y.Z.tar.gz

# Extract
tar -xzf iching-api-linux-amd64-vX.Y.Z.tar.gz

# Run
chmod +x iching-api-linux-amd64
./iching-api-linux-amd64
```

**Windows**:
```cmd
# Download and extract the ZIP file from GitHub Releases
# Then run:
iching-api-windows-amd64.exe
```

**macOS (ARM64)**:
```bash
# For Apple Silicon (M1/M2/M3)
tar -xzf iching-api-darwin-arm64-vX.Y.Z.tar.gz
chmod +x iching-api-darwin-arm64
./iching-api-darwin-arm64
```

#### Option 2: Build from Source

**Prerequisites**:
- Go 1.21 or later
- Make

**Steps**:
```bash
# Clone the repository
git clone https://github.com/peter7775/go-iching.git
cd go-iching

# Run the application
go run ./cmd/api
```

The app will start on `http://localhost:8080` by default.

#### Option 3: Docker

```bash
# Build Docker image
docker build -t iching-app .

# Run container
docker run -p 8080:8080 \
  -e STORAGE=sqlite \
  -e SQLITE_PATH=iching.db \
  -v $(pwd)/data:/app/data \
  iching-app
```

Access the app at `http://localhost:8080`.

---

## Configuration

Configure the application using environment variables:

```env
# Storage backend: sqlite, postgres, or memory (default: sqlite)
STORAGE=sqlite

# SQLite configuration
SQLITE_PATH=iching.db

# PostgreSQL configuration (if STORAGE=postgres)
POSTGRES_DSN=postgres://user:password@localhost:5432/iching

# Server address (default: :8080)
ADDR=:8080
```

### Storage Backend Details

**SQLite** (Default - recommended for single-user):
- No external dependencies
- Readings stored in local database
- Perfect for personal use or small deployments

**PostgreSQL** (For shared/production deployments):
- Supports multiple concurrent users
- Requires PostgreSQL server
- Connection via DSN: `postgres://user:password@host:port/database`

**In-Memory** (Testing only):
- No persistence
- Useful for testing
- Readings lost on restart

---

## Overview

The application is structured as a small web app. A Go backend handles reading generation and persistence, while the frontend is served as static files from `cmd/api/static`. The API currently exposes REST endpoints for CRUD operations on readings.

Each stored reading includes:
- question text
- cast method
- language
- raw line values
- primary and relating hexagram numbers
- serialized interpretation
- optional reflection rating, note, and timestamp
- creation timestamp

## Architecture

Main components in the current project:

- `cmd/api` — application entrypoint and startup logic
- `internal/httpfiber` — Fiber app configuration and HTTP routes
- `internal/service` — application and reading logic
- `internal/storage/sqlite` — SQLite persistence for saved readings
- `internal/storage/postgres` — PostgreSQL repository option
- `internal/storage/memory` — in-memory fallback repository
- `web/static` — browser UI served by the backend

## API Endpoints

The current HTTP layer exposes these routes:

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Returns a simple health response |
| GET | `/api/readings` | Returns the saved reading history |
| GET | `/api/readings/:id` | Returns one saved reading by ID |
| POST | `/api/readings` | Creates a reading from the request payload |
| POST | `/api/readings/random` | Creates a randomized reading using coin casting mode |

## Storage Backends

The repository is selected by application configuration. If `cfg.Storage` is set to `sqlite`, the app opens `cfg.SQLitePath` and uses the SQLite repository. If it is set to `postgres`, the app connects via the PostgreSQL DSN.

### SQLite

The SQLite repository creates a `readings` table if it does not exist and stores both structured fields and serialized JSON payloads. The current schema includes `question`, `method`, `language`, and other fields.

Reads are ordered by `created_at DESC`, which means the latest saved readings appear first in history. The repository also supports loading a single reading by ID and updating reflection data.

### PostgreSQL

PostgreSQL is available as an alternative repository and is activated when the storage mode is `postgres`. The application opens the connection using `sql.Open("pgx", cfg.PostgresDSN)`.

### In-Memory Mode

In-memory storage is only a fallback. It is useful for tests or temporary runs, but readings are not persisted across process restarts.

## Development

### Run Locally

Start the application with:

```bash
go run ./cmd/api
```

The current startup code loads configuration, initializes the selected repository, creates the reading service, builds the Fiber app, and starts listening on `cfg.Addr`. It also installs graceful shutdown handlers.

### Build for Specific Platform

**Windows x86_64**:
```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o iching.exe ./cmd/api
```

**Linux x86_64**:
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o iching-api ./cmd/api
```

**macOS (ARM64)**:
```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o iching-api ./cmd/api
```

## Release Management

### Local Release Build

To build release packages for all platforms (Linux, Windows, macOS):

```bash
make release-build
```

This will:
1. Run code quality checks (tidy, vet, lint, test)
2. Build binaries for all platforms (x86_64 and ARM64 variants)
3. Create compressed archives in `dist/` directory

Individual platform builds:
```bash
make build-linux              # Linux x86_64
make build-linux-arm64        # Linux ARM64
make build-windows            # Windows x86_64
make build-windows-arm64      # Windows ARM64
make build-darwin             # macOS x86_64
make build-darwin-arm64       # macOS ARM64
```

### Publish a Release

1. **Prepare code**: Ensure all changes are committed and tests pass
   ```bash
   make audit  # Runs tidy, vet, lint, test
   ```

2. **Create a release tag**: Tags should follow semantic versioning (e.g., `v0.1.0`, `v1.0.0-beta`)
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0: Initial release"
   git push origin v0.1.0
   ```

3. **GitHub Actions**: When you push a tag starting with `v`, GitHub Actions automatically:
   - Builds binaries for all platforms
   - Creates compressed archives
   - Creates a GitHub Release with all artifacts

4. **Manual Release Build** (if not using GitHub Actions):
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   make release-build
   ```

### Distribution Files

Release packages include:
- **Linux**: `iching-api-linux-amd64-vX.Y.Z.tar.gz`, `iching-api-linux-arm64-vX.Y.Z.tar.gz`
- **Windows**: `iching-api-windows-amd64-vX.Y.Z.zip`, `iching-api-windows-arm64-vX.Y.Z.zip`
- **macOS**: `iching-api-darwin-amd64-vX.Y.Z.tar.gz`, `iching-api-darwin-arm64-vX.Y.Z.tar.gz`

Each archive contains:
- Compiled binary for the platform
- `static/` directory with web UI files

## Browser Behavior

The application starts the server and logs the listening address. Open your browser and navigate to the address shown in the console (typically `http://localhost:8080`).

## Notes

- The current SQLite repository stores line data and interpretation data as JSON strings in the database
- Reading history is returned in reverse chronological order
- Reflection persistence is present at repository level through `SaveReflection(...)`, but the currently known HTTP layer does not yet expose dedicated reflection endpoints
- The current repository scans `language` as a plain string, so nullable historical rows may require a migration or a repository fix if older data contains `NULL` in that column

## License

This project is released under the Custom Attribution Software License v1.0.

## Attribution

If you reuse, redistribute, or modify this software, you must keep the copyright notice and credit:

Petr Štěpánek <petrstepanek99@proton.me>

See [LICENSE](LICENSE) for full details.
