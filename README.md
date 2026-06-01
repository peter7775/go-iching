# I Ching app template (Go)

Základ projektu pro I-ťing / Yi Jing aplikaci s jedním Go backendem a dvěma výstupy:
- web/cloud frontend
- standalone desktop frontend

## Architektura
- `cmd/api` — HTTP API server
- `internal/domain` — doménové typy a pravidla I-ťingu
- `internal/service` — aplikační logika
- `internal/storage` — repository vrstva
- `web/` — jednoduchý web frontend
- `desktop/wails_stub` — kostra pro Wails desktop app
- `data/` — referenční data hexagramů

## Rychlý start
```bash
go run ./cmd/api
```
Pak otevři http://localhost:8080
