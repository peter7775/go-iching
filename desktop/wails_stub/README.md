# Wails stub

Tato složka je kostra pro desktop klienta ve Wails.

## Inicializace
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails init -n iching-desktop -t vanilla
```

Potom přesuň nebo zkopíruj sdílenou logiku z `internal/service` a `internal/domain` do Wails projektu,
nebo použij tento repozitář jako monorepo a připoj Wails frontend na stejné servisní metody. Wails podle dokumentace balí Go backend a web frontend do jedné desktop aplikace. [viz README projektu]
