# Release Guide

## Příprava na release

Pokud chcete vytvořit nový release aplikace I Ching, postupujte dle tohoto průvodce.

### Krok 1: Ověřte kód

```bash
# Spusťte všechny kontroly
make audit

# Nebo jednotlivě:
make tidy      # go mod tidy
make vet       # go vet ./...
make lint      # golangci-lint
make test      # go test ./...
```

### Krok 2: Vytvořte git tag

```bash
# Případ pro verzi 1.0.0
git tag -a v1.0.0 -m "Release v1.0.0: Description"

# Nebo s podrobnější zprávou
git tag -a v1.0.0 -F - << 'EOF'
Release v1.0.0: Major features

- Feature 1
- Feature 2
- Bug fix: Something

Signed-by: Your Name <your@email.com>
EOF
```

### Krok 3: Push tagu na GitHub

**Lokálně:**
```bash
git push origin v1.0.0
```

**Jakmile pushneš tag, GitHub Actions automaticky:**
- Spustí build pro všechny platformy
- Vytvoří GitHub Release
- Nahraje binárky

### Krok 4: Ověřte release na GitHubu

1. Jděte na: https://github.com/your-org/go-iching/releases
2. Měli byste vidět nový release s objekty a popisy
3. Stáhněte si binárky a otestujte

---

## Lokální build release (bez GitHub)

Pokud budete dělat release lokálně bez GitHub Actions:

```bash
# Ujistěte se, že máte správný git tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Spusťte release build
make release-build
```

Výsledné balíčky budou v `dist/`:
- `iching-api-linux-amd64-v1.0.0.tar.gz`
- `iching-api-linux-arm64-v1.0.0.tar.gz`
- `iching-api-windows-amd64-v1.0.0.zip`
- `iching-api-windows-arm64-v1.0.0.zip`
- `iching-api-darwin-amd64-v1.0.0.tar.gz`

---

## Instalace z release

### Linux/macOS

```bash
# Stáhnout a rozbalit
tar -xzf iching-api-linux-amd64-v1.0.0.tar.gz

# Prvního spuštění
./iching-api-linux-amd64
```

### Windows

```cmd
# Rozbalit ZIP
# Spustit iching-api-windows-amd64.exe
iching-api-windows-amd64.exe
```

---

## Versionování (Semantic Versioning)

Dodržujte sémantické verzování: `MAJOR.MINOR.PATCH-PRERELEASE+BUILD`

Příklady:
- `v1.0.0` - prvnı stabilní release
- `v1.1.0` - nová features (backward compatible)
- `v1.0.1` - bug fix
- `v1.0.0-rc1` - release candidate
- `v1.0.0-alpha` - alpha verze

### Kdy zvýšit verzi?

- **MAJOR**: breaking changes (ne-kompatible změny)
- **MINOR**: nové features (kompatibilní)
- **PATCH**: bug fixes

---

## Checks před release

```bash
# 1. Ujistěte se, že je vše committováno
git status

# 2. Spusťte testy a lint
make audit

# 3. Prohlédněte si poslední commity
git log --oneline -n 10

# 4. Zkontrolujte existující tagy
git tag

# 5. Teprve pak vytvoříte nový tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 6. Push tagu
git push origin v1.0.0
```

---

## Troubleshooting

### Build selže s errory
- Spusťte `make audit` a opravte všechny problémy
- Checkněte si Go verzi: `go version` (měla by být 1.23.0+)

### GitHub Actions workflow nereaguje
- Checkněte `.github/workflows/release.yml`
- Jděte na GitHub Actions záložku a podívejte se na logy

### Binárka nepracuje na Windowsech
- Ujistěte se, že máte `static/` soubory vedle `.exe`
- Windows binárka potřebuje přístup k `static/` adresáři

---

Další otázky? Podívejte se na README.md v sekci "Release Management".

