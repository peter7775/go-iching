# Release Guide

## Preparing for Release

If you want to create a new release of the I Ching application, follow this guide.

### Step 1: Verify the Code

```bash
# Run all checks
make audit

# Or individually:
make tidy      # go mod tidy
make vet       # go vet ./...
make lint      # golangci-lint
make test      # go test ./...
```

### Step 2: Create a Git Tag

```bash
# Example for version 1.0.0
git tag -a v1.0.0 -m "Release v1.0.0: Description"

# Or with a more detailed message
git tag -a v1.0.0 -F - << 'EOF'
Release v1.0.0: Major features

- Feature 1
- Feature 2
- Bug fix: Something

Signed-by: Your Name <your@email.com>
EOF
```

### Step 3: Push Tag to GitHub

**Locally:**
```bash
git push origin v1.0.0
```

**Once you push the tag, GitHub Actions automatically:**
- Runs builds for all platforms
- Creates a GitHub Release
- Uploads binaries

### Step 4: Verify Release on GitHub

1. Go to: https://github.com/your-org/go-iching/releases
2. You should see the new release with assets and descriptions
3. Download the binaries and test them

---

## Local Release Build (without GitHub)

If you want to build a release locally without GitHub Actions:

```bash
# Make sure you have the correct git tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Run release build
make release-build
```

The resulting packages will be in `dist/`:
- `iching-api-linux-amd64-v1.0.0.tar.gz`
- `iching-api-linux-arm64-v1.0.0.tar.gz`
- `iching-api-windows-amd64-v1.0.0.zip`
- `iching-api-windows-arm64-v1.0.0.zip`
- `iching-api-darwin-amd64-v1.0.0.tar.gz`

---

## Installation from Release

### Linux/macOS

```bash
# Download and extract
tar -xzf iching-api-linux-amd64-v1.0.0.tar.gz

# First run
./iching-api-linux-amd64
```

### Windows

```cmd
# Extract ZIP
# Run iching-api-windows-amd64.exe
iching-api-windows-amd64.exe
```

---

## Versioning (Semantic Versioning)

Follow semantic versioning: `MAJOR.MINOR.PATCH-PRERELEASE+BUILD`

Examples:
- `v1.0.0` - first stable release
- `v1.1.0` - new features (backward compatible)
- `v1.0.1` - bug fix
- `v1.0.0-rc1` - release candidate
- `v1.0.0-alpha` - alpha version

### When to Bump Version?

- **MAJOR**: breaking changes (incompatible changes)
- **MINOR**: new features (compatible)
- **PATCH**: bug fixes

---

## Pre-Release Checks

```bash
# 1. Make sure everything is committed
git status

# 2. Run tests and linting
make audit

# 3. Review recent commits
git log --oneline -n 10

# 4. Check existing tags
git tag

# 5. Only then create a new tag
git tag -a v1.0.0 -m "Release v1.0.0"

# 6. Push the tag
git push origin v1.0.0
```

---

## Troubleshooting

### Build fails with errors
- Run `make audit` and fix all issues
- Check your Go version: `go version` (should be 1.23.0+)

### GitHub Actions workflow doesn't respond
- Check `.github/workflows/release.yml`
- Go to the GitHub Actions tab and check the logs

### Binary doesn't work on Windows
- Make sure you have `static/` files next to the `.exe`
- Windows binary needs access to the `static/` directory

---

More questions? Check out README.md in the "Release Management" section.

