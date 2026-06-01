# Wails stub

Wails je vhodný pro desktop výstup, protože staví desktop aplikaci nad Go backendem a web technologiemi v jednom projektu.

## Start
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails init -n iching-desktop -t vanilla
```

Pak můžeš sdílet `internal/domain` a `internal/service` mezi Fiber backendem a desktop klientem.
