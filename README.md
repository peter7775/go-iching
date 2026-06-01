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

## Doporučený další krok
1. Nahradit in-memory repository PostgreSQL/SQLite implementací.
2. Doplnit plný dataset 64 hexagramů a texty čar.
3. Přidat autentizaci a synchronizaci.
4. Rozšířit Wails desktop klient.
